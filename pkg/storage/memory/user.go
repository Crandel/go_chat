package memory

import (
	"errors"
	"time"

	"github.com/Crandel/go_chat/pkg/login"
	"github.com/Crandel/go_chat/pkg/reading"
	"github.com/Crandel/go_chat/pkg/signin"
	"github.com/google/uuid"
)

type Role string

const (
	Member Role = "Member"
	Admin  Role = "Admin"
)

type UserId string

type User struct {
	Email      UserId
	Name       string
	SecondName string
	Password   string
	Token      string
	Role       Role
	Created    time.Time
}

func (s *Storage) SigninUser(u signin.User) (signin.SigninResponse, error) {
	token := uuid.New().String()
	su := User{
		Email:      UserId(u.Email),
		Name:       u.Name,
		SecondName: u.SecondName,
		Password:   u.Password,
		Token:      token,
		Role:       Member,
		Created:    time.Now(),
	}
	s.Users[UserId(u.Email)] = su
	return signin.SigninResponse{Id: u.Email, Token: token}, nil
}

func (s *Storage) LoginUser(lu login.User) (string, error) {
	u, exists := s.Users[UserId(lu.Email)]
	if !exists {
		return "", errors.New("No user with email: " + lu.Email)
	}
	if u.Password != lu.Password {
		return "", errors.New("User with email" + lu.Email + "has wrong password")
	}
	return u.Token, nil
}

func (s *Storage) ReadUsers() ([]reading.User, error) {
	var users = []reading.User{}
	for id, u := range s.Users {
		users = append(users, reading.User{
			ID:         reading.UserId(id),
			Name:       u.Name,
			SecondName: u.SecondName,
			Email:      string(u.Email)})
	}
	return users, nil
}

func (s *Storage) ReadUser(uid reading.UserId) (reading.User, error) {
	umid := UserId(string(uid))
	s_user, exists := s.Users[umid]
	if !exists {
		return reading.User{}, errors.New("")
	}
	r_user := reading.User{
		ID:         reading.UserId(s_user.Email),
		Name:       s_user.Name,
		SecondName: s_user.SecondName,
		Email:      string(s_user.Email),
	}
	return r_user, nil
}
