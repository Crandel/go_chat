package auth

import (
	b64 "encoding/base64"
	"fmt"
)

const AuthKey = "authUser"

// LoginUser is for logging
type LoginUser struct {
	Nick     string `json:"nick"`
	Password string `json:"password"`
}

// SigninUser is for signin
type SigninUser struct {
	Nick       string  `json:"nick"`
	Name       *string `json:"name"`
	SecondName *string `json:"second_name"`
	Email      *string `json:"email"`
	Password   string  `json:"password"`
}

type AuthUser struct {
	Nick  string
	Token string
}

func (a *LoginUser) String() string {
	return fmt.Sprintf("Login user %s", a.Nick)
}

func (a *SigninUser) String() string {
	return fmt.Sprintf("Signin user %s", a.Nick)
}

func (a *AuthUser) String() string {
	return fmt.Sprintf("Auth user %s", a.Nick)
}

func MakeToken(nick string, password string) string {
	token := nick + ":" + password
	return b64.StdEncoding.EncodeToString([]byte(token))
}
