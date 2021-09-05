package memory

import (
	"time"

	r "github.com/Crandel/go_chat/pkg/reading"
	s "github.com/Crandel/go_chat/pkg/signin"
	"github.com/google/uuid"
)

type Role string

const (
	Member Role = "Member"
	Admin  Role = "Admin"
)

type UserId string

func (muid UserId) ConvertUserIdToReading() r.UserId {
	return r.UserId(string(muid))
}

func ConvertUserIdFromReading(rid r.UserId) UserId {
	return UserId(string(rid))
}

type User struct {
	Email      UserId
	Name       string
	SecondName string
	Password   string
	Token      string
	Role       Role
	Created    time.Time
}

func ConvertUserFromSigning(su s.User) User {
	id := UserId(su.Email)
	token := uuid.New().String()
	return User{
		Email:      id,
		Name:       su.Name,
		SecondName: su.SecondName,
		Password:   su.Password,
		Token:      token,
		Role:       Member,
		Created:    time.Now(),
	}
}

func (u User) ConvertUserToReading() r.User {
	return r.User{
		ID:         u.Email.ConvertUserIdToReading(),
		Name:       u.Name,
		SecondName: u.SecondName,
		Email:      string(u.Email),
	}
}
