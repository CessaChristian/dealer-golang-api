package vtype

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func normalize(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

func (s *Service) Create(name string) error {
	name = normalize(name)

	_, err := s.Repo.GetByName(context.Background(), name)
	if err == nil {
		return errors.New("type already exists")
	}

	return s.Repo.Create(context.Background(), name)
}

func (s *Service) GetAll() ([]Type, error) {
	return s.Repo.GetAll(context.Background())
}
