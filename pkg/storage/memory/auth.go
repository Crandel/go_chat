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
	_, exists := str.Users[u.Nick]
	if exists {
		return "", errs.New(
			op, errs.Info, fmt.Sprintf("User with nick: '%s' exists", u.Nick))
	}
	str.Users[u.Nick] = u
	str.Unlock()
	return u.Token, nil
}

func (str *Storage) LoginUser(lu auth.LoginUser) (string, error) {
	const op errs.Op = "memory.LoginUser"
	str.RLock()
	u, exists := str.Users[UserId(lu.Nick)]
	str.RUnlock()
	if !exists {
		return "", errs.New(
			op, errs.Info, fmt.Sprintf("No user with nick '%s'", lu.Nick))
	}
	if u.Password != lu.Password {
		return "", errs.New(
			op, errs.Info, fmt.Sprintf("User with nick '%s' has wrong password", lu.Nick))
	}
	return u.Token, nil
}
