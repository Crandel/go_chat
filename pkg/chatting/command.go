package chatting

type CommandID string

const (
	CMD_MSG   CommandID = "/msg"
	CMD_JOIN            = "/join"
	CMD_ROOMS           = "/rooms"
	CMD_USERS           = "/users"
	CMD_QUIT            = "/quit"
)

type Command struct {
	id   CommandID
	user *User
	args []string
}
