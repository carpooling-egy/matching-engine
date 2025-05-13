package pruning
import (
	"math"
	"matching-engine/internal/model"
	"github.com/dhconnelly/rtreego"
)
type Segment struct {
	A, B model.Coordinate
	rect rtreego.Rect
}

// Implements the rtreego.Spatial interface
func (s *Segment) Bounds() rtreego.Rect {
	return s.rect
}

func NewSegment(a, b model.Coordinate) *Segment {
	minLat, maxLat := math.Min(a.Lat(), b.Lat()), math.Max(a.Lat(), b.Lat())
	minLng, maxLng := math.Min(a.Lng(), b.Lng()), math.Max(a.Lng(), b.Lng())

	rect, _ := rtreego.NewRect(
		rtreego.Point{minLat, minLng},
		[]float64{maxLat - minLat, maxLng - minLng},
	)

	return &Segment{A: a, B: b, rect: rect}
	
}
