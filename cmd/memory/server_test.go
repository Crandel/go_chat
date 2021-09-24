package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

const user_id = "example@post.com"

type data map[string]string

func TestRestHandlers(t *testing.T) {
	memory := memory.NewStorage()
	ls := login.NewService(&memory)
	sis := signin.NewService(&memory)
	as := adding.NewService(&memory)
	rs := reading.NewService(&memory)
	router := rest.InitHandlers(ls, sis, as, rs)
	srv := httptest.NewServer(router)
	defer srv.Close()
	client := &http.Client{}
	tt := []struct {
		name          string
		url           string
		method        string
		data          data
		response_keys []string
		err           string
	}{
		{
			name:   "Signin",
			url:    "/api/signin",
			method: http.MethodPost,
			data: data{
				"name":        "user1",
				"email":       user_id,
				"second_name": "second",
				"password":    "pass",
			},
			response_keys: []string{"token"},
		},
		{
			name:          "Login",
			url:           "/api/login",
			method:        http.MethodPost,
			data:          data{"email": "example@post.com", "password": "pass"},
			response_keys: []string{"token"},
		},
		{
			name:          "List users",
			url:           "/api/users",
			method:        http.MethodGet,
			response_keys: []string{"id", "email", "name", "second_name"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			data := new(bytes.Buffer)
			if tc.data != nil {
				json.NewEncoder(data).Encode(tc.data)
			}
			fmt.Println(srv.URL + tc.url)
			req, _ := http.NewRequest(tc.method, srv.URL+tc.url, data)
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", res.Status)
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			sb := string(body)
			fmt.Println(sb)
		})
	}
}
