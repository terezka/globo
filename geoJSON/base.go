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

//Coordinates is a set of coordinate
type Coordinates []Coordinate

func (cc Coordinates) tos2() (s2.Loop, error) {
	pts := []s2.Point{}
	for _, c := range cc {
		p := c.tos2point()
		pts = append(pts, p)
	}
	origin := s2.OriginPoint()
	for i := range pts {
		j := 1 + i
		k := 2 + i
		if i == len(pts)-2 {
			k = 0
		}
		if i == len(pts)-1 {
			j = 0
			k = 1
		}
		if !s2.OrderedCCW(pts[i], pts[j], pts[k], origin) {
			err := fmt.Errorf("Polygon not ordered")
			return *s2.LoopFromPoints(pts), err
		}
	}
	return *s2.LoopFromPoints(pts), nil
}

// Point rapresent a geoJSON point geometry object
type Point struct {
	Type       string `json:"type"`
	Coordinate `json:"coordinates"`
}

// ToS2 converts lat/long to S2Cellid. If p.Precision is specified then the
// parent cellid at specfied level is returned.
// levels go from 0 to 30:
// for reference :
// 30 covers 0.48cm2
// 12 covers 3.31km2
// 0 covers 85,011,012 km2
func (p Point) ToS2(precision int) [][]uint64 {
	c := p.Coordinate
	ll := s2.LatLngFromDegrees(c[1], c[0])
	cellID := s2.CellIDFromLatLng(ll)
	// approximate returns
	if precision != 30 {
		log.Print("Rounding down")
		return [][]uint64{[]uint64{uint64(cellID.Parent(precision))}}
	}
	return [][]uint64{[]uint64{uint64(cellID)}}
}

// ToGeoJSON converts point to geoJSON
func (p Point) ToGeoJSON(in [][]uint64) FeatureCollection {
	ff := FeatureCollection{}
	ff.Type = "FeatureCollection"
	var features []Feature
	for _, id := range in[0] {
		// add center point
		feature := cellIDToCenterPoint(id)
		features = append(features, feature)
		// add bounding box
		feature = cellIDToPolygon(id)
		features = append(features, feature)
	}
	ff.Feat = features
	return ff
}
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
