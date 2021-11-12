package chatting

type User struct {
	Email    string
	commands chan<- Command
}

func (u *User) readInput() {
	f
}
