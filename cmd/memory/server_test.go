package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	add "github.com/Crandel/go_chat/pkg/adding"
	ath "github.com/Crandel/go_chat/pkg/auth"
	rst "github.com/Crandel/go_chat/pkg/http/rest"
	rdn "github.com/Crandel/go_chat/pkg/reading"
	mem "github.com/Crandel/go_chat/pkg/storage/memory"
)

const user_id = "example@post.com"

type data map[string]string

func TestRestHandlers(t *testing.T) {
	uid := mem.UserId(user_id)
	testUsers := make(map[mem.UserId]mem.User)
	testRooms := make(map[string]mem.Room)
	testUser := mem.User{
		Email:      uid,
		Name:       "name",
		SecondName: "second",
		Password:   "pass",
		Token:      "token",
		Role:       mem.Member,
		Created:    time.Now(),
	}
	testMessages := make(map[mem.UserId][]mem.Message)
	testMessages[uid] = []mem.Message{}
	testRoom := mem.Room{
		Name:     "test room",
		Messages: testMessages,
	}
	testUsers[testUser.Email] = testUser
	testRooms[testRoom.Name] = testRoom
	mStorage := mem.FilledStorage(testUsers, testRooms)
	fmt.Printf("\n\n %v\n", mStorage.Users)
	aths := ath.NewService(&mStorage)
	adds := add.NewService(&mStorage)
	rdns := rdn.NewService(&mStorage)
	router := rst.InitHandlers(aths, adds, rdns)
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
			url:    "/api/users/signin",
			method: http.MethodPost,
			data: data{
				"name":        "user1",
				"email":       user_id + "1",
				"second_name": "second",
				"password":    "pass",
			},
		},
		{
			name:   "Login",
			url:    "/api/users/login",
			method: http.MethodPost,
			data:   data{"email": "example@post.com", "password": "pass"},
		},
		{
			name:   "List users",
			url:    "/api/users",
			method: http.MethodGet,
		},
		{
			name:   "Add rooms",
			url:    "/api/rooms",
			method: http.MethodPost,
			data:   data{"name": "room 1", "users": `[{"id":"example@post.com"}]`},
		},
		{
			name:   "List rooms",
			url:    "/api/rooms",
			method: http.MethodGet,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			data := new(bytes.Buffer)
			if tc.data != nil {
				json.NewEncoder(data).Encode(tc.data)
			}
			fmt.Println(srv.URL + tc.url)
			fmt.Println(data)
			req, _ := http.NewRequest(tc.method, srv.URL+tc.url, data)
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", res.Status)
			}
		})
	}
}
