package generator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"iter"
	"matching-engine/internal/model"
	"math/rand"
	"time"
)

type RandomTopologicalGenerator struct {
	// RandomTopologicalGenerator is a path generator that generates paths using random sampling
	k     int                         // Number of random samples to generate
	graph *model.TopologicalPathGraph // Graph representing the paths
}

func NewRandomTopologicalGenerator() PathGenerator {
	k := getNumberOfSamples()
	return &RandomTopologicalGenerator{
		k:     k,
		graph: model.NewTopologicalPathGraph(),
	}
}

// GeneratePaths generates k random paths from the path graph
func (rsg *RandomTopologicalGenerator) GeneratePaths(
	path []model.PathPoint,
	pickup, dropoff *model.PathPoint,
) (iter.Seq2[[]model.PathPoint, error], error) {
	// Create an iterator that generates k random paths
	if rsg.k <= 0 {
		return nil, nil // No paths to generate
	}
	pathLength := len(path)
	if pathLength < 2 {
		return nil, fmt.Errorf("path must contain at least two points")
	}
	if pickup == nil || dropoff == nil {
		return nil, fmt.Errorf("pickup and dropoff points must be provided")
	}
	// We are initializing the graph with the provided path, pickup, and drop off points.
	// We are adding the start node of the graph as the first point in the path which is the driver source
	// and the end node as the last point which is the driver destination.
	rsg.graph.InitPathGraph(path, pickup, dropoff, &path[0], &path[len(path)-1])
	return func(yield func([]model.PathPoint, error) bool) {
		r := rand.New(rand.NewSource(time.Now().UnixNano())) // Initialize random seed
		count := 0
		newPath := []model.PathPoint{*rsg.graph.StartNode()}
		visited := make(map[model.PathPointID]bool)
		visited[rsg.graph.StartNode().ID()] = true
		tempInDegree := rsg.graph.CopyInDegree()
		log.Debug().Msgf("RandomTopologicalGenerator: generating paths with k = %d", rsg.k)
		if !rsg.randomBacktrack(tempInDegree, visited, newPath, yield, &count, rsg.k, r) {
			log.Debug().Msgf("RandomTopologicalGenerator: stopped generating paths after %d samples", count)
			rsg.graph.Clear()
			return // Stop generating paths if yield returns false or count exceeds k
		}
	}, nil
}

func (rsg *RandomTopologicalGenerator) randomBacktrack(tempInDegree map[model.PathPointID]int, visited map[model.PathPointID]bool, path []model.PathPoint, yield func([]model.PathPoint, error) bool, count *int, k int, r *rand.Rand) bool {
	if len(path) == int(rsg.graph.Nodes().Size()) && path[len(path)-1].ID() == rsg.graph.EndNode().ID() {
		cp := make([]model.PathPoint, len(path))
		copy(cp, path)
		*count++
		if !yield(cp, nil) || *count >= k {
			return false
		}
	}

	var candidates []model.PathPointID
	err := rsg.graph.Nodes().Range(func(nodeID model.PathPointID, node *model.PathPoint) error {
		nodeInDegree := tempInDegree[nodeID]
		nodeVisited := visited[nodeID]
		if !nodeVisited && nodeInDegree == 0 {
			candidates = append(candidates, nodeID)
		}
		return nil
	})
	if err != nil {
		log.Error().Msg("error iterating over nodes")
		return false
	}

	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for _, nodeID := range candidates {
		visited[nodeID] = true
		node, exists := rsg.graph.GetNode(nodeID)
		if !exists {
			log.Error().Msgf("node %s does not exist in the graph", nodeID)
			return false // Node does not exist in the graph
		}
		path = append(path, *node)
		for _, neigh := range rsg.graph.GetEdges(nodeID) {
			tempInDegree[neigh]--
		}
		if !rsg.randomBacktrack(tempInDegree, visited, path, yield, count, k, r) {
			return false // Stop if yield returns false or count exceeds k
		}
		for _, neigh := range rsg.graph.GetEdges(nodeID) {
			tempInDegree[neigh]++
		}
		path = path[:len(path)-1]
		visited[nodeID] = false
	}
	return true
}
