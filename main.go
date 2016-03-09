// GLOBO is a microservice that converts lat/long to s2
// main.go
// 2016 giulio <giulioungaretti@me.com>

package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/giulioungaretti/globo/geoJSON"
	"github.com/giulioungaretti/globo/middleware"
	"github.com/giulioungaretti/globo/point"
	"github.com/justinas/alice"
)

func main() {
	defaultMiddleware := alice.New(middleware.TimeOut, middleware.IsJSON, middleware.PostOnly)

	// point
	pointHandler := http.HandlerFunc(point.Handler)
	http.Handle("/tos2/point", defaultMiddleware.Then(pointHandler))

	// geojson
	geoJSONHandler := http.HandlerFunc(geoJSON.Handl, r)
	http.Handle(geoJSON.Endpoint, defaultMiddleware.Then(geoJSONHandler))

	// server
	port := os.Getenv("PORT")
	lvl, err := ParseLogLevel(os.Getenv("LOGLEVEL"))
	if err != nil {
		panic(err)
	}
	log.SetLevel(lvl)
	addr := fmt.Sprintf(":%v", port)
	http.ListenAndServe(addr, nil)
}
