package chatting

type Room struct {
	Name    string
	Members map[string]*User
}

func (r *Room) broadcast() {
	for _, m := range r.Members {
		m.
	}
}
