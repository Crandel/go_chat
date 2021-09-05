package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Crandel/go_chat/pkg/adding"
)

func addRoomHandler(as adding.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ar adding.Room
		if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		name, err := as.AddRoom(ar)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(name)
	}

}
