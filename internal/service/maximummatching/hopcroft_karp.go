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
func NewHopcroftKarp() *HopcroftKarp {
	return &HopcroftKarp{}
}

// FindMaximumMatching finds maximum bipartite matching between offer and request nodes.
func (hk *HopcroftKarp) FindMaximumMatching(
	graph *model.Graph,
) ([]collections.Tuple2[*model.OfferNode, *model.Edge], error) {
	if graph == nil {
		return nil, fmt.Errorf("graph cannot be nil")
	}

	offers, requests, requestIndex := initialize(graph)
	offerCount, requestCount := len(offers), len(requests)
	if offerCount == 0 || requestCount == 0 {
		return []collections.Tuple2[*model.OfferNode, *model.Edge]{}, nil
	}

	// Pre-allocate slices
	hk.offerMatches = make([]int, offerCount+1)
	hk.requestMatches = make([]int, requestCount+1)
	hk.distances = make([]int, offerCount+1)
	queue := collections.NewQueueWithCapacity[int](offerCount)

	adj := buildAdjacencyList(offers, requestIndex, offerCount)

	matchingCount := 0

	// Main matching loop
	for hk.bfs(adj, queue, offerCount) {
		for offerIndex := 1; offerIndex <= offerCount; offerIndex++ {
			if hk.offerMatches[offerIndex] == NIL && hk.dfs(offerIndex, adj) {
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

func initialize(graph *model.Graph) ([]*model.OfferNode, []*model.RequestNode, *collections.SyncMap[*model.RequestNode, int]) {
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

func buildAdjacencyList(offers []*model.OfferNode, requestIndex *collections.SyncMap[*model.RequestNode, int], n int) [][]int {
	adj := make([][]int, n+1)
	for offerIndex, offer := range offers {
		edges := offer.Edges()
		adj[offerIndex+1] = make([]int, 0, len(edges))
		for _, edge := range edges {
			if requestIdx, exists := requestIndex.Get(edge.RequestNode()); exists {
				adj[offerIndex+1] = append(adj[offerIndex+1], requestIdx)
			}
		}
	}
	return adj
}

func (hk *HopcroftKarp) bfs(adj [][]int, queue *collections.Queue[int], n int) bool {
	// Reset distances
	distNIL := INF

	// Initialize distances
	for u := 1; u <= n; u++ {
		if hk.offerMatches[u] == NIL {
			hk.distances[u] = 0
			queue.Enqueue(u)
		} else {
			hk.distances[u] = INF
		}
	}

	// BFS loop
	for queue.Size() > 0 {
		u, _ := queue.Dequeue()

		if hk.distances[u] < distNIL {
			for _, v := range adj[u] {
				if hk.requestMatches[v] == NIL {
					distNIL = hk.distances[u] + 1
				} else if hk.distances[hk.requestMatches[v]] == INF {
					hk.distances[hk.requestMatches[v]] = hk.distances[u] + 1
					queue.Enqueue(hk.requestMatches[v])
				}
			}
		}
	}

	return distNIL != INF
}

func (hk *HopcroftKarp) dfs(offer int, adj [][]int) bool {

	// Try to find an augmenting path
	for _, request := range adj[offer] {
		if hk.requestMatches[request] == NIL || (hk.distances[hk.requestMatches[request]] == hk.distances[offer]+1 && hk.dfs(hk.requestMatches[request], adj)) {
			hk.offerMatches[offer] = request
			hk.requestMatches[request] = offer
			return true
		}
	}

	hk.distances[offer] = INF
	return false
}

func buildMatchingResult(graph *model.Graph, offers []*model.OfferNode, requests []*model.RequestNode, offerMatches []int) ([]collections.Tuple2[*model.OfferNode, *model.Edge], error) {
	result := make([]collections.Tuple2[*model.OfferNode, *model.Edge], len(offerMatches)-1)
	for offerIndex, requestIndex := range offerMatches {
		if offerIndex > 0 && requestIndex > 0 {
			offer := offers[offerIndex-1]
			reqNode := requests[requestIndex-1]
			chosen, exists := graph.GetEdge(offer.Offer(), reqNode.Request())
			if !exists || chosen == nil {
				return nil, fmt.Errorf("edge not found for match Offer[%d] -> Request[%d]", offerIndex-1, requestIndex-1)
			}
			result[offerIndex-1] = collections.NewTuple2(offer, chosen)
		}
	}
	return result, nil
}
