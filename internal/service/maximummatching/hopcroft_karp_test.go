package maximummatching

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"matching-engine/internal/enums"
	"matching-engine/internal/model"
)

func minimalOffer(id string) *model.Offer {
	coord, _ := model.NewCoordinate(0, 0)
	return model.NewOffer(
		id,
		"user1",
		*coord,
		*coord,
		time.Now(),
		0,
		1,
		*model.NewPreference(enums.Male, false, false, false),
		time.Now().Add(time.Hour),
		0,
		[]model.PathPoint{*model.NewPathPoint(*coord, enums.Source, time.Now(), nil, 0), *model.NewPathPoint(*coord, enums.Destination, time.Now().Add(time.Hour), nil, 0)},
		nil,
	)
}

func minimalRequest(id string) *model.Request {
	coord, _ := model.NewCoordinate(0, 0)
	return model.NewRequest(
		id,
		"user2",
		*coord,
		*coord,
		time.Now(),
		time.Now().Add(time.Hour),
		0,
		1,
		*model.NewPreference(enums.Female, false, false, false),
	)
}

func minimalEdge(requestNode *model.RequestNode) *model.Edge {
	return model.NewEdge(requestNode, nil)
}

func TestHopcroftKarp_FindMaximumMatching(t *testing.T) {
	tests := []struct {
		name     string
		graph    *model.Graph
		wantSize int
		wantErr  bool
	}{
		{
			name:     "empty graph",
			graph:    model.NewGraph(),
			wantSize: 0,
			wantErr:  false,
		},
		{
			name: "single match",
			graph: func() *model.Graph {
				g := model.NewGraph()
				offer := minimalOffer("offer1")
				request := minimalRequest("request1")
				offerNode := model.NewOfferNode(offer)
				requestNode := model.NewRequestNode(request)
				offerNode.SetEdges([]*model.Edge{minimalEdge(requestNode)})
				g.AddOfferNode(offerNode)
				g.AddRequestNode(requestNode)
				g.AddEdge(offer, request, offerNode.Edges()[0])
				return g
			}(),
			wantSize: 1,
			wantErr:  false,
		},
		{
			name: "no possible matches",
			graph: func() *model.Graph {
				g := model.NewGraph()
				offer := minimalOffer("offer1")
				request := minimalRequest("request1")
				offerNode := model.NewOfferNode(offer)
				requestNode := model.NewRequestNode(request)
				g.AddOfferNode(offerNode)
				g.AddRequestNode(requestNode)
				return g
			}(),
			wantSize: 0,
			wantErr:  false,
		},
		{
			name: "multiple possible matches",
			graph: func() *model.Graph {
				g := model.NewGraph()
				offer1 := minimalOffer("offer1")
				offer2 := minimalOffer("offer2")
				request1 := minimalRequest("request1")
				request2 := minimalRequest("request2")
				offerNode1 := model.NewOfferNode(offer1)
				offerNode2 := model.NewOfferNode(offer2)
				requestNode1 := model.NewRequestNode(request1)
				requestNode2 := model.NewRequestNode(request2)
				offerNode1.SetEdges([]*model.Edge{minimalEdge(requestNode1), minimalEdge(requestNode2)})
				offerNode2.SetEdges([]*model.Edge{minimalEdge(requestNode1), minimalEdge(requestNode2)})
				g.AddOfferNode(offerNode1)
				g.AddOfferNode(offerNode2)
				g.AddRequestNode(requestNode1)
				g.AddRequestNode(requestNode2)
				g.AddEdge(offer1, request1, offerNode1.Edges()[0])
				g.AddEdge(offer1, request2, offerNode1.Edges()[1])
				g.AddEdge(offer2, request1, offerNode2.Edges()[0])
				g.AddEdge(offer2, request2, offerNode2.Edges()[1])
				return g
			}(),
			wantSize: 2,
			wantErr:  false,
		},
		{
			name: "all possible matches",
			graph: func() *model.Graph {
				g := model.NewGraph()
				offer1 := minimalOffer("offer1")
				offer2 := minimalOffer("offer2")
				request1 := minimalRequest("request1")
				request2 := minimalRequest("request2")
				offerNode1 := model.NewOfferNode(offer1)
				offerNode2 := model.NewOfferNode(offer2)
				requestNode1 := model.NewRequestNode(request1)
				requestNode2 := model.NewRequestNode(request2)
				offerNode1.SetEdges([]*model.Edge{minimalEdge(requestNode1)})
				offerNode2.SetEdges([]*model.Edge{minimalEdge(requestNode2)})
				g.AddOfferNode(offerNode1)
				g.AddOfferNode(offerNode2)
				g.AddRequestNode(requestNode1)
				g.AddRequestNode(requestNode2)
				g.AddEdge(offer1, request1, offerNode1.Edges()[0])
				g.AddEdge(offer2, request2, offerNode2.Edges()[0])
				return g
			}(),
			wantSize: 2,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hk := NewHopcroftKarp()
			got, err := hk.FindMaximumMatching(tt.graph)
			if (err != nil) != tt.wantErr {
				t.Errorf("HopcroftKarp.FindMaximumMatching() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				if tt.wantSize != 0 {
					t.Errorf("HopcroftKarp.FindMaximumMatching() got nil result, want size %v", tt.wantSize)
				}
				return
			}
			if got.Size() != tt.wantSize {
				t.Errorf("HopcroftKarp.FindMaximumMatching() got size = %v, want %v", got.Size(), tt.wantSize)
			}
		})
	}
}

func createSmallTestGraph() *model.Graph {
	offer1 := minimalOffer("offer1")
	offer2 := minimalOffer("offer2")
	offer3 := minimalOffer("offer3")
	offer4 := minimalOffer("offer4")
	offer5 := minimalOffer("offer5")
	offer6 := minimalOffer("offer6")
	request1 := minimalRequest("request1")
	request2 := minimalRequest("request2")
	request3 := minimalRequest("request3")
	request4 := minimalRequest("request4")
	request5 := minimalRequest("request5")
	request6 := minimalRequest("request6")
	offerNode1 := model.NewOfferNode(offer1)
	offerNode2 := model.NewOfferNode(offer2)
	offerNode3 := model.NewOfferNode(offer3)
	offerNode4 := model.NewOfferNode(offer4)
	offerNode5 := model.NewOfferNode(offer5)
	offerNode6 := model.NewOfferNode(offer6)
	requestNode1 := model.NewRequestNode(request1)
	requestNode2 := model.NewRequestNode(request2)
	requestNode3 := model.NewRequestNode(request3)
	requestNode4 := model.NewRequestNode(request4)
	requestNode5 := model.NewRequestNode(request5)
	requestNode6 := model.NewRequestNode(request6)
	g := model.NewGraph()
	g.AddOfferNode(offerNode1)
	g.AddOfferNode(offerNode2)
	g.AddOfferNode(offerNode3)
	g.AddOfferNode(offerNode4)
	g.AddOfferNode(offerNode5)
	g.AddOfferNode(offerNode6)
	g.AddRequestNode(requestNode1)
	g.AddRequestNode(requestNode2)
	g.AddRequestNode(requestNode3)
	g.AddRequestNode(requestNode4)
	g.AddRequestNode(requestNode5)
	g.AddRequestNode(requestNode6)
	offerNode1.SetEdges([]*model.Edge{minimalEdge(requestNode2), minimalEdge(requestNode3)})
	offerNode3.SetEdges([]*model.Edge{minimalEdge(requestNode1), minimalEdge(requestNode4)})
	offerNode4.SetEdges([]*model.Edge{minimalEdge(requestNode3)})
	offerNode5.SetEdges([]*model.Edge{minimalEdge(requestNode3), minimalEdge(requestNode4)})
	offerNode6.SetEdges([]*model.Edge{minimalEdge(requestNode6)})
	g.AddEdge(offer1, request2, offerNode1.Edges()[0])
	g.AddEdge(offer1, request3, offerNode1.Edges()[1])
	g.AddEdge(offer3, request1, offerNode3.Edges()[0])
	g.AddEdge(offer4, request3, offerNode4.Edges()[0])
	g.AddEdge(offer5, request3, offerNode5.Edges()[0])
	g.AddEdge(offer5, request4, offerNode5.Edges()[1])
	g.AddEdge(offer6, request6, offerNode6.Edges()[0])
	return g
}

func TestHopcroftKarp_FindMaximumMatching_SmallCase(t *testing.T) {
	g := createSmallTestGraph()
	hk := NewHopcroftKarp()
	got, err := hk.FindMaximumMatching(g)
	if err != nil {
		t.Fatalf("HopcroftKarp.FindMaximumMatching() error: %v", err)
	}
	if got == nil {
		t.Fatalf("HopcroftKarp.FindMaximumMatching() got nil result")
	}
	if got.Size() != 5 {
		t.Errorf("HopcroftKarp.FindMaximumMatching() got size = %v, want 5", got.Size())
	}
	got.Range(func(offerNode *model.OfferNode, edge *model.Edge) bool {
		requestID := edge.RequestNode().Request().ID()
		t.Logf("Offer %s matched with request %s", offerNode.Offer().ID(), requestID)
		return true
	})
}

func TestHopcroftKarp_FindMaximumMatching_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		graph    *model.Graph
		wantSize int
		wantErr  bool
	}{
		{
			name: "graph with missing edge",
			graph: func() *model.Graph {
				g := model.NewGraph()
				offer := minimalOffer("offer1")
				request := minimalRequest("request1")
				offerNode := model.NewOfferNode(offer)
				requestNode := model.NewRequestNode(request)
				g.AddOfferNode(offerNode)
				g.AddRequestNode(requestNode)
				// No edge added
				return g
			}(),
			wantSize: 0,
			wantErr:  false,
		},
		{
			name: "graph with only offers",
			graph: func() *model.Graph {
				g := model.NewGraph()
				offer1 := minimalOffer("offer1")
				offer2 := minimalOffer("offer2")
				g.AddOfferNode(model.NewOfferNode(offer1))
				g.AddOfferNode(model.NewOfferNode(offer2))
				return g
			}(),
			wantSize: 0,
			wantErr:  false,
		},
		{
			name: "graph with only requests",
			graph: func() *model.Graph {
				g := model.NewGraph()
				request1 := minimalRequest("request1")
				request2 := minimalRequest("request2")
				g.AddRequestNode(model.NewRequestNode(request1))
				g.AddRequestNode(model.NewRequestNode(request2))
				return g
			}(),
			wantSize: 0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hk := NewHopcroftKarp()
			got, err := hk.FindMaximumMatching(tt.graph)
			if (err != nil) != tt.wantErr {
				t.Errorf("HopcroftKarp.FindMaximumMatching() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				if tt.wantSize != 0 {
					t.Errorf("HopcroftKarp.FindMaximumMatching() got nil result, want size %v", tt.wantSize)
				}
				return
			}
			if got.Size() != tt.wantSize {
				t.Errorf("HopcroftKarp.FindMaximumMatching() got size = %v, want %v", got.Size(), tt.wantSize)
			}
		})
	}
}

func TestHopcroftKarp_LargeCase(t *testing.T) {
	const (
		nOffers   = 1000
		nRequests = 1000
		nRuns     = 10
	)

	var totalMatches int
	var totalTime time.Duration

	for run := 0; run < nRuns; run++ {
		// Get initial memory stats
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		initialAlloc := m.TotalAlloc

		g := model.NewGraph()
		offerNodes := make([]*model.OfferNode, nOffers)
		requestNodes := make([]*model.RequestNode, nRequests)

		// Create all offers and requests first
		for i := 0; i < nOffers; i++ {
			offer := minimalOffer(fmt.Sprintf("offer%d", i))
			offerNode := model.NewOfferNode(offer)
			g.AddOfferNode(offerNode)
			offerNodes[i] = offerNode
		}

		for i := 0; i < nRequests; i++ {
			request := minimalRequest(fmt.Sprintf("request%d", i))
			requestNode := model.NewRequestNode(request)
			g.AddRequestNode(requestNode)
			requestNodes[i] = requestNode
		}

		// Get memory stats after creating nodes
		runtime.ReadMemStats(&m)
		nodesAlloc := m.TotalAlloc

		// Create a complete bipartite graph
		// Each offer can match with every request
		for i := 0; i < nOffers; i++ {
			edges := make([]*model.Edge, nRequests)
			for j := 0; j < nRequests; j++ {
				edge := minimalEdge(requestNodes[j])
				edges[j] = edge
				g.AddEdge(offerNodes[i].Offer(), requestNodes[j].Request(), edge)
			}
			offerNodes[i].SetEdges(edges)
		}

		// Get memory stats after creating edges
		runtime.ReadMemStats(&m)
		edgesAlloc := m.TotalAlloc

		// Run the matching algorithm
		hk := NewHopcroftKarp()
		start := time.Now()
		result, err := hk.FindMaximumMatching(g)
		elapsed := time.Since(start)

		// Get final memory stats
		runtime.ReadMemStats(&m)
		finalAlloc := m.TotalAlloc
		finalHeap := m.HeapAlloc

		if err != nil {
			t.Fatalf("HopcroftKarp.FindMaximumMatching() error: %v", err)
		}

		// Calculate statistics
		totalPossibleMatches := nOffers * nRequests // Complete graph

		// Print memory statistics
		fmt.Printf("\nRun %d Memory Statistics:\n", run+1)
		fmt.Printf("- Initial memory allocation: %.2f MB\n", float64(initialAlloc)/1024/1024)
		fmt.Printf("- Memory after creating nodes: %.2f MB (delta: %.2f MB)\n",
			float64(nodesAlloc)/1024/1024,
			float64(nodesAlloc-initialAlloc)/1024/1024)
		fmt.Printf("- Memory after creating edges: %.2f MB (delta: %.2f MB)\n",
			float64(edgesAlloc)/1024/1024,
			float64(edgesAlloc-nodesAlloc)/1024/1024)
		fmt.Printf("- Final memory allocation: %.2f MB (delta: %.2f MB)\n",
			float64(finalAlloc)/1024/1024,
			float64(finalAlloc-edgesAlloc)/1024/1024)
		fmt.Printf("- Peak heap usage: %.2f MB\n", float64(finalHeap)/1024/1024)

		// Print matching statistics
		fmt.Printf("\nRun %d Matching Statistics:\n", run+1)
		fmt.Printf("- Total offers: %d\n", nOffers)
		fmt.Printf("- Total requests: %d\n", nRequests)
		fmt.Printf("- Total possible matches: %d\n", totalPossibleMatches)
		fmt.Printf("- Actual matches found: %d\n", result.Size())
		fmt.Printf("- Time taken: %v\n", elapsed)

		// Verify that the matching is valid
		matchedRequests := make(map[string]bool)
		result.Range(func(offerNode *model.OfferNode, edge *model.Edge) bool {
			requestID := edge.RequestNode().Request().ID()
			if matchedRequests[requestID] {
				t.Errorf("Request %s is matched multiple times", requestID)
			}
			matchedRequests[requestID] = true
			return true
		})

		// Verify that we got the maximum possible matching
		expectedMatches := nOffers // Since nOffers = nRequests
		if result.Size() != expectedMatches {
			t.Errorf("Expected %d matches (maximum possible), got %d", expectedMatches, result.Size())
		}

		// Force garbage collection and get final memory stats
		runtime.GC()
		runtime.ReadMemStats(&m)
		fmt.Printf("\nRun %d Memory after GC: %.2f MB\n", run+1, float64(m.HeapAlloc)/1024/1024)

		totalMatches += result.Size()
		totalTime += elapsed
	}

	fmt.Printf("\nAverage matches over %d runs: %.2f\n", nRuns, float64(totalMatches)/float64(nRuns))
	fmt.Printf("Average time over %d runs: %v\n", nRuns, totalTime/time.Duration(nRuns))
}
