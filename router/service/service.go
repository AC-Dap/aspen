package service

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/rs/zerolog/log"
)

// runningServices tracks ref counts of services that are currently running, determined by their ID.
var runningServices = make(map[string]int)
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
	repo Repo
}

func NewService(id, remote, commitHash string) *Service {
	return &Service{
		id:     id,
		status: NotInitialized,
		repo:   NewRepo(getServiceFolder(id), remote, commitHash),
	}
}

func (s *Service) GetID() string {
	return s.id
}

func (s *Service) GetStatus() Status {
	return s.status
}

// Build initializes the service by cloning the source code and building the docker image.
// If the currently built image is out of date, this will update it.
// This method should be called once before starting the service.
func (s *Service) Build() error {
	if s.status != NotInitialized {
		return fmt.Errorf("trying to build service %s again, current status is %s", s.id, s.status)
	}
	log.Info().Str("service", s.id).Msg("Building service")
	s.status = Building

	// First clone repo if it's not up to date
	if !s.repo.Updated() {
		s.repo.Clone()
	}

	// Run docker build in repo
	cmd := exec.Command("docker", "compose", "build")
	cmd.Dir = s.repo.folder
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error building service %s: %w, output: %s", s.id, err, string(output))
	}

	s.status = Built
	return nil
}

// Start increases the ref count for this service and starts the docker container.
// If an outdated version is running, this will reload it with the updated version.
// This method should be called once after Build; if called multiple times, it will return an error.
func (s *Service) Start() error {
	if s.status != Built {
		return fmt.Errorf("trying to build unbuilt or already started service %s, current status is %s", s.id, s.status)
	}
	log.Info().Str("service", s.id).Msg("Starting service")
	s.status = Starting

	// First update ref count to ensure the service isn't killed while starting
	runningServicesLock.Lock()
	runningServices[s.id]++
	runningServicesLock.Unlock()

	// Launch service using docker compose
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = s.repo.folder
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error starting service %s: %w, output: %s", s.id, err, string(output))
	}

	s.status = Started

	return nil
}

// Stop decrements the ref count for this service and, if the count reaches zero, stops the service.
// This method should be called once after Start; if called multiple times, it will return an error.
func (s *Service) Stop() error {
	if s.status != Started {
		return fmt.Errorf("trying to stop not-running service %s, current status is %s", s.id, s.status)
	}
	log.Info().Str("service", s.id).Msg("Stopping service")
	s.status = Stopping

	// Decrease ref count, and if it reaches zero, stop the service
	runningServicesLock.Lock()
	defer runningServicesLock.Unlock()

	if count := runningServices[s.id]; count <= 0 {
		return fmt.Errorf("service %s has 0 ref count", s.id)
	}

	runningServices[s.id]--
	if runningServices[s.id] == 0 {
		// Stop service with docker compose
		cmd := exec.Command("docker", "compose", "down")
		cmd.Dir = s.repo.folder
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("error stopping service %s: %w, output: %s", s.id, err, string(output))
		}
	}
	s.status = Stopped

	return nil
}
