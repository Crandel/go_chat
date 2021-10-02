package rest

import (
	"encoding/json"
	"net/http"

	rdg "github.com/Crandel/go_chat/pkg/reading"
	"github.com/gorilla/mux"
)

func ListUsersHandler(rs rdg.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := rs.ReadUsers()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		}
		json.NewEncoder(w).Encode(users)
	}
}

func GetUserHandler(rs rdg.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user_id, exists := vars["user_id"]
		if !exists {
			json.NewEncoder(w).Encode("user_id parameter is not exists or not valid")
		} else {
			user, err := rs.ReadUser(rdg.UserId(user_id))
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
			} else {
				json.NewEncoder(w).Encode(user)
			}
		}
	}
}

func ListRoomsHandler(rs rdg.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := rs.ReadRooms()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		} else {
			json.NewEncoder(w).Encode(rooms)
		}
	}
}

func GetRoomHandler(rs rdg.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		room_id, exists := vars["room_id"]
		if !exists {
			json.NewEncoder(w).Encode("room_id parameter is not exists or not valid")
		} else {
			room, err := rs.ReadRoom(room_id)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
			} else {
				json.NewEncoder(w).Encode(room)
			}
		}
	}
}
