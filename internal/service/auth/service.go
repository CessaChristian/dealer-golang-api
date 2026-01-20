package auth

import (
	"context"
	"dealer_golang_api/internal/service/user"
	"dealer_golang_api/utils"
	"errors"
)

type Service struct {
	UserRepo user.UserRepository
}

func NewService(repo user.UserRepository) *Service {
	return &Service{
		UserRepo: repo,
	}
}

func (s *Service) Register(req RegisterRequest) error {
	// check duplicate email
	_, err := s.UserRepo.GetByEmail(context.Background(), req.Email)
	if err == nil {
		return errors.New("email already registered")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	newUser := user.User{
		Name: 	  req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     "user",
	}

	return s.UserRepo.Create(context.Background(), newUser)
}

func (s *Service) Login(req LoginRequest) (string, user.UserResponse, error) {
	u, err := s.UserRepo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		return "", user.UserResponse{}, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, u.Password) {
		return "", user.UserResponse{}, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(u.ID, u.Role)
	if err != nil {
		return "", user.UserResponse{}, err
	}

	resp := user.UserResponse{
		ID:       u.ID,
		Name: u.Name,
		Email:    u.Email,
		Role:     u.Role,
	}

	return token, resp, nil
}
