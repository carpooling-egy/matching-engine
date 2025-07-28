// Package maximummatching implements the MaximumMatching interface using Hopcroft–Karp.
package maximummatching

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
	"math"
)

const (
	// NIL is a constant used to represent a null value in the matching algorithm.
	NIL = 0
	INF = math.MaxInt32
)

// HopcroftKarp implements MaximumMatching.
type HopcroftKarp struct {
	offerMatches   []int
	requestMatches []int
	distances      []int
}

// NewHopcroftKarp returns a hopcroftKarp using the Hopcroft–Karp algorithm.
func NewHopcroftKarp() MaximumMatching {
	return &HopcroftKarp{}
}

// FindMaximumMatching finds maximum bipartite matching between offer and request nodes.
func (hk *HopcroftKarp) FindMaximumMatching(
	graph *model.MaximumMatchingGraph,
) ([]collections.Tuple2[*model.OfferNode, *model.Edge], error) {
	if graph == nil {
		return nil, fmt.Errorf("graph cannot be nil")
	}

	offers, requests, requestIndex := initialize(graph)
	offerCount, requestCount := len(offers), len(requests)
	if offerCount == 0 || requestCount == 0 {
		return []collections.Tuple2[*model.OfferNode, *model.Edge]{}, nil
	}

	// Pre-allocate slices and queue
	hk.preAllocate(offerCount, requestCount)

	adj := buildAdjacencyList(offers, requestIndex, offerCount)

	matchingCount := 0

	// Main matching loop
	for hk.computeLayerDistances(adj, offerCount) {
		for offerIndex := 1; offerIndex <= offerCount; offerIndex++ {
			if hk.offerMatches[offerIndex] == NIL && hk.findAugmentingPath(offerIndex, adj) {
				matchingCount++
			}
		}
	}
	// Check if we found a matching
	if matchingCount == 0 {
		return []collections.Tuple2[*model.OfferNode, *model.Edge]{}, nil
	}

	return buildMatchingResult(graph, offers, requests, hk.offerMatches)
}

func initialize(graph *model.MaximumMatchingGraph) ([]*model.OfferNode, []*model.RequestNode, *collections.SyncMap[*model.RequestNode, int]) {
	offers := make([]*model.OfferNode, 0, graph.OfferNodes().Size())
	graph.OfferNodes().Range(func(_ string, node *model.OfferNode) error {
		offers = append(offers, node)
		return nil
	})
	requests := make([]*model.RequestNode, 0, graph.RequestNodes().Size())
	requestIndexMap := collections.NewSyncMap[*model.RequestNode, int]()
	index := 0
	graph.RequestNodes().Range(func(_ string, node *model.RequestNode) error {
		requests = append(requests, node)
		index++
		requestIndexMap.Set(node, index)
		return nil
	})
	return offers, requests, requestIndexMap
}

func (hk *HopcroftKarp) preAllocate(offerCount, requestCount int) {
	hk.offerMatches = make([]int, offerCount+1)
	hk.requestMatches = make([]int, requestCount+1)
	hk.distances = make([]int, offerCount+1)
}

func buildAdjacencyList(offers []*model.OfferNode, requestIndexMap *collections.SyncMap[*model.RequestNode, int], n int) [][]int {
	adj := make([][]int, n+1)
	for idx, offer := range offers {
		offerIndex := idx + 1
		edges := offer.Edges()
		adj[offerIndex] = make([]int, 0, len(edges))
		for _, edge := range edges {
			if requestIdx, exists := requestIndexMap.Get(edge.RequestNode()); exists {
				adj[offerIndex] = append(adj[offerIndex], requestIdx)
			}
		}
	}
	return adj
}

func (hk *HopcroftKarp) computeLayerDistances(adj [][]int, offerCount int) bool {
	// Reset distances
	distNIL := INF
	queue := collections.NewQueueWithCapacity[int](offerCount)

	// Initialize distances
	for offerIndex := 1; offerIndex <= offerCount; offerIndex++ {
		if hk.offerMatches[offerIndex] == NIL {
			hk.distances[offerIndex] = 0
			queue.Enqueue(offerIndex)
		} else {
			hk.distances[offerIndex] = INF
		}
	}

	// BFS loop
	for queue.Size() > 0 {
		currentOffer, _ := queue.Dequeue()

		if hk.distances[currentOffer] < distNIL {
			for _, requestIndex := range adj[currentOffer] {
				if hk.requestMatches[requestIndex] == NIL {
					distNIL = hk.distances[currentOffer] + 1
				} else if hk.distances[hk.requestMatches[requestIndex]] == INF {
					hk.distances[hk.requestMatches[requestIndex]] = hk.distances[currentOffer] + 1
					queue.Enqueue(hk.requestMatches[requestIndex])
				}
			}
		}
	}

	return distNIL != INF
}

func (hk *HopcroftKarp) findAugmentingPath(offer int, adj [][]int) bool {

	// Try to find an augmenting path
	for _, request := range adj[offer] {
		if hk.requestMatches[request] == NIL || (hk.distances[hk.requestMatches[request]] == hk.distances[offer]+1 && hk.findAugmentingPath(hk.requestMatches[request], adj)) {
			hk.offerMatches[offer] = request
			hk.requestMatches[request] = offer
			return true
		}
	}

	hk.distances[offer] = INF
	return false
}

func buildMatchingResult(graph *model.MaximumMatchingGraph, offers []*model.OfferNode, requests []*model.RequestNode, offerMatches []int) ([]collections.Tuple2[*model.OfferNode, *model.Edge], error) {
	var result []collections.Tuple2[*model.OfferNode, *model.Edge]
	for offerIndex, requestIndex := range offerMatches {
		if offerIndex > 0 && requestIndex > 0 {
			offer := offers[offerIndex-1]
			reqNode := requests[requestIndex-1]
			chosen, exists := graph.GetEdge(offer.Offer(), reqNode.Request())
			if !exists || chosen == nil {
				return nil, fmt.Errorf("edge not found for match Offer[%d] -> Request[%d]", offerIndex-1, requestIndex-1)
			}
			result = append(result, collections.NewTuple2(offer, chosen))
		}
	}
	return result, nil
}
