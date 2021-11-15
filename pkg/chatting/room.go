package chatting

import errs "github.com/Crandel/go_chat/pkg/errors"

type Room struct {
	Name    string
	Members map[string]*User
}

func (r *Room) broadcast(message string) {
	for _, m := range r.Members {
		m.WriteMsg(message)
	}
}

func (r *Room) addUser(u *User) error {
	const op errs.Op = "chatting.Room.addUser"
	if r.haveUser(u) {
		return errs.New(op, errs.Info, "User already in room")
	}
	r.Members[u.Nick] = u
	return nil
}

func (r *Room) haveUser(u *User) bool {
	const op errs.Op = "chatting.Room.haveUser"
	u, ok := r.Members[u.Nick]
	if ok {
		return true
	}
	return false
}
