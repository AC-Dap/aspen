package main

import (
	"fmt"
	"log"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
}

func testAPIHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println("Request received", r)
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/api/test", testAPIHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
