package auth

import b64 "encoding/base64"

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

func MakeToken(nick string, password string) string {
	token := nick + ":" + password
	return b64.StdEncoding.EncodeToString([]byte(token))
}
