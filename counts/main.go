// Package counts handles a multipolygon requests
// and returns the event counts inside the polygon
// in the current time range.
package counts

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	// TODO
	// use new library asap
	// pg/pg
	// use postgres driver
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "s2"
)

func connect() (db *sql.DB, err error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
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

// NOTE should use
// date range for dates
// and int8range for ID
func query(db *sql.DB, minID int, maxID int, start time.Time, end time.Time) (rows *sql.Rows, err error) {
	rows, err = db.Query(`
	SELECT count, day FROM data
	WHERE
	day > '2010-08-08'::date`)
	return
}

// NOTE should use
// date range for dates
// and int8range for ID
//func query(db *sql.DB, minID int, maxID int, start time.Time, end time.Time) (rows *sql.Rows, err error) {
//rows, err = db.Query(`
//SELECT count FROM data
//WHERE
//day <= $1
//AND
//day >= $2
//AND
//s2cellid <= $3
//AND
//s2cellid >= $4
//`, end, start,  maxID, minID)
//return
//}
//Handler takes care of the request
func Handler(w http.ResponseWriter, r *http.Request) {
}
