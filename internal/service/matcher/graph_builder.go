package matcher

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	stage2WorkerCount = runtime.GOMAXPROCS(0) * 4
	stage3WorkerCount = runtime.GOMAXPROCS(0) * 2
	channelBufferSize = 2000
)

func (matcher *Matcher) buildMatchingGraph(
	graph *model.MaximumMatchingGraph,
) (bool, error) {
	var hasNewEdge atomic.Bool

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	offersChannel := make(chan collections.Tuple2[*model.OfferNode, *collections.Set[string]], channelBufferSize)
	offerRequestNodePairsChannel := make(chan collections.Tuple2[*model.OfferNode, *model.RequestNode], channelBufferSize)

	eg.Go(func() error {
		defer close(offersChannel)
		return matcher.produceOffers(ctx, offersChannel)
	})

	filterWG := sync.WaitGroup{}
	filterWG.Add(stage2WorkerCount)
	for i := 0; i < stage2WorkerCount; i++ {
		eg.Go(func() error {
			defer filterWG.Done()
			return matcher.flattenAndPopulate(ctx, offersChannel, offerRequestNodePairsChannel)
		})
	}
	go func() {
		filterWG.Wait()
		close(offerRequestNodePairsChannel)
	}()

	for i := 0; i < stage3WorkerCount; i++ {
		eg.Go(func() error {
			return matcher.evaluateAndBuild(ctx, offerRequestNodePairsChannel, graph, &hasNewEdge)
		})
	}

	if err := eg.Wait(); err != nil {
		return false, err
	}
	return hasNewEdge.Load(), nil
}

// produceOffers ranges over potentialOfferRequests and sends to offersChannel
func (matcher *Matcher) produceOffers(
	ctx context.Context,
	offersChannel chan<- collections.Tuple2[*model.OfferNode, *collections.Set[string]],
) error {
	log.Error().Msg("Stage 1: Producing offers")
	return matcher.potentialOfferRequests.Range(func(offerID string, requestSet *collections.Set[string]) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		offerNode, exists := matcher.availableOffers.Get(offerID)
		if !exists || offerNode == nil {
			matcher.potentialOfferRequests.Delete(offerID)
			return nil
		}
		offersChannel <- collections.NewTuple2(offerNode, requestSet)
		return nil
	})
}

// flattenAndPopulate resolves IDs to nodes, cleans sets, populates matrix, emits node pairs
func (matcher *Matcher) flattenAndPopulate(
	ctx context.Context,
	offersChannel <-chan collections.Tuple2[*model.OfferNode, *collections.Set[string]],
	offerRequestNodePairsChannel chan<- collections.Tuple2[*model.OfferNode, *model.RequestNode],
) error {
	log.Debug().Msg("Stage 2: Flattening and populating time matrix cache")
	for offerRequestSetPair := range offersChannel {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		offerNode := offerRequestSetPair.First
		requestSet := offerRequestSetPair.Second
		ids := requestSet.ToSlice()

		requestNodes := make([]*model.RequestNode, 0, len(ids))
		for _, requestID := range ids {
			if requestNode, ok := matcher.availableRequests.Get(requestID); ok && requestNode != nil {
				requestNodes = append(requestNodes, requestNode)
			} else {
				requestSet.Remove(requestID)
			}
		}

		if len(requestNodes) == 0 {
			return nil
		}

		if err := matcher.timeMatrixCachePopulator.Populate(offerNode, requestNodes); err != nil {
			return err
		}

		for _, requestNode := range requestNodes {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case offerRequestNodePairsChannel <- collections.NewTuple2(offerNode, requestNode):
			}
		}
	}
	return nil
}

// evaluateAndBuild calls Evaluate, cleans invalids, and updates graph on valid
func (matcher *Matcher) evaluateAndBuild(
	ctx context.Context,
	offerRequestNodePairsChannel <-chan collections.Tuple2[*model.OfferNode, *model.RequestNode],
	graph *model.MaximumMatchingGraph,
	hasNewEdge *atomic.Bool,
) error {
	log.Debug().Msg("Stage 3: Evaluating matches and building graph")
	for offerRequestNodePair := range offerRequestNodePairsChannel {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		offerNode := offerRequestNodePair.First
		requestNode := offerRequestNodePair.Second

		path, valid, err := matcher.matchEvaluator.Evaluate(offerNode, requestNode)
		if err != nil {
			return fmt.Errorf("error evaluating the match %v", err)
		}

		if !valid {
			if requestSet, ok := matcher.potentialOfferRequests.Get(offerNode.Offer().ID()); ok && requestSet != nil {
				requestSet.Remove(requestNode.Request().ID())
			}
			continue
		}

		hasNewEdge.Store(true)
		edge := model.NewEdge(requestNode, path)
		graph.AddOfferNode(offerNode)
		graph.AddRequestNode(requestNode)
		graph.AddEdge(offerNode, requestNode, edge)
	}
	return nil
}
