// Package geoJSON converts geoJSON TO s2
package geoJSON

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// Endpoint is  the name of the geojson handler endpoint
const Endpoint = "/tos2/geojson/"

//Handler handles a request for a geojsonPoint
func Handler(w http.ResponseWriter, r *http.Request) {
	// request
	resp, err := matcher(r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO extract from query
	var prec int
	s2 := resp.ToS2(&prec)
	// response
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	err = encoder.Encode(s2)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Matcher exctract from the url witch geoJSON object we want
func matcher(r *http.Request) (p geojson, err error) {
	objectType := r.URL.Path[len(Endpoint):]
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	switch objectType {
	case "point":
		// TODO this is ugly
		pp := Point{}
		err = dec.Decode(&pp)
		p = pp
	case "polygon":
		p = Polygon{}
	case "multipolygon":
		p = MultiPolygon{}
	default:
		err = fmt.Errorf("Bad geoJSON object type")
	}
	return p, err
}
