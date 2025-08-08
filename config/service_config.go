package config

import "aspen/router/service"

type ServiceConfig struct {
	Id         string
	Remote     string
	CommitHash string
}

func (sc ServiceConfig) Parse() (*service.Service, error) {
	return service.NewService(sc.Id, sc.Remote, sc.CommitHash), nil
}
