package chatting

type CommandID string

const (
	CmdMsg   CommandID = "/msg"
	CmdPing            = "/ping"
	CmdJoin            = "/join"
	CmdRooms           = "/rooms"
	CmdUsers           = "/users"
	CmdQuit            = "/quit"
)

type Command struct {
	id     CommandID
	client *Client
	args   []string
}
