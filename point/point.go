// point.go
// 2016 giulio <giulioungaretti@me.com>

package point

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/giulioungaretti/geo/s2"

	log "github.com/Sirupsen/logrus"
)

// Point umarshal the json requests
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	// 30 is max precision i.e. cell leaf
	// default to nil, which is the max precision
	Precision *int `json:"precision"`
}

// S2Point is a point in s2 coordinates
type S2Point struct {
	CellID uint64 `json:"cellid"`
}

// IsValid vaildates point struct
func (p Point) IsValid() bool {
	if p.Lat == 0 || p.Lng == 0 {
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
func (p Point) ToS2() uint64 {
	ll := s2.LatLngFromDegrees(p.Lat, p.Lng)
	cellID := s2.CellIDFromLatLng(ll)
	// approximate returns
	if p.Precision != nil {
		return uint64(cellID.Parent(*p.Precision))
	}
	return uint64(cellID)
}

// ToCell returns s2 cell id from lat long precision
func (p Point) ToCell() (s2.LatLng, s2.Cell) {
	ll := s2.LatLngFromDegrees(p.Lat, p.Lng)
	cellID := s2.CellFromLatLng(ll)
	return ll, cellID
}

func parseJSON(r *http.Request) (Point, error) {
	var p Point
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	return p, err
}

//Handler handles a request for a geojsonPoint
func Handler(w http.ResponseWriter, r *http.Request) {
	// request
	p, err := parseJSON(r)
	if err != nil {
		log.Errorf("Malformed json in request:%v", r.Body)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !p.IsValid() {
		err := fmt.Errorf("Malformed request expected no zero lat long")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// response
	resp := S2Point{
		CellID: p.ToS2(),
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	err = encoder.Encode(resp)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
