package brand

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

	// check if exists
	_, err := s.Repo.GetByName(context.Background(), name)
	if err == nil {
		return errors.New("brand already exists")
	}

	return s.Repo.Create(context.Background(), name)
}

func (s *Service) GetAll() ([]Brand, error) {
	return s.Repo.GetAll(context.Background())
}
