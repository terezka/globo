// GLOBO is a microservice that converts lat/long to s2
// main.go
// 2016 giulio <giulioungaretti@me.com>

package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/tos2/point", point)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%v", port)
	http.ListenAndServe(addr, nil)
}
