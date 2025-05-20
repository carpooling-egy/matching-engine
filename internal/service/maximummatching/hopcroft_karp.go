// Package maximummatching implements the MaximumMatching interface using Hopcroft–Karp.
package maximummatching

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
	"math"
)

// HopcroftKarp implements MaximumMatching.
type HopcroftKarp struct {
	pairU []int
	pairV []int
	dist  []int
	queue []int
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
	n, m := len(offers), len(requests)
	if n == 0 || m == 0 {
		return []collections.Tuple2[*model.OfferNode, *model.Edge]{}, nil
	}

	// Pre-allocate slices
	hk.pairU = make([]int, n+1)
	hk.pairV = make([]int, m+1)
	hk.dist = make([]int, n+1)
	hk.queue = make([]int, 0, n)

	adj := buildAdjacencyList(offers, requestIndex, n)

	const NIL = 0
	matching := 0

	// Main matching loop
	for hk.bfs(adj, n, NIL) {
		for u := 1; u <= n; u++ {
			if hk.pairU[u] == NIL && hk.dfs(u, adj, NIL) {
				matching++
			}
		}
	}
	// Check if we found a matching
	if matching == 0 {
		return []collections.Tuple2[*model.OfferNode, *model.Edge]{}, nil
	}

	return buildResultMap(graph, offers, requests, hk.pairU)
}

func initialize(graph *model.Graph) ([]*model.OfferNode, []*model.RequestNode, *collections.SyncMap[*model.RequestNode, int]) {
	offers := make([]*model.OfferNode, 0, graph.OfferNodes().Size())
	graph.OfferNodes().Range(func(_ string, node *model.OfferNode) error {
		offers = append(offers, node)
		return nil
	})
	requests := make([]*model.RequestNode, 0, graph.RequestNodes().Size())
	requestIndex := collections.NewSyncMap[*model.RequestNode, int]()
	counter := 0
	graph.RequestNodes().Range(func(_ string, node *model.RequestNode) error {
		requests = append(requests, node)
		counter++
		requestIndex.Set(node, counter)
		return nil
	})
	return offers, requests, requestIndex
}

func buildAdjacencyList(offers []*model.OfferNode, requestIndex *collections.SyncMap[*model.RequestNode, int], n int) [][]int {
	adj := make([][]int, n+1)
	for i, u := range offers {
		edges := u.Edges()
		adj[i+1] = make([]int, 0, len(edges))
		for _, edge := range edges {
			if idx, ok := requestIndex.Get(edge.RequestNode()); ok {
				adj[i+1] = append(adj[i+1], idx)
			}
		}
	}
	return adj
}

func (hk *HopcroftKarp) bfs(adj [][]int, n, NIL int) bool {
	// Reset queue and distances
	hk.queue = hk.queue[:0]
	distNIL := math.MaxInt32

	// Initialize distances
	for u := 1; u <= n; u++ {
		if hk.pairU[u] == NIL {
			hk.dist[u] = 0
			hk.queue = append(hk.queue, u)
		} else {
			hk.dist[u] = math.MaxInt32
		}
	}

	// BFS loop
	for len(hk.queue) > 0 {
		u := hk.queue[0]
		hk.queue = hk.queue[1:]

		if hk.dist[u] < distNIL {
			for _, v := range adj[u] {
				if hk.pairV[v] == NIL {
					distNIL = hk.dist[u] + 1
				} else if hk.dist[hk.pairV[v]] == math.MaxInt32 {
					hk.dist[hk.pairV[v]] = hk.dist[u] + 1
					hk.queue = append(hk.queue, hk.pairV[v])
				}
			}
		}
	}

	return distNIL != math.MaxInt32
}

func (hk *HopcroftKarp) dfs(u int, adj [][]int, NIL int) bool {

	// Try to find an augmenting path
	for _, v := range adj[u] {
		if hk.pairV[v] == NIL || (hk.dist[hk.pairV[v]] == hk.dist[u]+1 && hk.dfs(hk.pairV[v], adj, NIL)) {
			hk.pairU[u] = v
			hk.pairV[v] = u
			return true
		}
	}

	hk.dist[u] = math.MaxInt32
	return false
}

func buildResultMap(graph *model.Graph, offers []*model.OfferNode, requests []*model.RequestNode, pairU []int) ([]collections.Tuple2[*model.OfferNode, *model.Edge], error) {
	result := make([]collections.Tuple2[*model.OfferNode, *model.Edge], len(pairU)-1)
	for uIdx, vIdx := range pairU {
		if uIdx > 0 && vIdx > 0 {
			offer := offers[uIdx-1]
			reqNode := requests[vIdx-1]
			chosen, exists := graph.GetEdge(offer.Offer(), reqNode.Request())
			if !exists || chosen == nil {
				return nil, fmt.Errorf("edge not found for match Offer[%d] -> Request[%d]", uIdx-1, vIdx-1)
			}
			result[uIdx-1] = collections.NewTuple2(offer, chosen)
		}
	}
	return result, nil
}
