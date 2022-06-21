package auth

import (
	"log"
)

type Repository interface {
	LoginUser(LoginUser) (string, error)
	SigninUser(SigninUser) (string, error)
}

type Service interface {
	LoginUser(LoginUser) (AuthResponse, error)
	SigninUser(SigninUser) (AuthResponse, error)
}

type service struct {
	r Repository
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) LoginUser(u LoginUser) (AuthResponse, error) {
	token, err := s.r.LoginUser(u)
	if err != nil {
		log.Println("Error while login:", err)
		return AuthResponse{}, err
	}
	return AuthResponse{Token: token}, nil
}

func (s *service) SigninUser(u SigninUser) (AuthResponse, error) {
	token, err := s.r.SigninUser(u)
	if err != nil {
		log.Println("Error while signin:", err)
		return AuthResponse{}, err
	}
	return AuthResponse{Token: token}, nil
}
