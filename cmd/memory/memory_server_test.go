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
	cht "github.com/Crandel/go_chat/pkg/chatting"
	ntw "github.com/Crandel/go_chat/pkg/network"
	rdn "github.com/Crandel/go_chat/pkg/reading"
	mem "github.com/Crandel/go_chat/pkg/storage/memory"
)

const (
	user_id   = "example@post.com"
	room_name = "test room"
)

type data map[string]interface{}

func TestHandlers(t *testing.T) {
	cnow := time.Now()
	uid := mem.UserId(user_id)
	testUsers := make(map[mem.UserId]mem.User)
	testRooms := make(map[string]mem.Room)
	testMessages := make(map[int]mem.Message)
	testUser := mem.User{
		Email:      uid,
		Name:       "name",
		SecondName: "second",
		Password:   "pass",
		Token:      "token",
		Role:       mem.Member,
		Created:    cnow,
	}
	testMessages[0] = mem.Message{
		ID:       0,
		UserId:   uid,
		RoomName: room_name,
		Payload:  "Test message",
		Created:  cnow,
	}
	testRoom := mem.Room{
		Name:    room_name,
		Created: cnow,
	}
	testUsers[testUser.Email] = testUser
	testRooms[testRoom.Name] = testRoom
	mStorage := mem.FilledStorage(testUsers, testRooms, testMessages)
	aths := ath.NewService(&mStorage)
	adds := add.NewService(&mStorage)
	rdns := rdn.NewService(&mStorage)
	chts := cht.NewService(&mStorage)
	router := ntw.InitHandlers(aths, adds, rdns, chts)
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
			data:   data{"email": user_id, "password": "pass"},
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
			data:   data{"name": "room 1"},
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
