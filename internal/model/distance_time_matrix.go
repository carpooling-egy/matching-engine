package model

import (
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"
)

type DistanceTimeMatrix struct {
	distances [][]Distance
	times     [][]time.Duration
}

func NewDistanceTimeMatrix(
	distances [][]Distance,
	times [][]time.Duration,
) (*DistanceTimeMatrix, error) {

	if err := validateMatrixDimensions(distances, times); err != nil {
		return nil, err
	}

	if err := validateSquareMatrix(distances); err != nil {
		return nil, err
	}
	if err := validateSquareMatrix(times); err != nil {
		return nil, err
	}

	if err := validateTimeMatrixElements(times); err != nil {
		return nil, err
	}

	return &DistanceTimeMatrix{
		distances: distances,
		times:     times,
	}, nil
}

func (dtm *DistanceTimeMatrix) Distances() [][]Distance {
	return dtm.distances
}

func (dtm *DistanceTimeMatrix) Times() [][]time.Duration {
	return dtm.times
}

func validateSquareMatrix[T any](matrix [][]T) error {
	rows := len(matrix)
	for i := range matrix {
		if len(matrix[i]) != rows {
			return fmt.Errorf("matrix must be square, "+
				"row %d has an incorrect number of columns", i)
		}
	}
	return nil
}

func validateMatrixDimensions(distances [][]Distance, times [][]time.Duration) error {
	if len(distances) == 0 || len(times) == 0 {
		return errors.New("distances or times matrix is empty")
	}

	if len(distances) != len(times) {
		return errors.New("distances and times matrix must have the same number of rows")
	}

	for i := range distances {
		if len(distances[i]) != len(times[i]) {
			return fmt.Errorf("distances and times matrices must have the same number of columns at row %d", i)
		}
	}
	return nil
}

func validateTimeMatrixElements(times [][]time.Duration) error {
	for i := range times {
		for j := range times[i] {
			if times[i][j] < 0 {
				return fmt.Errorf("time at [%d][%d] cannot be negative", i, j)
			}
		}
	}
	return nil
}

func (dtm *DistanceTimeMatrix) String() string {
	var builder strings.Builder
	tw := tabwriter.NewWriter(&builder, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprint(tw, "\t")
	for col := range dtm.distances[0] {
		_, _ = fmt.Fprintf(tw, "To %d\t", col)
	}
	_, _ = fmt.Fprintln(tw)

	for rowIndex, row := range dtm.distances {
		_, _ = fmt.Fprintf(tw, "From %d\t", rowIndex)
		for colIndex, distance := range row {
			duration := dtm.times[rowIndex][colIndex]
			entry := fmt.Sprintf("%s (%s)", distance.String(), duration)
			_, _ = fmt.Fprintf(tw, "%s\t", entry)
		}
		_, _ = fmt.Fprintln(tw)
	}

	if err := tw.Flush(); err != nil {
		return fmt.Sprintf("error flushing table: %v", err)
	}

	return builder.String()
}
