package service

import (
	"aspen/external/git"
	"fmt"
	"sync"
)

// runningServices tracks ref counts of services that are currently running, determined by their ID and hash.
var runningServices = make(map[string]map[string]int)
var runningServicesLock sync.Mutex

type Status int

const (
	NotInitialized Status = iota
	Building
	Built
	Starting
	Started
	Stopping
	Stopped
)

type Service struct {
	id     string
	status Status

	// Remote git repo that contains the service code
	repo git.Repo

	buildCommand string
	startCommand string

	// Hash of the repo + build/start commands, to detect changes
	// Calculated when the service is built
	hash string
}

func NewService(id, remote, commitHash, buildCommand, startCommand string) *Service {
	return &Service{
		id:           id,
		status:       NotInitialized,
		repo:         git.NewRepo(getServiceFolder(id), remote, commitHash),
		buildCommand: buildCommand,
		startCommand: startCommand,
	}
}

func (s *Service) GetID() string {
	return s.id
}

func (s *Service) GetStatus() Status {
	return s.status
}

// Build initializes the service by cloning the source code and building the docker image.
// If the service source code is already cloned or built, these steps will be skipped.
// This method should be called once before starting the service.
func (s *Service) Build() error {
	if s.status != NotInitialized {
		return fmt.Errorf("trying to build service %s again, current status is %s", s.id, s.status)
	}

	s.status = Building

	// TODO: Implement actual build logic here, e.g., cloning the repo, running build commands, etc.
	// TODO: Also calculate hash

	// First clone repo if it's not up to date
	if !s.repo.Updated() {
		s.repo.Clone()
	}

	// Run docker build in repo

	// Calculate hash

	s.status = Built
	return nil
}

// Start increases the ref count for this service and starts the docker container.
// If the docker container is already running, nothing will happen.
// This method should be called once after Build; if called multiple times, it will return an error.
func (s *Service) Start() error {
	if s.status != Built {
		return fmt.Errorf("trying to build unbuilt or already started service %s, current status is %s", s.id, s.status)
	}

	s.status = Starting
	// First update ref count to ensure the service isn't killed while starting
	runningServicesLock.Lock()
	if _, exists := runningServices[s.id]; !exists {
		runningServices[s.id] = make(map[string]int)
	}
	runningServices[s.id][s.hash]++
	runningServicesLock.Unlock()

	// TODO: Implement actual start logic here, e.g., running the start command, etc.
	s.status = Started

	return nil
}

// Stop decrements the ref count for this service and, if the count reaches zero, stops the service.
// This method should be called once after Start; if called multiple times, it will return an error.
func (s *Service) Stop() error {
	if s.status != Started {
		return fmt.Errorf("trying to stop not-running service %s, current status is %s", s.id, s.status)
	}

	s.status = Stopping
	// Decrease ref count, and if it reaches zero, stop the service
	runningServicesLock.Lock()
	defer runningServicesLock.Unlock()

	if _, exists := runningServices[s.id]; !exists {
		return fmt.Errorf("service %s is missing ref count", s.id)
	}
	if count := runningServices[s.id][s.hash]; count <= 0 {
		return fmt.Errorf("service %s has 0 ref count", s.id)
	}

	runningServices[s.id][s.hash]--
	if runningServices[s.id][s.hash] == 0 {
		// TODO: Implement actual stop logic here, e.g., killing the process, etc.
	}
	s.status = Stopped

	return nil
}
