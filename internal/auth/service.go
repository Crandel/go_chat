package auth

import (
	"log/slog"

	"gitlab.com/greyxor/slogor"
)

// Repository for auth
type Repository interface {
	LoginUser(LoginUser) (string, error)
	SigninUser(SigninUser) (string, error)
	ReadAuthUsers() []AuthUser
}

// Service for auth
type Service interface {
	LoginUser(LoginUser) (Response, error)
	SigninUser(SigninUser) (Response, error)
	ReadAuthUsers() []AuthUser
}

type service struct {
	r   Repository
	log slog.Logger
}

// Response return token
type Response struct {
	Token string `json:"token"`
}

// NewService will return Service
func NewService(r Repository) Service {
	authLog := slog.With(
		slog.Group("auth"),
	)
	return &service{r, *authLog}
}

func (s *service) LoginUser(u LoginUser) (Response, error) {
	token, err := s.r.LoginUser(u)
	if err != nil {
		s.log.Error("Error while login:", slogor.Err(err))
		return Response{}, err
	}
	return Response{Token: token}, nil
}

func (s *service) SigninUser(u SigninUser) (Response, error) {
	token, err := s.r.SigninUser(u)
	if err != nil {
		s.log.Error("Error while signin:", slogor.Err(err))
		return Response{}, err
	}
	return Response{Token: token}, nil
}

func (s *service) ReadAuthUsers() []AuthUser {
	return s.r.ReadAuthUsers()
}
