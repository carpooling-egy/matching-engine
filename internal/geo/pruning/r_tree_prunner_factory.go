package pruning

import (
	"errors"
	"github.com/dhconnelly/rtreego"
	"matching-engine/internal/model"
)

// R-tree configuration constants
const (
	TreeDimension                = 2
	MinNodeEntries               = 25
	MaxNodeEntries               = 50
	MinThresholdDistanceInMeters = 50
)

// RTreePrunerFactory implements the RoutePrunerFactory interface
type RTreePrunerFactory struct{}

// NewRTreePrunerFactory creates a new factory instance
func NewRTreePrunerFactory() RoutePrunerFactory {
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
	indexMap := make(map[string]int, len(route))

	for i := 0; i < len(route)-1; i++ {
		a, b := route[i], route[i+1]

		seg := NewSegment(a, b)
		tree.Insert(seg)

		ka, kb := a.Key(), b.Key()
		if _, ok := indexMap[ka]; !ok {
			indexMap[ka] = i
		}
		if _, ok := indexMap[kb]; !ok {
			indexMap[kb] = i + 1
		}
	}

	return &RTreePruner{
		tree:     tree,
		indexMap: indexMap,
	}, nil
}
