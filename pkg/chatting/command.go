package chatting

type CommandID int

const (
	CMD_MSG CommandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_USERS
	CMD_QUIT
)

type Command struct {
	id   CommandID
	user *User
	args []string
}
