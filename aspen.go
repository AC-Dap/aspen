package main

import (
	"aspen/resources"
	"aspen/router"
	"log"
	"net/http"
)

func main() {
	// Init resources
	var resources = map[string]router.Resource{
		"/info":   resources.NewStaticFile("info", "README.md"),
		"/design": resources.NewStaticFile("design", "design.md"),
		"/code":   resources.NewStaticDirectory("resources", "resources", []string{"static_file.go"}, false),
		"/update": resources.NewUpdateRouterResource("update"),
	}

	// Init router
	router.UpdateRouter(resources)

	// Start server
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", &router.GlobalRouter))
}
