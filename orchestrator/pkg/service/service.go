package service

import (
	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/repository"
)

type Service struct {
	repo *repository.Repository
	expChannel *chan *orchestrator.Result
}

func NewService(repository *repository.Repository, expChannel *chan *orchestrator.Result) *Service {
	return &Service{repo: repository, expChannel: expChannel}
}