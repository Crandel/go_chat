package auth

import (
	"log"
)

// Repository for auth
type Repository interface {
	LoginUser(LoginUser) (string, error)
	SigninUser(SigninUser) (string, error)
}

// Service for auth
type Service interface {
	LoginUser(LoginUser) (Response, error)
	SigninUser(SigninUser) (Response, error)
}

type service struct {
	r Repository
}

// Response return token
type Response struct {
	Token string `json:"token"`
}

// NewService will return Service
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) LoginUser(u LoginUser) (Response, error) {
	token, err := s.r.LoginUser(u)
	if err != nil {
		log.Println("Error while login:", err)
		return Response{}, err
	}
	return Response{Token: token}, nil
}

func (s *service) SigninUser(u SigninUser) (Response, error) {
	token, err := s.r.SigninUser(u)
	if err != nil {
		log.Println("Error while signin:", err)
		return Response{}, err
	}
	return Response{Token: token}, nil
}
