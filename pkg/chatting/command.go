package chatting

// CommandID helper type for commands
type CommandID string

const (
	CmdMsg   CommandID = "/msg"
	CmdPing            = "/ping"
	CmdJoin            = "/join"
	CmdRooms           = "/rooms"
	CmdUsers           = "/users"
	CmdQuit            = "/quit"
)

// Command used by client to manage interaction with chat
type Command struct {
	id     CommandID
	client *Client
	args   []string
}
