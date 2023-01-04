package chatting

import (
	"strings"

	lg "github.com/Crandel/go_chat/internal/logging"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	commands chan<- Command
	Nick     string
}

func (c *Client) ReadCommands() {
	const op = "chatting#user#ReadCommands "
	defer c.conn.Close()
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Log(lg.Warning, op, err.Error())
			c.commands <- Command{
				id:     CmdQuit,
				client: c,
			}
			return
		}
		rawCommand := string(p)
		args := strings.Split(rawCommand, " ")
		cmd := strings.TrimSpace(args[0])
		log.Log(lg.Debug, op, "Command: ", cmd)
		var cmdID CommandID
		if !strings.HasPrefix(cmd, "/") {
			cmdID = CmdMsg
		} else {
			switch cmd {
			case string(CmdPing):
				cmdID = CmdPing
			case string(CmdJoin):
				cmdID = CmdJoin
			case string(CmdUsers):
				cmdID = CmdUsers
			case string(CmdRooms):
				cmdID = CmdRooms
			case string(CmdQuit):
				cmdID = CmdQuit
			default:
				c.WriteMsg("ERR: Unknown command " + cmd)
				continue
			}
		}
		c.commands <- Command{
			id:     cmdID,
			client: c,
			args:   args,
		}
	}
}

func (c *Client) WriteMsg(message string) {
	messageType := websocket.TextMessage
	if err := c.conn.WriteMessage(messageType, []byte(message)); err != nil {
		log.Log(lg.Warning, "WriteMsg ", err.Error())
		return
	}
}
