package earlypruning

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
	"matching-engine/internal/service/checker"
	"math"
	"runtime"
	"sync"
)

type CandidateIterator struct {
	offers   []*model.Offer
	requests []*model.Request
	checker  checker.Checker
}

func NewCandidateIterator(offers []*model.Offer, requests []*model.Request, checker checker.Checker) *CandidateIterator {
	return &CandidateIterator{
		offers:   offers,
		requests: requests,
		checker:  checker,
	}
}

var (
	candidateGenStage1WorkerCount = runtime.GOMAXPROCS(0)
	candidateGenStage2WorkerCount = runtime.GOMAXPROCS(0) * 10
	channelBufferSize             = 2000
)

func (ci *CandidateIterator) Candidates(
	ctx context.Context,
	eg *errgroup.Group,
) <-chan *model.MatchCandidate {
	candidatesChannel := make(chan *model.MatchCandidate, channelBufferSize)
	offerPairChannel := make(chan collections.Tuple2[*model.Offer, *model.Request], channelBufferSize)

	ci.produceOfferPairs(ctx, eg, offerPairChannel)

	var consWG sync.WaitGroup
	consWG.Add(candidateGenStage2WorkerCount)
	for i := 0; i < candidateGenStage2WorkerCount; i++ {
		eg.Go(func() error {
			defer consWG.Done()
			return ci.processOfferPairs(ctx, offerPairChannel, candidatesChannel)
		})
	}

	eg.Go(func() error {
		consWG.Wait()
		close(candidatesChannel)
		return nil
	})

	return candidatesChannel
}

func (ci *CandidateIterator) produceOfferPairs(
	ctx context.Context,
	eg *errgroup.Group,
	offerPairChannel chan<- collections.Tuple2[*model.Offer, *model.Request],
) {
	var prodWG sync.WaitGroup
	prodWG.Add(candidateGenStage1WorkerCount)

	nOffers := len(ci.offers)
	nReqs := len(ci.requests)

	rows := min(candidateGenStage1WorkerCount, nOffers)
	cols := int(math.Ceil(float64(candidateGenStage1WorkerCount) / float64(rows)))

	offersPerRow := int(math.Ceil(float64(nOffers) / float64(rows)))
	reqsPerCol := int(math.Ceil(float64(nReqs) / float64(cols)))

	for tile := 0; tile < candidateGenStage1WorkerCount; tile++ {
		row := tile / cols
		col := tile % cols

		oStart := row * offersPerRow
		oEnd := min(oStart+offersPerRow, nOffers)

		rStart := col * reqsPerCol
		rEnd := min(rStart+reqsPerCol, nReqs)

		if oStart >= oEnd || rStart >= rEnd {
			prodWG.Done()
			continue
		}

		eg.Go(func() error {
			defer prodWG.Done()
			for _, offer := range ci.offers[oStart:oEnd] {
				for _, request := range ci.requests[rStart:rEnd] {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case offerPairChannel <- collections.NewTuple2[*model.Offer, *model.Request](offer, request):
					}
				}
			}
			return nil
		})
	}

	eg.Go(func() error {
		prodWG.Wait()
		close(offerPairChannel)
		return nil
	})
}

func (ci *CandidateIterator) processOfferPairs(
	ctx context.Context,
	offerPairChannel <-chan collections.Tuple2[*model.Offer, *model.Request],
	candidatesChannel chan<- *model.MatchCandidate,
) error {
	for offerRequestPair := range offerPairChannel {
		if err := ctx.Err(); err != nil {
			return err
		}

		offer := offerRequestPair.First
		request := offerRequestPair.Second

		isPotential, err := ci.checker.Check(offer, request)
		if err != nil {
			return fmt.Errorf("checker failed for %s/%s: %w",
				offer.ID(), request.ID(), err)
		}
		if !isPotential {
			continue
		}

		matchCandidate := model.NewMatchCandidate(request, offer)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case candidatesChannel <- matchCandidate:
		}
	}
	return nil
}
