package main

import (
	"dashboard/config"
	"dashboard/integrations"
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

func testAuthHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println("Request received", r)
	fmt.Fprintf(w, "Authorized!")
}

func main() {
	http.HandleFunc("/api/test", testAPIHandler)
	http.HandleFunc("/auth", testAuthHandler)

	resources := config.Read()
	for _, resource := range resources.Resources {
		log.Println(resource.Name)
		log.Println("  Route:", resource.Route)
		log.Println("  Source:", resource.Source)
		log.Println("  Restricted:", resource.Restricted)
	}

	config.Validate(resources)

	nginx_conf := integrations.GenerateNginxConfig("8080", resources.Resources)
	log.Println(nginx_conf)
	integrations.ReloadNginxConfig("nginx.conf", nginx_conf)

	log.Fatal(http.ListenAndServe(":3001", nil))
}
