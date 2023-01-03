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

func (u *User) GetName() string {
	name := ""
	if u.Name != nil {
		name = *u.Name
	}
	return name
}

func (u *User) GetSecondName() string {
	secondName := ""
	if u.SecondName != nil {
		secondName = *u.SecondName
	}
	return secondName
}

func (u *User) GetEmail() string {
	email := ""
	if u.Email != nil {
		email = *u.Email
	}
	return email
}

func (u User) ConvertUserToReading() r.User {
	return r.User{
		Name:       u.GetName(),
		SecondName: u.GetSecondName(),
		Email:      u.GetEmail(),
		Nick:       u.Nick.ConvertUserIdToReading(),
	}
}
