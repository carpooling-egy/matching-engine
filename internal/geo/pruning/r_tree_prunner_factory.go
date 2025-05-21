package pruning

import (
	"errors"
	"github.com/dhconnelly/rtreego"
	"matching-engine/internal/model"
)

// R-tree configuration constants
const (
	TreeDimension  = 2
	MinNodeEntries = 25
	MaxNodeEntries = 50
)

// RTreePrunerFactory implements the RoutePrunerFactory interface
type RTreePrunerFactory struct{}

// CreateRTreePrunerFactory creates a new factory instance
func CreateRTreePrunerFactory() RoutePrunerFactory {
	return &RTreePrunerFactory{}
}

// NewRoutePruner creates a new RoutePruner for the given route
func (f *RTreePrunerFactory) NewRoutePruner(route model.LineString) (RoutePruner, error) {
	return NewRTreePruner(route)
}

// NewRTreePruner creates a new RTreePruner with the specified route
func NewRTreePruner(route model.LineString) (*RTreePruner, error) {
	if len(route) == 0 {
		return nil, errors.New("route cannot be empty")
	}

	tree := rtreego.NewTree(TreeDimension, MinNodeEntries, MaxNodeEntries)

	// Index all segments from the route
	for i := 0; i < len(route)-1; i++ {
		seg := NewSegment(route[i], route[i+1])
		tree.Insert(seg)
	}

	return &RTreePruner{
		tree: tree,
	}, nil
}
