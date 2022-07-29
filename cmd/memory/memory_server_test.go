package main_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

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

type Login struct {
	Token string `json:"token"`
}

func runRequest(d data, method string, url string) ([]byte, error) {
	data := new(bytes.Buffer)
	if d != nil {
		err := json.NewEncoder(data).Encode(d)
		if err != nil {
			return []byte{}, err
		}
	}
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return []byte{}, err
	}
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return []byte{}, errors.New(fmt.Sprintf("Expected status OK; got %v", res.Status))
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return bytes.TrimSpace(b), nil
}

func TestHandlers(t *testing.T) {
	mStorage := mem.NewStorage()

	aths := ath.NewService(&mStorage)
	adds := add.NewService(&mStorage)
	rdns := rdn.NewService(&mStorage)
	chts := cht.NewService(&mStorage)
	router := ntw.NewRouter(aths, adds, rdns, chts)
	srv := httptest.NewServer(router)
	defer srv.Close()

	signRes, err := runRequest(
		data{
			"name":        "user",
			"email":       user_id,
			"second_name": "second",
			"password":    "pass",
		},
		http.MethodPost,
		srv.URL+"/api/users/signin",
	)
	if err != nil {
		t.Fatal(err)
	}
	var token Login
	err = json.Unmarshal(signRes, &token)
	if err != nil {
		t.Fatal(err)
	}
	logRes, err := runRequest(
		data{
			"email":    user_id,
			"password": "pass",
		},
		http.MethodPost,
		srv.URL+"/api/users/login",
	)
	if err != nil {
		t.Fatal(err)
	}
	var login Login
	err = json.Unmarshal(logRes, &login)

	if token != login {
		t.Fatalf("Token from signing '%s' != '%s' token from login", token.Token, login.Token)
	}

	tt := []struct {
		name        string
		url         string
		method      string
		data        data
		compareResp func(resp interface{}) bool
		err         string
	}{
		{
			name:   "Health",
			url:    "/health",
			method: http.MethodGet,
			compareResp: func(resp interface{}) bool {
				respConv := resp.(map[string]interface{})
				testResp := data{
					"status": "OK",
				}
				if len(testResp) != len(respConv) {
					return false
				}

				if reflect.DeepEqual(respConv, testResp) {
					fmt.Printf("Test is not equal with resp:\n%v != %v\n", respConv, testResp)
					return false
				}
				return true
			},
		},
		{
			name:   "List users",
			url:    "/api/users",
			method: http.MethodGet,
			compareResp: func(resp interface{}) bool {
				respConv := resp.([]interface{})
				testResp := []data{
					{
						"email":       "example@post.com",
						"name":        "user",
						"second_name": "second",
					},
				}
				if len(testResp) != len(respConv) {
					return false
				}

				for i, d := range respConv {
					if reflect.DeepEqual(d, testResp[i]) {
						fmt.Printf("Test is not equal with resp:\n%v != %v\n", d, testResp[i])
						return false
					}
				}
				return true
			},
		},
		{
			name:   "Add rooms",
			url:    "/api/rooms",
			method: http.MethodPost,
			data:   data{"name": "room 1"},
			compareResp: func(resp interface{}) bool {
				respConv := resp.(map[string]interface{})
				testResp := map[string]string{
					"name": "room 1",
				}

				if reflect.DeepEqual(testResp, respConv) {
					fmt.Printf("Test is not equal with resp:\n%v != %v\n", testResp, respConv)
					return false
				}
				return true
			},
		},
		{
			name:   "List rooms",
			url:    "/api/rooms",
			method: http.MethodGet,
			compareResp: func(resp interface{}) bool {
				respConv := resp.([]interface{})
				testResp := []data{
					{
						"messages": nil,
						"name":     "room 1",
					},
				}
				if len(testResp) != len(respConv) {
					fmt.Printf("Different length of test and resp:\n%v\n\n%v\n", testResp, respConv)
					return false
				}
				for i, d := range testResp {
					if reflect.DeepEqual(d, respConv[i]) {
						fmt.Printf("Test is not equal with resp:\n%v != %v\n", d, respConv[i])
						return false
					}
				}
				return true
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			response, err := runRequest(tc.data, tc.method, srv.URL+tc.url)
			if err != nil {
				t.Fatal(err)
			}
			var unmResponse interface{}
			err = json.Unmarshal(response, &unmResponse)
			if err != nil {
				t.Fatal(err)
			}

			if !tc.compareResp(unmResponse) {
				t.Fatalf("Response :%v is not equal with testing", unmResponse)
			}
		})
	}
}
