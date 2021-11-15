package chatting

import (
	"fmt"
	"strings"

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
	for c := range s.commands {
		switch c.id {
		case CMD_MSG:
			for _, r := range s.rooms {
				if r.haveUser(c.user) {
					var msg strings.Builder
					msg.WriteString(c.user.Nick)
					msg.WriteString(" ")
					msg.WriteString(strings.Join(c.args, " "))
					r.broadcast(msg.String())
				}
			}
		case CMD_JOIN:
			if len(c.args) > 2 {
				c.user.WriteMsg("Please provide only correct room name")
			}
			roomName := c.args[1]
			r, ok := s.rooms[roomName]
			if ok {
				if r.haveUser(c.user) {
					c.user.WriteMsg("You are in room " + roomName)
				} else {
					r.addUser(c.user)
				}
			} else {
				r.addUser(c.user)
				s.rooms[roomName] = r
			}
		case CMD_ROOMS:
			names := make([]string, 0, len(s.rooms))
			for name := range s.rooms {
				names = append(names, name)
			}
			c.user.WriteMsg("Rooms:\n" + strings.Join(names, "\n"))
		case CMD_USERS:
			for _, r := range s.rooms {
				if r.haveUser(c.user) {
					mbrs := make([]string, 0, len(r.Members))
					for member := range r.Members {
						mbrs = append(mbrs, member)
					}
					var msg strings.Builder
					msg.WriteString("Users in room ")
					msg.WriteString(r.Name)
					msg.WriteString("\n")
					msg.WriteString(strings.Join(mbrs, "\n"))
					c.user.WriteMsg(msg.String())
				}
			}
		case CMD_QUIT:
			for _, r := range s.rooms {
				if r.haveUser(c.user) {
					delete(r.Members, c.user.Nick)
					r.broadcast("User " + c.user.Nick + " leave the room")
				}
			}
		}
	}
}
