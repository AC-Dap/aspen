package main

import (
	"fmt"
	"log"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	// TODO: Choose a more restrictive origin
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func testAPIHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println("Request received", r)
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/api/test", testAPIHandler)

	log.Fatal(http.ListenAndServe(":3001", nil))
}
