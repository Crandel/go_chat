package memory

import (
	"time"

	"github.com/Crandel/go_chat/pkg/auth"
	r "github.com/Crandel/go_chat/pkg/reading"
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
	Created    time.Time
	Role       Role
	Nick       UserId
	Name       string
	SecondName *string
	Password   string
	Token      string
	Email      *string
}

func ConvertUserFromSigning(su auth.SigninUser) User {
	id := UserId(su.Nick)
	token := uuid.New().String()
	return User{
		Email:      su.Email,
		Name:       su.Name,
		SecondName: su.SecondName,
		Nick:       id,
		Password:   su.Password,
		Token:      token,
		Role:       Member,
		Created:    time.Now(),
	}
}

func (u User) ConvertUserToReading() r.User {
	return r.User{
		Name:       u.Name,
		SecondName: *u.SecondName,
		Email:      *u.Email,
		Nick:       u.Nick.ConvertUserIdToReading(),
	}
}
