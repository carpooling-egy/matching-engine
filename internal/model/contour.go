package model

import (
	"errors"
	"fmt"
	"strings"
)

type ContourMetric string

const (
	ContourMetricTimeInMinutes        ContourMetric = "time"
	ContourMetricDistanceInKilometers ContourMetric = "distance"
)

func NewContourMetric(raw string) (ContourMetric, error) {
	metric := ContourMetric(strings.ToLower(strings.TrimSpace(raw)))
	if !metric.isValid() {
		return "", fmt.Errorf("invalid ContourMetric %q", raw)
	}
	return metric, nil
}

func (m ContourMetric) isValid() bool {
	switch m {
	case ContourMetricTimeInMinutes, ContourMetricDistanceInKilometers:
		return true
	}
	return false
}

func (m ContourMetric) String() string {
	return string(m)
}

func (m ContourMetric) Unit() string {
	switch m {
	case ContourMetricTimeInMinutes:
		return "minutes"
	case ContourMetricDistanceInKilometers:
		return "kilometers"
	default:
		return ""
	}
}

type Contour struct {
	value  float32
	metric ContourMetric
}

func NewContour(value float32, metric ContourMetric) (*Contour, error) {
	if value < 0 {
		return nil, errors.New("contour value must be non-negative")
	}

	if !metric.isValid() {
		return nil, fmt.Errorf("invalid contour metric %q", metric)
	}

	return &Contour{
		value:  value,
		metric: metric,
	}, nil
}

func (c *Contour) Value() float32 {
	return c.value
}

func (c *Contour) Metric() ContourMetric {
	return c.metric
}

func (c *Contour) String() string {
	return fmt.Sprintf("%.2f %s", c.value, c.metric)
}
