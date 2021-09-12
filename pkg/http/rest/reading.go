package rest

import (
	"encoding/json"
	"net/http"

	rdg "github.com/Crandel/go_chat/pkg/reading"
	"github.com/gorilla/mux"
)

func rootHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Welcome")
	}
}

func listUsersHandler(rs rdg.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, _ := rs.ReadUsers()
		json.NewEncoder(w).Encode(users)
	}
}

func getUserHandler(rs rdg.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user_id, exists := vars["user_id"]
		if !exists {
			json.NewEncoder(w).Encode("user_id parameter is not exists or not valid")
		} else {
			user, err := rs.ReadUser(rdg.UserId(user_id))
			if err != nil {
				json.NewEncoder(w).Encode(err)
			} else {
				json.NewEncoder(w).Encode(user)
			}
		}
	}
}

func listRoomsHandler(rs rdg.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := rs.ReadRooms()
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
		} else {
			json.NewEncoder(w).Encode(rooms)
		}
	}
}

func getRoomHandler(rs rdg.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		room_id, exists := vars["room_id"]
		if !exists {
			json.NewEncoder(w).Encode("room_id parameter is not exists or not valid")
		} else {
			room, err := rs.ReadRoom(room_id)
			if err != nil {
				json.NewEncoder(w).Encode(err)
			} else {
				json.NewEncoder(w).Encode(room)
			}
		}
	}
}
