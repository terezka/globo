//
// middleware.go
// 2016 giulio <giulioungaretti@me.com>
//

package middleware

import (
	"net/http"
	"time"
)

// TimeOut handles timeout if request takes more than timeout
func TimeOut(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 1*time.Second, "timed out")
}

// IsJSON make sure that we get the request as tpye json
func IsJSON(next http.Handler) http.Handler {
	hanlderfunc := func(w http.ResponseWriter, r *http.Request) {
		if value := r.Header.Get("Content-Type"); value != "application/json" {
			http.Error(w, http.StatusText(415), 415)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hanlderfunc)
}

// PostOnly checks if the request tpye is post
func PostOnly(next http.Handler) http.Handler {
	handlerfunc := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handlerfunc)
}

// placeHolder is a middle ware that takes
// care of unmarshaling and validating the
// incoming json
func placeHolder(next http.Handler) http.Handler {
	hanlderfunc := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hanlderfunc)
}
