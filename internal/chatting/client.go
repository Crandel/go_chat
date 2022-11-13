package chatting

import (
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	Nick     *string
	conn     *websocket.Conn
	commands chan<- Command
}

func (c *Client) GetNick() string {
	if c.Nick == nil {
		return *c.Nick
	}
	return ""
}

func (c *Client) ReadCommands() {
	log.SetPrefix("chatting#user#ReadCommands ")
	defer c.conn.Close()
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			c.commands <- Command{
				id:     CmdQuit,
				client: c,
			}
			return
		}
		rawCommand := string(p)
		args := strings.Split(rawCommand, " ")
		cmd := strings.TrimSpace(args[0])
		log.Println("Command: " + cmd)
		var cmdID CommandID
		if !strings.HasPrefix(cmd, "/") {
			cmdID = CmdMsg
		} else {
			switch cmd {
			case CmdPing:
				cmdID = CmdPing
			case CmdJoin:
				cmdID = CmdJoin
			case CmdUsers:
				cmdID = CmdUsers
			case CmdRooms:
				cmdID = CmdRooms
			case CmdQuit:
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
		log.Println("WriteMsg " + err.Error())
		return
	}
}
