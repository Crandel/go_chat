package chatting

// CommandID helper type for commands
type CommandID string

const (
	CmdMsg   CommandID = "/msg"
	CmdJoin            = "/join"
	CmdPing            = "/ping"
	CmdQuit            = "/quit"
	CmdRooms           = "/rooms"
	CmdUsers           = "/users"
)

// Command used by client to manage interaction with chat
type Command struct {
	id     CommandID
	client *Client
	args   []string
}
