package earlypruning

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"matching-engine/internal/model"
	"matching-engine/internal/service/checker"
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
	workerCount       = runtime.GOMAXPROCS(0) * 4
	channelBufferSize = 2000
)

func (ci *CandidateIterator) Candidates(
	ctx context.Context,
	g *errgroup.Group,
) <-chan *model.MatchCandidate {
	candidatesChannel := make(chan *model.MatchCandidate, channelBufferSize)
	sem := make(chan struct{}, workerCount)
	var generatorWaitGroup sync.WaitGroup

	for _, offer := range ci.offers {
		for _, request := range ci.requests {
			generatorWaitGroup.Add(1)

			g.Go(func() error {
				defer generatorWaitGroup.Done()

				select {
				case <-ctx.Done():
					return ctx.Err()
				case sem <- struct{}{}:
				}
				defer func() { <-sem }()

				isPotential, err := ci.checker.Check(offer, request)
				if err != nil {
					return fmt.Errorf("checker failed for %s/%s: %w", offer.ID(), request.ID(), err)
				}

				if !isPotential {
					return nil
				}

				matchCandidate := model.NewMatchCandidate(request, offer)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case candidatesChannel <- matchCandidate:
					return nil
				}
			})
		}
	}

	g.Go(func() error {
		generatorWaitGroup.Wait()
		close(candidatesChannel)
		return nil
	})

	return candidatesChannel
}
