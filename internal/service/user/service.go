package user

import "context"

type Service struct {
	Repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) GetAllUsers() ([]UserResponse, error) {
	return s.Repo.GetAll(context.Background())
}

func (s *Service) GetByID(id int) (User, error) {
	return s.Repo.GetByID(context.Background(), id)
}
