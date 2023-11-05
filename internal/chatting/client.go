package chatting

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	commands chan<- Command
	Nick     string
}

func (c *Client) ReadCommands() {
	logRead := log.With(slog.Group("ReadCommands"))
	defer c.conn.Close()
	for {
		var message ChatMessage
		err := c.conn.ReadJSON(&message)
		if err != nil {
			logRead.Warn(err.Error())
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
	logWrite := log.With(slog.Group("WriteMsg"))
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
		logWrite.Warn(err.Error())
		return
	}
}
