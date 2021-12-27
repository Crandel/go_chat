package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Crandel/go_chat/pkg/adding"
)

type RoomResponse struct {
	Name string `json:"name"`
}

func AddRoomHandler(as adding.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ar adding.Room
		if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		name, err := as.AddRoom(ar.Name)
		if name == "" {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		resp := RoomResponse{
			Name: name,
		}
		json.NewEncoder(w).Encode(resp)
	}

}
