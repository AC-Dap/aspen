package service

import (
	"aspen/config"
)

func ParseService(sc config.ServiceConfig)  (*Service, error) {
	return NewService(sc.Id, sc.Remote, sc.CommitHash), nil
}
