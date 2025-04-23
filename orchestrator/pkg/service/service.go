package service

import (
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	return &Service{repo: repository}
}