package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Crandel/go_chat/pkg/adding"
)

type RoomResponse struct {
	Name   string   `json:"name"`
	Errors []string `json:"errors"`
}

func AddRoomHandler(as adding.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ar adding.Room
		if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		name, err := as.AddRoom(ar)
		if name == "" {
			err_msg := ""
			for _, e := range err {
				err_msg = err_msg + e.Error()
			}
			http.Error(w, err_msg, http.StatusBadGateway)
			return
		}
		var errors []string
		for _, e := range err {
			errors = append(errors, e.Error())
		}
		resp := RoomResponse{
			Name: name, Errors: errors,
		}
		json.NewEncoder(w).Encode(resp)
	}

}
