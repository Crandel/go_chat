package main_test

import (
	"bytes"
	"encoding/json"
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

type data map[string]interface{}

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
	testAddingUser := add.User{ID: user_id}
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
			data:   data{"email": user_id + "1", "password": "pass"},
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
			data:   data{"name": "room 1", "users": []add.User{testAddingUser}},
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
				err := json.NewEncoder(data).Encode(tc.data)
				if err != nil {
					t.Fatal(err)
				}
			}
			req, err := http.NewRequest(tc.method, srv.URL+tc.url, data)
			if err != nil {
				t.Fatalf("Could not create request %v", err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Could not send %s request: %v", tc.method, err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK; got %v", res.Status)
			}
			// b, err := ioutil.ReadAll(res.Body)
			// if err != nil {
			// 	t.Fatalf("could not read response: %v", err)
			// }
			// fmt.Println(string(bytes.TrimSpace(b)))
		})
	}
}
