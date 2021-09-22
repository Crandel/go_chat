package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/http/rest"
	"github.com/Crandel/go_chat/pkg/login"
	"github.com/Crandel/go_chat/pkg/reading"
	"github.com/Crandel/go_chat/pkg/signin"
	"github.com/Crandel/go_chat/pkg/storage/memory"
)

func TestRestHandlers(t *testing.T) {
	memory := memory.NewStorage()
	ls := login.NewService(&memory)
	sis := signin.NewService(&memory)
	as := adding.NewService(&memory)
	rs := reading.NewService(&memory)
	router := rest.InitHandlers(ls, sis, as, rs)
	srv := httptest.NewServer(router)
	defer srv.Close()

	tt := []struct {
		name          string
		url           string
		method        string
		data          []byte
		response_keys []string
		err           string
	}{
		{
			name:          "Login",
			url:           "%s/api/login",
			method:        http.MethodPost,
			data:          []byte(`{"email": "example@post.com","password": "pass"}`),
			response_keys: []string{"token"},
		},
		{
			name:          "Signin",
			url:           "%s/api/signin",
			method:        http.MethodPost,
			data:          []byte(`{"name": "Name","second_name": "Second","email": "example@post.com","password": "pass"}`),
			response_keys: []string{"token"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
