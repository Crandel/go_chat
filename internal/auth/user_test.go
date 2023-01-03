package auth_test

import (
	b64 "encoding/base64"
	"testing"

	"github.com/Crandel/go_chat/internal/auth"
)

func TestMakeToken(t *testing.T) {
	nick := "test_nick"

	password := "tst pass"
	token := b64.StdEncoding.EncodeToString([]byte(nick + ":" + password))
	gToken := auth.MakeToken(nick, password)

	if token != gToken {
		t.Errorf("Got token: '%s', expected: '%s'", gToken, token)
	}
}
