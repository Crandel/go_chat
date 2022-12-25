package memory

import (
	"time"

	"github.com/Crandel/go_chat/internal/auth"
	r "github.com/Crandel/go_chat/internal/reading"
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
	Name       *string
	SecondName *string
	Email      *string
	Role       Role
	Nick       UserId
	Token      string
}

func ConvertUserFromSigning(su auth.SigninUser) User {
	id := UserId(su.Nick)
	token := auth.MakeToken(su.Nick, su.Password)
	return User{
		Email:      su.Email,
		Name:       su.Name,
		SecondName: su.SecondName,
		Nick:       id,
		Token:      token,
		Role:       Member,
		Created:    time.Now(),
	}
}

func (u User) ConvertUserToReading() r.User {
	return r.User{
		Name:       *u.Name,
		SecondName: *u.SecondName,
		Email:      *u.Email,
		Nick:       u.Nick.ConvertUserIdToReading(),
	}
}
