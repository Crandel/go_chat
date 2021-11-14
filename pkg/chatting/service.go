package chatting

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type service struct {
	rooms    map[string]*Room
	commands chan Command
}

func NewService() *service {
	return &service{
		rooms:    make(map[string]*Room),
		commands: make(chan Command),
	}
}

func (s *service) NewUser(conn websocket.Conn, email string) {
	u := &User{
		Email:    email,
		conn:     conn,
		commands: s.commands,
	}
	fmt.Println("Client was successfuly connected")
	u.ReadCommands()
}

func (s *service) Run() {
	for {
		c := <-s.commands
		switch c.id {
		case CMD_MSG:

		}
	}
}
