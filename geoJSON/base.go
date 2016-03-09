// Package geoJSON converts geoJSON TO s2
// 2016 giulio <giulioungaretti@me.com>
package geoJSON

import "github.com/golang/geo/s2"

//Coordinates is a set of coordinate
type Coordinates []Coordinate

//Coordinate is a [longitude, latitude]
type Coordinate [2]float64

// Point rapresent a geojson point geometry object
type Point struct {
	Type       string `json:"type"`
	Coordinate `json:"coordinates"`
}

// IsValid vaildates coordinate
func (c *Coordinate) IsValid() bool {
	if c[0] == 0 || c[1] == 0 {
		return false
	}
	return true
}

// ToS2 converts lat/long to S2Cellid. If p.Precision is specified then the
// parent cellid at specfied level is returned.
// levels go from 0 to 30:
// for reference :
// 0 covers 0.48cm2
// 12 covers 3.31km2
func (p Point) ToS2(precision *int) uint64 {
	c := p.Coordinate
	ll := s2.LatLngFromDegrees(c[1], c[0])
	cellID := s2.CellIDFromLatLng(ll)
	// approximate returns
	if precision != nil {
		return uint64(cellID.Parent(*precision))
	}
	return uint64(cellID)
}

// Polygon rapresent a geojson polygon geometry object
type Polygon struct {
	Type string        `json:"type"`
	C    []Coordinates `json:"coordinates"`
}

// ToS2 converts lat/long to S2Cellid.
func (p Polygon) ToS2(precision *int) uint64 {
	return uint64(1)
}

// MultiPolygon rapresent a geojson mulitpolygon  geometry object
type MultiPolygon struct {
	Type string          `json:"type"`
	C    [][]Coordinates `json:"coordinates"`
}

// ToS2 converts lat/long to S2Cellid.
func (p MultiPolygon) ToS2(precision *int) uint64 {
	return uint64(1)
}

type geojson interface {
	ToS2(*int) uint64
}
