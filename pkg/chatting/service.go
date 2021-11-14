package chatting

import (
	"fmt"
	"log"

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
		commands: s.commands,
	}
	fmt.Println("Client was successfuly connected")
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))
		c := Command{
			id:   CMD_MSG,
			user: u,
		}
		u.commands <- c
	}
}
