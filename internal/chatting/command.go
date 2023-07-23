package chatting

import "fmt"

// CommandID helper type for commands
type CommandID string

const (
	CmdMsg   CommandID = "msg"
	CmdJoin  CommandID = "join"
	CmdPing  CommandID = "ping"
	CmdQuit  CommandID = "quit"
	CmdRooms CommandID = "rooms"
	CmdUsers CommandID = "users"
)

// Command used by client to manage interaction with chat
type Command struct {
	id     CommandID
	client *Client
	args   []string
}

func ConvertToCommandID(input string) (CommandID, error) {
	commandID := CommandID(input)
	switch commandID {
	case CmdMsg, CmdJoin, CmdPing, CmdQuit, CmdRooms, CmdUsers:
		return commandID, nil
	default:
		return "", fmt.Errorf("invalid CommandID: %s", input)
	}
}
