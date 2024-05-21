package routing

import (
	"dashboard/auth"
	"dashboard/config"
	"dashboard/util"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
)

type Server struct {
	srv       http.Server
	resources []config.Resource
	users     []auth.User
	quitChn   chan struct{}
}

func NewServer() (*Server, error) {
	resources, err := config.Read()
	if err != nil {
		return nil, err
	}
	users, err := auth.Read()
	if err != nil {
		return nil, err
	}

	return &Server{
		resources: resources,
		users:     users,
		quitChn:   make(chan struct{}),
	}, nil
}

// Start
//
//	Starts the server on the next available port, returning the address if successful
func (s *Server) Start() (string, error) {
	util.Assert(s.srv.Handler == nil, "Server already started")

	serveMux := http.NewServeMux()
	for _, resource := range s.resources {
		log.Println("Adding handler for", resource.Name)
		if resource.Name == "auth" {
			auth.AddRoutes(resource.Route, serveMux, s.users)
			continue
		}

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
