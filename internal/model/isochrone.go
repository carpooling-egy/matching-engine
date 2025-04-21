package model

import "errors"

type Isochrone struct {
	contour  float64
	metric   string
	geometry GeoJSONGeom
}

func NewIsochrone(contour float64, metric string, geometry GeoJSONGeom) (*Isochrone, error) {
	if contour < 0 {
		return nil, errors.New("contour is negative")
	}

	if metric == "" {
		return nil, errors.New("metric is empty")
	}

	return &Isochrone{
		contour:  contour,
		metric:   metric,
		geometry: geometry,
	}, nil
}

func (I *Isochrone) Contour() (float64, error) {
	if I == nil {
		return 0, errors.New("nil isochrone reference")
	}

	return I.contour, nil
}

func (I *Isochrone) Metric() (string, error) {
	if I == nil {
		return "", errors.New("nil isochrone reference")
	}

	return I.metric, nil
}

func (I *Isochrone) Geometry() (GeoJSONGeom, error) {
	if I == nil {
		return GeoJSONGeom{}, errors.New("nil isochrone reference")
	}

	return I.geometry, nil
}

type GeoJSONGeom struct {
	geometryType string
	coordinates  [][][]float64
}

func NewGeoJSONGeom(geometryType string, coordinates [][][]float64) (*GeoJSONGeom, error) {
	if geometryType == "" {
		return nil, errors.New("geometry type is empty")
	}

	if len(coordinates) == 0 {
		return nil, errors.New("coordinates are empty")
	}

	return &GeoJSONGeom{
		geometryType: geometryType,
		coordinates:  coordinates,
	}, nil
}

func (g *GeoJSONGeom) GeometryType() (string, error) {
	if g == nil {
		return "", errors.New("nil GeoJSONGeom reference")
	}

	return g.geometryType, nil
}

func (g *GeoJSONGeom) Coordinates() ([][][]float64, error) {
	if g == nil {
		return nil, errors.New("nil GeoJSONGeom reference")
	}

	return g.coordinates, nil
}
