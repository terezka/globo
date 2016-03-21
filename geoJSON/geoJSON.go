// Package geoJSON converts geoJSON TO s2 and back to GeoJSON
// this is mostly an endpoint to visualize the simplifications
package geoJSON

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Endpoint is  the name of the geojson handler endpoint
const Endpoint = "/tos2/geojson/"

//Handler handles a request for a geojsonPoint
func Handler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	// parse form
	var precision int
	var err error
	values := r.URL.Query()
	if p, ok := values["precision"]; ok {
		precision, err = strconv.Atoi(p[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		// set max precision
		precision = 30
	}
	log.Debugf("Request with precision: %v", precision)
	// request
	resp, err := Matcher(r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	geoj, err := resp.ToGeoJSON(precision)
	// response
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	err = encoder.Encode(geoj)
	log.Debug(time.Since(t))
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Matcher exctract from the url witch geoJSON object we want
func Matcher(r *http.Request) (p GeoJSON, err error) {
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
		pp := Polygon{}
		err = dec.Decode(&pp)
		p = pp
	case "multipolygon":
		pp := MultiPolygon{}
		err = dec.Decode(&pp)
		p = pp
	default:
		err = fmt.Errorf("Bad geoJSON object type")
	}
	return p, err
}
