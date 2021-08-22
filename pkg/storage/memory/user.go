package memory

import (
	"errors"
	"time"

	"github.com/Crandel/go_chat/pkg/login"
	"github.com/Crandel/go_chat/pkg/signin"
	"github.com/google/uuid"
)

type Role string

const (
	Member Role = "Member"
	Admin  Role = "Admin"
)

type User struct {
	ID         string
	Name       string
	SecondName string
	Email      string
	Password   string
	Token      string
	Role       Role
	Created    time.Time
}

func (s *Storage) SigninUser(u signin.User) (signin.SigninResponse, error) {
	id := uuid.New().String()
	token := uuid.New().String()
	su := User{id, u.Name, u.SecondName, u.Email, u.Password, token, Member, time.Now()}
	_ = append(s.Users, su)
	return signin.SigninResponse{id, token}, nil
}

func (s *Storage) LoginUser(lu login.User) (string, error) {
	for _, u := range s.Users {
		if lu.Email == u.Email && lu.Password == u.Password {
			return u.Token, nil
		}
	}
	return "", errors.New("No user with email: " + lu.Email)
}
