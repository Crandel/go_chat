package chatting

type Message struct {
	user *User
	room *Room
	msg  string
}
