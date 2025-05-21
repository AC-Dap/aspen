package main

import (
	"aspen/resources"
	"aspen/router"
	"log"
	"net/http"
)

func main() {
	// Init resources
	var resources = map[string]resources.Resource{
		"/info":   resources.NewStaticFile("info", "README.md"),
		"/design": resources.NewStaticFile("design", "design.md"),
		"/code":   resources.NewStaticDirectory("resources", "resources", []string{"static_file.go"}, false),
	}

	// Init router
	var router router.Router
	router.Init(resources)

	// Start server
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", &router))
}
