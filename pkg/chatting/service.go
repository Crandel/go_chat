package chatting

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type service struct {
	rooms    map[string]*Room
	commands chan Command
}

type Service interface {
	Run()
	NewUser(conn *websocket.Conn, nick string)
}

func NewService() *service {
	return &service{
		rooms:    make(map[string]*Room),
		commands: make(chan Command),
	}
}

func (s *service) NewUser(conn *websocket.Conn, nick string) {
	u := &User{
		Nick:     nick,
		conn:     conn,
		commands: s.commands,
	}
	fmt.Printf("User %s was successfuly connected", nick)
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
