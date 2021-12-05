package memory

import (
	"fmt"

	"github.com/Crandel/go_chat/pkg/auth"
	errs "github.com/Crandel/go_chat/pkg/errors"
)

func (str *Storage) SigninUser(su auth.SigninUser) (string, error) {
	const op errs.Op = "memory.Signin"
	u := ConvertUserFromSigning(su)
	if str.Users == nil {
		str.Users = make(map[UserId]User)
	}
	str.Lock()
	_, exists := str.Users[u.Email]
	if exists {
		return "", errs.New(
			op, errs.Info, fmt.Sprintf("User with email: '%s' exists", u.Email))
	}
	str.Users[u.Email] = u
	str.Unlock()
	return u.Token, nil
}

func (str *Storage) LoginUser(lu auth.LoginUser) (string, error) {
	const op errs.Op = "memory.LoginUser"
	str.RLock()
	u, exists := str.Users[UserId(lu.Email)]
	str.RUnlock()
	if !exists {
		return "", errs.New(
			op, errs.Info, fmt.Sprintf("No user with email '%s'", lu.Email))
	}
	if u.Password != lu.Password {
		return "", errs.New(
			op, errs.Info, fmt.Sprintf("User with email '%s' has wrong password", lu.Email))
	}
	return u.Token, nil
}
