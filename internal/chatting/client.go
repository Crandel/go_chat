package chatting

import (
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
		var message ChatMessage
		err := c.conn.ReadJSON(&message)
		if err != nil {
			log.Log(lg.Warning, op, err.Error())
			c.commands <- Command{
				id:     CmdQuit,
				client: c,
			}
			return
		}
		c.commands <- Command{
			id:     message.CommandId,
			client: c,
			args:   message.Args,
		}
	}
}

func (c *Client) WriteMsg(message string, userName ...string) {
	var user *string
	if len(userName) > 0 {
		user = &userName[0]
	}
	chatMessage := ChatMessage{
		CmdMsg,
		user,
		[]string{
			message,
		},
	}
	if err := c.conn.WriteJSON(chatMessage); err != nil {
		log.Log(lg.Warning, "WriteMsg ", err.Error())
		return
	}
}
