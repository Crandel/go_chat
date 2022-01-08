package chatting

import (
	"fmt"

	errs "github.com/Crandel/go_chat/pkg/errors"
)

type Room struct {
	Name    string
	Clients map[*Client]struct{}
}

func (r *Room) broadcast(sender *Client, message string) {
	for m := range r.Clients {
		if m != sender {
			fmt.Printf("Broadcast message %s \n", message)
			m.WriteMsg(message)
		}
	}
}

func (r *Room) addUser(u *Client) error {
	const op errs.Op = "chatting.Room.addUser"
	if r.haveUser(u) {
		return errs.New(op, errs.Info, "User already in room")
	}
	r.Clients[u] = struct{}{}
	return nil
}

func (r *Room) haveUser(u *Client) bool {
	const op errs.Op = "chatting.Room.haveUser"
	if _, ok := r.Clients[u]; ok {
		return true
	}
	return false
}
