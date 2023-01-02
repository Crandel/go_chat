package memory_test

import (
	"testing"

	"github.com/Crandel/go_chat/internal/auth"
	m "github.com/Crandel/go_chat/internal/storage/memory"
)

func TestSigninUser(t *testing.T) {
	s := NewTestStorage()
	nick := "test_nick"
	pass := "pass"
	token := auth.MakeToken(nick, pass)
	sUser := auth.SigninUser{
		Nick:       nick,
		Name:       new(string),
		SecondName: new(string),
		Email:      new(string),
		Password:   "pass",
	}

	emptyNickUser := auth.SigninUser{
		Nick:       "",
		Name:       new(string),
		SecondName: new(string),
		Email:      new(string),
		Password:   "pass",
	}

	_, err := s.SigninUser(emptyNickUser)
	if err == nil {
		t.Error("SigninUser should reject the user with empty nick")
	} else {
		if err.Error() != m.EmptyNickSigninError {
			t.Errorf("SigninUser should reject with error '%s', got '%s'", m.EmptyNickSigninError, err.Error())
		}
	}

	emptyPassUser := auth.SigninUser{
		Nick:       "nick",
		Name:       new(string),
		SecondName: new(string),
		Email:      new(string),
		Password:   "",
	}

	_, err = s.SigninUser(emptyPassUser)
	if err == nil {
		t.Error("SigninUser should reject the user with empty nick")
	} else {
		if err.Error() != m.EmptyPassSigninError {
			t.Errorf("SigninUser should reject with error '%s', got '%s'", m.EmptyPassSigninError, err.Error())
		}
	}

	expToken, err := s.SigninUser(sUser)
	if err != nil {
		t.Errorf("Error %s in signin", err.Error())
	}
	if expToken != token {
		t.Errorf("Expected token %s are not equal to result token %s", token, expToken)
	}
	_, err = s.SigninUser(sUser)
	if err == nil {
		t.Error("SigninUser should reject the same user")
	}
}

func TestLoginUser(t *testing.T) {
	s := NewTestStorage()
	nick := "test_nick"
	pass := "pass"
	token := auth.MakeToken(nick, pass)
	emptyNickUser := auth.LoginUser{
		Nick:     "",
		Password: pass,
	}

	_, err := s.LoginUser(emptyNickUser)
	if err == nil {
		t.Error("LoginUser should reject the user with empty nick")
	} else {
		if err.Error() != m.EmptyNickLoginError {
			t.Errorf("LoginUser should reject with error '%s', got '%s'", m.EmptyNickLoginError, err.Error())
		}
	}

	emptyPassUser := auth.LoginUser{
		Nick:     nick,
		Password: "",
	}

	_, err = s.LoginUser(emptyPassUser)
	if err == nil {
		t.Error("LoginUser should reject the user with empty password")
	} else {
		if err.Error() != m.EmptyPassLoginError {
			t.Errorf("SigninUser should reject with error '%s', got '%s'", m.EmptyPassLoginError, err.Error())
		}
	}

	sUser := auth.SigninUser{
		Nick:       nick,
		Name:       new(string),
		SecondName: new(string),
		Email:      new(string),
		Password:   "pass",
	}

	_, err = s.SigninUser(sUser)
	if err != nil {
		t.Errorf("Error %s in signin", err.Error())
	}

	lUser := auth.LoginUser{
		Nick:     nick,
		Password: pass,
	}
	expToken, err := s.LoginUser(lUser)
	if err != nil {
		t.Errorf("Error %s in login", err.Error())
	}
	if expToken != token {
		t.Errorf("Expected token %s are not equal to result token %s",
			token, expToken)
	}
}

func TestReadAuthUsers(t *testing.T) {
	s := NewTestStorage()
	// Read empty storage
	users := s.ReadAuthUsers()
	if len(users) != 0 {
		t.Errorf("Should return empty list")
	}
	nick := "test_nick"
	pass := "pass"
	sUser := auth.SigninUser{
		Nick:       nick,
		Name:       new(string),
		SecondName: new(string),
		Email:      new(string),
		Password:   pass,
	}
	_, _ = s.SigninUser(sUser)
	s2User := auth.SigninUser{
		Nick:     nick + "2",
		Password: pass,
	}
	_, _ = s.SigninUser(s2User)
	users2 := s.ReadAuthUsers()
	if len(users2) == 0 {
		t.Errorf("Should return list with 2 users, got '%d'", len(users2))
	}
}
