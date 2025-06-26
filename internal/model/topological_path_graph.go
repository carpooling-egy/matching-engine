package model

import (
	"matching-engine/internal/collections"
	"matching-engine/internal/enums"
)

type TopologicalPathGraph struct {
	// TopologicalPathGraph represents a graph structure for path generation
	nodes     *collections.SyncMap[PathPointID, *PathPoint]
	adjList   *collections.SyncMap[PathPointID, []PathPointID]
	inDegress *collections.SyncMap[PathPointID, int]
	startNode *PathPoint
	endNode   *PathPoint
}

func NewTopologicalPathGraph() *TopologicalPathGraph {
	return &TopologicalPathGraph{
		nodes:     collections.NewSyncMap[PathPointID, *PathPoint](),
		adjList:   collections.NewSyncMap[PathPointID, []PathPointID](),
		inDegress: collections.NewSyncMap[PathPointID, int](),
	}
}

// AddNode adds a node to the path graph
func (g *TopologicalPathGraph) AddNode(node *PathPoint) {
	g.nodes.Set(node.ID(), node)
	g.adjList.Set(node.ID(), []PathPointID{})
	g.inDegress.Set(node.ID(), 0)
}

func (g *TopologicalPathGraph) GetNode(nodeID PathPointID) (*PathPoint, bool) {
	// GetNode retrieves a node by its ID from the path graph
	node, exists := g.nodes.Get(nodeID)
	return node, exists
}

func (g *TopologicalPathGraph) StartNode() *PathPoint {
	// StartNode returns the start node of the path graph
	return g.startNode
}

func (g *TopologicalPathGraph) SetStartNode(node *PathPoint) {
	// SetStartNode sets the start node of the path graph
	g.startNode = node
	g.AddNode(node) // Ensure the start node is added to the graph
}

func (g *TopologicalPathGraph) EndNode() *PathPoint {
	// EndNode returns the end node of the path graph
	return g.endNode
}

func (g *TopologicalPathGraph) SetEndNode(node *PathPoint) {
	// SetEndNode sets the end node of the path graph
	g.endNode = node
	g.AddNode(node) // Ensure the end node is added to the graph
}

func (g *TopologicalPathGraph) Nodes() *collections.SyncMap[PathPointID, *PathPoint] {
	// GetNodes returns all nodes in the path graph
	return g.nodes
}

// AddEdge adds a directed edge from source to target in the path graph
func (g *TopologicalPathGraph) AddEdge(source, target PathPointID) {
	if _, exists := g.nodes.Get(source); !exists {
		return // Source node does not exist
	}
	if _, exists := g.nodes.Get(target); !exists {
		return // Target node does not exist
	}

	// Add the edge to the adjacency list
	edges, _ := g.adjList.Get(source)
	g.adjList.Set(source, append(edges, target))

	// Increment the in-degree of the target node
	inDegree, _ := g.inDegress.Get(target)
	g.inDegress.Set(target, inDegree+1)
}

// GetEdges returns the edges for a given node
func (g *TopologicalPathGraph) GetEdges(nodeID PathPointID) []PathPointID {
	if edges, exists := g.adjList.Get(nodeID); exists {
		return edges
	}
	return nil // No edges found for the node
}

// GetInDegree returns the in-degree of a given node
func (g *TopologicalPathGraph) GetInDegree(nodeID PathPointID) int {
	if inDegree, exists := g.inDegress.Get(nodeID); exists {
		return inDegree
	}
	return 0 // The Node does not exist or has no in-degree
}

func (g *TopologicalPathGraph) SetInDegree(nodeID PathPointID, inDegree int) {
	// SetInDegree sets the in-degree of a given node
	g.inDegress.Set(nodeID, inDegree)
}

// Clear clears the path graph
func (g *TopologicalPathGraph) Clear() {
	g.nodes.Clear()
	g.adjList.Clear()
	g.inDegress.Clear()
	g.startNode = nil
	g.endNode = nil
}

// Size returns the number of nodes in the path graph
func (g *TopologicalPathGraph) Size() int64 {
	return g.nodes.Size()
}

func (g *TopologicalPathGraph) CopyInDegree() map[PathPointID]int {
	tempInDegree := make(map[PathPointID]int)
	err := g.inDegress.Range(func(key PathPointID, value int) error {
		tempInDegree[key] = value
		return nil
	})
	if err != nil {
		return nil
	}
	return tempInDegree
}

func (g *TopologicalPathGraph) InitPathGraph(
	path []PathPoint,
	pickup, dropoff *PathPoint,
	startNode *PathPoint,
	endNode *PathPoint,
) {

	g.Clear()

	// Add the start and end nodes to the graph that guarantees that any order of the path can be generated
	// will start with the startNode and end with the endNode
	g.SetStartNode(startNode)
	g.SetEndNode(endNode)

	// Add the pickup and dropoff points to the graph
	g.AddNode(pickup)
	g.AddNode(dropoff)
	g.AddEdge(pickup.ID(), dropoff.ID())

	requestToPickup := collections.NewSyncMap[string, PathPointID]()
	pathLength := len(path)
	for i := 0; i < pathLength; i++ {
		current := path[i]
		g.AddNode(&current)
		if current.PointType() == enums.Pickup {
			requestToPickup.Set(current.GetOwnerID(), current.ID())
		} else if current.PointType() == enums.Dropoff {
			if pickupID, exists := requestToPickup.Get(current.GetOwnerID()); exists {
				g.AddEdge(pickupID, current.ID())
			} else {
				// If no pickup point exists for this dropoff, we can skip adding the edge
				continue
			}
		}
	}
}
