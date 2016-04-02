// Package counts handles a multipolygon requests
// and returns the event counts inside the polygon
// in the current time range.
package counts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/giulioungaretti/geo/s2"
	"github.com/giulioungaretti/globo/geoJSON"
	// use postgres driver
	_ "github.com/lib/pq"
)

const (
	dbUser     = "giulio"
	dbPassword = "postgres"
	dbName     = "postgres"
	// Endpoint is  the name of the counts endpoint
	Endpoint = "/v1/counts/"
)

var db *sql.DB
var err error

type row struct {
	value  uint64
	cellid uint64
}

func checker(c, t *uint64, done chan struct{}) {
	for {
		time.Sleep(10 * time.Microsecond)
		log.Debugf("⚡️  %v, %v", *c, *t)
		if *c == *t {
			log.Debugf("⚡️  done")
			close(done)
			return
		}
	}
}

func processor(ch chan row, count, processed *uint64, loop s2.Loop, bound s2.Rect) {
	for row := range ch {
		if containsCell(loop, bound, row.cellid) {
			atomic.AddUint64(count, row.value)
		}
		atomic.AddUint64(processed, 1)
	}
}

func init() {
	db, err = connect()
	if err != nil {
		log.Fatal(err)
	}
}

func connect() (db *sql.DB, err error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		return
	}
	err = db.Ping()
	if err != nil {
		return
	}
	return
}

func query(db *sql.DB, minID, maxID uint64, start string, end string) (rows *sql.Rows, err error) {
	rows, err = db.Query(`
	SELECT sum(count), s2cellid
	FROM data
	WHERE
		daterange($1::date, $2::date, '[]')
		@> day
	AND
		numrange($3, $4)
		@>s2cellid
	GROUP BY
		s2cellid
	`, start, end, minID, maxID)
	return
}

type res struct {
	loopid int
	res    uint64
}

//Handler takes care of the request
func Handler(w http.ResponseWriter, r *http.Request) {

	workers := runtime.NumCPU()
	// parsq query
	values := r.URL.Query()
	var start string
	var end string
	var precision int
	var err error
	if p, ok := values["precision"]; ok {
		precision, err = strconv.Atoi(p[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		// set max precision
		precision = 30
	}
	if starts, ok := values["start"]; ok {
		start = starts[0]
	} else {
		http.Error(w, "Missing start date", http.StatusBadRequest)
	}
	if ends, ok := values["end"]; ok {
		end = ends[0]
	} else {
		http.Error(w, "Missing end date", http.StatusBadRequest)
	}

	// match geoJSON type
	resp, err := Matcher(r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// get data from geojson
	cellUnions, loops, err := resp.ToS2(precision)
	return
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results := make(map[string]uint64)
	t := time.Now()
loop:
	for i, ids := range cellUnions {
		min := ids[0]
		max := ids[len(ids)-1]
		rows, err := query(db, min, max, start, end)
		if err != nil {
			log.Error(err)
			log.Debug(min, max, nil, start, end)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}
		var rowsnumber uint64
		// get loop
		loop := loops[i]
		// get bound
		bound := loop.RectBound()
		// channels
		var processed uint64
		var count uint64
		ch := make(chan row, 1000)
		done := make(chan struct{})

		go checker(&processed, &rowsnumber, done)
		for i := 0; i < workers; i++ {
			go processor(ch, &count, &processed, loop, bound)
		}

		for rows.Next() {
			var row row
			if err := rows.Scan(&row.value, &row.cellid); err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break loop
			} else {
				ch <- row
				atomic.AddUint64(&rowsnumber, 1)
			}
		}
		<-done
		log.Debugf("🚥")
		results[fmt.Sprint(i)] = count
		log.Debugf("total row recieved:%v", rowsnumber)
		log.Debugf("total  recieved:%v", results)
	}
	log.Debug(time.Since(t))
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	err = encoder.Encode(results)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func containsCell(l s2.Loop, bound s2.Rect, cellid uint64) (res bool) {
	// fast check: if cell not in bounding rect
	// return false
	id := s2.CellID(cellid)
	cell := s2.CellFromCellID(id)
	cellBound := cell.RectBound()
	if !bound.Contains(cellBound) {
		return false
	}
	ll := id.LatLng()
	if l.ContainsPoint(s2.PointFromLatLng(ll)) {
		return true
	}
	return res
}

// Matcher exctract from the url witch geoJSON object we want
func Matcher(r *http.Request) (p geoJSON.GeoJSON, err error) {
	objectType := r.URL.Path[len(Endpoint):]
	dec := json.NewDecoder(r.Body)
	switch objectType {
	case "point":
		// TODO this is ugly
		pp := geoJSON.Point{}
		err = dec.Decode(&pp)
		p = pp
	case "polygon":
		pp := geoJSON.Polygon{}
		err = dec.Decode(&pp)
		p = pp
	case "multipolygon":
		pp := geoJSON.MultiPolygon{}
		err = dec.Decode(&pp)
		p = pp
	default:
		err = fmt.Errorf("Bad geoJSON object type")
	}
	return p, err
}
