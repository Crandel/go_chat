package chatting

type User struct {
	Email    string
	commands chan<- Command
}

func (u *User) readCommand() {
	for {
		c <- u.commands
		switch c.id {
		case CMD_MSG:
			c.args

		}
	}

}
