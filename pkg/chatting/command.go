package chatting

type CommandID int

const (
	CMD_JOIN CommandID = iota
	CMD_ROOMS
	CMD_USERS
	CMD_QUIT
)

type command struct {
	id   CommandID
	user User
}
