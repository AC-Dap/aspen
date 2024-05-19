package main

import (
	"dashboard/config"
	"dashboard/routing"
	"dashboard/types"
	"dashboard/util"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

var ProxyServer struct {
	router     *routing.Server
	routerAddr string
	routerLock sync.RWMutex
}

func UpdateRouter(resources []types.Resource) error {
	// Start new routing server
	newRouter := routing.NewServer(resources)
	addr, err := newRouter.Start()
	if err != nil {
		log.Println("Error starting server:", err)
		return err
	}

	// Update proxy server with new router
	ProxyServer.routerLock.Lock()
	oldRouter := ProxyServer.router
	ProxyServer.router = newRouter
	ProxyServer.routerAddr = addr
	ProxyServer.routerLock.Unlock()

	// Stop old routing server
	if oldRouter != nil {
		oldRouter.Stop()
	}

	return nil
}

func ProxyToRouter(w http.ResponseWriter, r *http.Request) {
	ProxyServer.routerLock.RLock()
	router := ProxyServer.router
	ProxyServer.routerLock.RUnlock()

	util.Assert(router != nil, "No routing server available")

	// Create a new request with the same method and body
	r.URL.Scheme = "http"
	r.URL.Host = ProxyServer.routerAddr
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		errMsg := fmt.Sprintln("Unable to create proxy request:", err)
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Copy headers
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errMsg := fmt.Sprintln("Unable to send proxy request:", err)
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy request to w
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

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
	// Read resources from config
	resources := config.Read()
	for _, resource := range resources.Resources {
		log.Println(resource.Name)
		log.Println("  Route:", resource.Route)
		log.Println("  Source:", resource.Source)
		log.Println("  Restricted:", resource.Restricted)
	}

	config.Validate(resources)

	// Start routing server
	log.Println("Starting routing server")
	err := UpdateRouter(resources.Resources)
	if err != nil {
		log.Fatal(err)
	}

	// Start proxy server
	log.Println("Starting proxy server")
	server := http.Server{
		Addr:    ":3001",
		Handler: http.HandlerFunc(ProxyToRouter),
	}
	log.Fatal(server.ListenAndServe())

	//nginxConf := integrations.GenerateNginxConfig("8080", resources.Resources)
	//log.Println(nginxConf)
	//integrations.ReloadNginxConfig("nginx.conf", nginxConf)
}
