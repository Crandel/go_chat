package chatting

import (
	"fmt"

	errs "github.com/Crandel/go_chat/pkg/errors"
)

type Room struct {
	Name    string
	Clients []*Client
}

func (r *Room) broadcast(sender *Client, message string) {
	for _, m := range r.Clients {
		if m != sender {
			fmt.Printf("Broadcast message %s \n", message)
			m.WriteMsg(message)
		}
	}
}

func (r *Room) addUser(u *Client) error {
	const op errs.Op = "chatting.Room.addUser"
	if ok, _ := r.haveUser(u); ok {
		return errs.New(op, errs.Info, "User already in room")
	}
	r.Clients = append(r.Clients, u)
	return nil
}

func (r *Room) haveUser(u *Client) (bool, int) {
	const op errs.Op = "chatting.Room.haveUser"
	for i, m := range r.Clients {
		if m == u {
			return true, i
		}
	}
	return false, 0
}
