package routing

import (
	"dashboard/types"
	"dashboard/util"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
)

type Server struct {
	srv       http.Server
	resources []types.Resource
	quitChn   chan struct{}
}

func NewServer(resources []types.Resource) *Server {
	return &Server{
		resources: resources,
		quitChn:   make(chan struct{}),
	}
}

// Start
//
//	Starts the server on the next available port, returning the address if successful
func (s *Server) Start() (string, error) {
	util.Assert(s.srv.Handler == nil, "Server already started")

	serveMux := http.NewServeMux()
	for _, resource := range s.resources {
		log.Println("Adding handler for", resource.Name)
		serveMux.HandleFunc(resource.Route, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%v\n\tSource:%v\n\tRestricted:%v", resource.Name, resource.Source, resource.Restricted)
		})
	}
	s.srv.Handler = serveMux

	// Create listener on any open port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	// Start the server
	go func() {
		if err = s.srv.Serve(listener); !errors.Is(http.ErrServerClosed, err) {
			log.Println("Error starting server:", err)
		}
	}()

	// Setup listener to close server
	go func() {
		<-s.quitChn
		s.srv.Shutdown(nil)
		log.Println("Server shutdown")
	}()

	// Return the address
	addr := listener.Addr().String()
	return addr, nil
}

// Stop
//
//	Stops the server
func (s *Server) Stop() {
	close(s.quitChn)
}
