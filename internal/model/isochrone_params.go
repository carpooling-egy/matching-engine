package model

import (
	"errors"
	"fmt"
)

type IsochroneParams struct {
	origin  *Coordinate
	contour *Contour
	profile RoutingProfile
}

func NewIsochroneParams(
	origin *Coordinate,
	contour *Contour,
	profile RoutingProfile,
) (*IsochroneParams, error) {
	if origin == nil {
		return nil, errors.New("origin is nil")
	}

	if contour == nil {
		return nil, errors.New("contour is nil")
	}

	if profile == "" {
		return nil, errors.New("profile is empty")
	}

	return &IsochroneParams{
		origin:  origin,
		contour: contour,
		profile: profile,
	}, nil
}

func (ip *IsochroneParams) Origin() *Coordinate {
	return ip.origin
}

func (ip *IsochroneParams) Contour() *Contour {
	return ip.contour
}

func (ip *IsochroneParams) Profile() RoutingProfile {
	return ip.profile
}

func (ip *IsochroneParams) String() string {
	return fmt.Sprintf(
		"IsochroneParams{origin: %s, contour: %s, profile: %s}",
		ip.origin.String(), ip.contour.String(), ip.profile.String(),
	)
}
