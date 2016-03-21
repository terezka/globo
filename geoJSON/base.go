// Package geoJSON converts geoJSON TO s2
// 2016 giulio <giulioungaretti@me.com>
package geoJSON

import (
	"fmt"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/giulioungaretti/geo/s2"
	"github.com/giulioungaretti/globo/point"
)

//Coordinate is a [longitude, latitude]
type Coordinate [2]float64

//Coordinates is a set of coordinate
type Coordinates []Coordinate

// Point rapresent a geoJSON point geometry object
type Point struct {
	Type       string `json:"type"`
	Coordinate `json:"coordinates"`
}

// Polygon rapresent a geoJSON polygon geometry object
type Polygon struct {
	Type string        `json:"type"`
	C    []Coordinates `json:"coordinates"`
}

// MultiPolygon rapresent a geoJSON mulitpolygon  geometry object
type MultiPolygon struct {
	Type string          `json:"type"`
	C    [][]Coordinates `json:"coordinates"`
}

// Prop is a geoJSON property
type Prop map[string]string

// Feature is a geoJSON feature
type Feature struct {
	Type       string      `json:"type"`
	Geometry   interface{} `json:"geometry"`
	Properties Prop        `json:"properties"`
}

//FeatureCollection is a geoJSON Feature colelction
type FeatureCollection struct {
	Type string    `json:"type"`
	Feat []Feature `json:"features"`
}

// geoJSON is the interface that allows any geojson to be unmarshaled
// converted to s2, and marhsaled back to be visualized
type geoJSON interface {
	ToS2(int) ([][]uint64, []s2.Loop, error)
	ToGeoJSON(int) (FeatureCollection, error)
}

func (c Coordinate) tos2point() s2.Point {
	ll := s2.LatLngFromDegrees(c[1], c[0])
	p := s2.PointFromLatLng(ll)
	return p
}

// IsValid vaildates coordinates (approximate)
func (c Coordinate) IsValid() bool {
	if c[0] == 0 || c[1] == 0 {
		return false
	}
	return true
}

// tos2 transforms Coordinates to a s2 loop and
// checks if the loop is CCW.
func (cc Coordinates) tos2() (s2.Loop, error) {
	pts := []s2.Point{}
	for _, c := range cc {
		p := c.tos2point()
		pts = append(pts, p)
	}
	//origin := s2.OriginPoint()
	//for i := range pts {
	//j := 1 + i
	//k := 2 + i
	//if i == len(pts)-2 {
	//k = 0
	//}
	//if i == len(pts)-1 {
	//j = 0
	//k = 1
	//}
	//if !s2.OrderedCCW(pts[i], pts[j], pts[k], origin) {
	//err := fmt.Errorf("Polygon not ordered")
	//log.Errorf("%v,%v,%v, %v", i, j, k, err)
	////return *s2.LoopFromPoints(pts), err
	//}
	//}
	return *s2.LoopFromPoints(pts), nil
}

// ToS2 converts a Point to S2Cellid.
// If precision is different than 30 (max) then the
// parent cellid at specified level is returned.
// levels go from 0 to 30:
// for reference :
// 30 covers 0.48cm2
// 12 covers 3.31km2
// 0 covers 85,011,012 km2
func (p Point) ToS2(precision int) (ids [][]uint64, loops []s2.Loop, err error) {
	c := p.Coordinate
	ll := s2.LatLngFromDegrees(c[1], c[0])
	cellID := s2.CellIDFromLatLng(ll)
	// approximate returns
	if precision != 30 {
		log.Debug("Rounding down")
		return [][]uint64{[]uint64{uint64(cellID.Parent(precision))}}, loops, nil
	}
	return [][]uint64{[]uint64{uint64(cellID)}}, loops, nil
}

// ToGeoJSON converts point to geoJSON
func (p Point) ToGeoJSON(precision int) (ff FeatureCollection, err error) {
	// ingore loops, as a point has obviously no loops
	in, _, err := p.ToS2(precision)
	if err != nil {
		return
	}
	ff = FeatureCollection{}
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
	return
}

// innerLoop(s) return an s2.loop representation of the inner loop of the
// geoJSON polygon or multi polygon. We do not support holes in polygons, so the inner ring of the geojson is discarded.
// At this point winding order of the loops is not specified but it **must** be
// counterclockwise, else we return an error.

func (p Polygon) innerLoops() (loops []s2.Loop, err error) {
	cds := p.C[0]
	loop, err := cds.tos2()
	loops = append(loops, loop)
	return loops, err
}

func (pp MultiPolygon) innerLoops() ([]s2.Loop, error) {
	var loops []s2.Loop
	var err error
	for _, c := range pp.C {
		cds := c[0]
		loop, err := cds.tos2()
		if err != nil {
			return loops, err
		}
		loops = append(loops, loop)
	}
	return loops, err
}

// Contains check if the polygon contains the  point
// TODO needs polish
func (p Polygon) Contains(point point.Point) bool {
	// TODO this has a bug
	// maybe is the point to s2 conversion
	var contains bool
	ll, cell := point.ToCell()
	loop, err := p.innerLoops()
	if err != nil {
		log.Error(err)
		return false
	}

	contains = loop[0].RectBound().ContainsCell(cell)
	if !contains {
		return false
	}
	contains = loop[0].ContainsPoint(s2.PointFromLatLng(ll))
	return contains
}

func cellIDToPolygon(id uint64) (f Feature) {
	cellid := s2.CellID(id)
	cell := s2.CellFromCellID(cellid)
	rect := cell.RectBound()
	var coordinates Coordinates
	var coordinates2 []Coordinates
	for i := 0; i < 4; i++ {
		ll := rect.Vertex(i)
		// ll is the vertex of the  cellid
		var ld Coordinate
		ld = [2]float64{ll.Lng.Degrees(), ll.Lat.Degrees()}
		coordinates = append(coordinates, ld)
	}
	// add first point as last point to close polygon
	coordinates = append(coordinates, coordinates[0])
	coordinates2 = append(coordinates2, coordinates)
	polygon := Polygon{}
	polygon.Type = "Polygon"
	polygon.C = coordinates2
	f.Type = "Feature"
	prop := make(Prop)
	prop["cellid"] = fmt.Sprintf("%v", cellid)
	f.Properties = prop
	f.Geometry = polygon
	return f
}

func cellIDToCenterPoint(id uint64) (f Feature) {
	cellid := s2.CellID(id)
	// ll is the center of the  cellid
	ll := cellid.LatLng()
	var ld Coordinate
	ld = [2]float64{ll.Lng.Degrees(), ll.Lat.Degrees()}
	point := Point{}
	point.Type = "Point"
	point.Coordinate = ld
	f.Type = "Feature"
	prop := make(Prop)
	prop["is_center"] = "true"
	f.Properties = prop
	f.Geometry = point
	return f
}

func boundingbox(loop s2.Loop, id int) (f Feature) {
	rect := loop.RectBound()
	var coordinates Coordinates
	var coordinates2 []Coordinates
	for i := 0; i < 4; i++ {
		ll := rect.Vertex(i)
		// ll is the vertex of the  cellid
		var ld Coordinate
		ld = [2]float64{ll.Lng.Degrees(), ll.Lat.Degrees()}
		coordinates = append(coordinates, ld)
	}
	// add last vertex
	coordinates = append(coordinates, coordinates[0])
	coordinates2 = append(coordinates2, coordinates)
	polygon := Polygon{}
	polygon.Type = "Polygon"
	polygon.C = coordinates2
	f = Feature{}
	f.Type = "Feature"
	prop := make(Prop)
	prop["boundingbox"] = fmt.Sprint(id)
	f.Properties = prop
	f.Geometry = polygon
	return
}

// transforms a loop to a sorted cellID collection
func loopCoverer(loop s2.Loop, precision int) []uint64 {
	debug := false
	// now for bounds
	rc := &s2.RegionCoverer{MinLevel: 0, MaxLevel: precision, MaxCells: 100}
	r := s2.Region(loop.RectBound())
	covering := rc.Covering(r)
	// now approximate the bounding rect
	var polygon []uint64
	for _, val := range covering {
		polygon = append(polygon, uint64(val))
	}
	// NOTE ++ TODO
	// is this 100% correct ?
	// in theory it is.
	sort.Sort(intArray(polygon))
	// NOTE this will be the actual polygonsponse
	// when I/google finish updating the s2 library
	// right now the region coverer for loops
	// is buggy in the sense that it reutns
	// just the boundary of the loop.
	// by setting debug to true, one can
	// visualize the json
	if debug {
		min := polygon[0]
		max := polygon[len(polygon)-1]
		rc := &s2.RegionCoverer{MinLevel: 15, MaxLevel: 30, MaxCells: 50}
		covering := rc.Covering(s2.Region(loop))
		var boundary []uint64
		// now approximate the polygon
		for _, val := range covering {
			boundary = append(boundary, uint64(val))
		}
		for _, val := range boundary {
			if val > max {
				log.Printf("strange %v", val)
			}
			if val < min {
				log.Printf("strange %v", val)
			}
		}
		return boundary
	}
	return polygon
}

// ToS2 converts a geoJSON polygon to a set of cellUnions
func (p Polygon) ToS2(precision int) (ids [][]uint64, loops []s2.Loop, err error) {
	var polygons [][]uint64
	loops, err = p.innerLoops()
	if err != nil {
		return polygons, loops, err
	}
	polygon := loopCoverer(loops[0], precision)
	polygons = append(polygons, polygon)
	return polygons, loops, err
}

// ToS2 converts a geoJSON multi polygon to a set of cellUnions
func (pp MultiPolygon) ToS2(precision int) (ids [][]uint64, loops []s2.Loop, err error) {
	var polygons [][]uint64
	loops, err = pp.innerLoops()
	if err != nil {
		return polygons, loops, err
	}
	for _, loop := range loops {
		polygon := loopCoverer(loop, precision)
		polygons = append(polygons, polygon)
	}
	return polygons, loops, err
}

// ToGeoJSON converts polygon to geoJSON
func (p Polygon) ToGeoJSON(precision int) (ff FeatureCollection, err error) {
	in, loops, err := p.ToS2(precision)
	if err != nil {
		return
	}
	// NOTE  we expect one polygon/ one loop
	polygon := in[0]
	inner := loops[0]
	ff = FeatureCollection{}
	ff.Type = "FeatureCollection"
	var features []Feature
	// add bbox
	for _, id := range polygon {
		feature := cellIDToPolygon(id)
		features = append(features, feature)
	}
	bbox := boundingbox(inner, 0)
	features = append(features, bbox)
	ff.Feat = features
	return
}

// ToGeoJSON converts back a multipolygon to geoJSON
func (pp MultiPolygon) ToGeoJSON(precision int) (ff FeatureCollection, err error) {
	in, loops, err := pp.ToS2(precision)
	if err != nil {
		return
	}
	ff = FeatureCollection{}
	ff.Type = "FeatureCollection"
	var features []Feature
	for _, polygon := range in {
		for _, id := range polygon {
			feature := cellIDToPolygon(id)
			features = append(features, feature)
		}
	}
	// add bbox
	for i, inner := range loops {
		bbox := boundingbox(inner, i)
		features = append(features, bbox)
	}
	ff.Feat = features
	return
}

// HELPERS
// sortable array
type intArray []uint64

func (s intArray) Len() int           { return len(s) }
func (s intArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s intArray) Less(i, j int) bool { return s[i] < s[j] }

// convert to string token
func toToken(ci uint64) string {
	s := strings.TrimRight(fmt.Sprintf("%016x", uint64(ci)), "0")
	if len(s) == 0 {
		return "X"
	}
	return s
}
