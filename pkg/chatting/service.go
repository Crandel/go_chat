package chatting

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type Repository interface {
	WriteMessage(u User, r Room, msg string) error
	ExcludeFromRoom(name string, u User) error
	AddUserToRoom(name string, u User) error
	RoomHasUser(name string, u User) bool
}

type Service interface {
	Run()
	NewUser(conn *websocket.Conn, nick string)
	WriteMessage(u User, r Room, msg string) error
	ExcludeFromRoom(name string, u User) error
	AddUserToRoom(name string, u User) error
	RoomHasUser(name string, u User) bool
}

type service struct {
	rooms    map[string]*Room
	commands chan Command
	r        Repository
}

func NewService(r Repository) *service {
	return &service{
		rooms:    make(map[string]*Room),
		commands: make(chan Command),
		r:        r,
	}
}

func (s *service) NewUser(conn *websocket.Conn, nick string) {
	u := &User{
		Nick:     nick,
		conn:     conn,
		commands: s.commands,
	}
	fmt.Printf("chatting#NewUser User %s was successfuly connected\n", nick)
	u.ReadCommands()
}

func (s *service) Run() {
	fmt.Println("chatting#Run Before loop")
	for c := range s.commands {
		fmt.Println("chatting#Run#command " + c.id)
		switch c.id {
		case CMD_MSG:
			fmt.Println("chatting#Run#command#MSG ")
			for _, r := range s.rooms {
				if r.haveUser(c.user) {
					var msg strings.Builder
					msg.WriteString(c.user.Nick)
					msg.WriteString(" ")
					msg.WriteString(strings.Join(c.args, " "))
					fmt.Printf("%s\n", msg.String())

					r.broadcast(c.user, msg.String())
				}
			}
		case CMD_PING:
			c.user.WriteMsg("pong")

		case CMD_JOIN:
			if len(c.args) > 2 {
				c.user.WriteMsg("Please provide only correct room name")
				continue
			}
			roomName := c.args[1]
			r, exists := s.rooms[roomName]
			if exists {
				if r.haveUser(c.user) {
					c.user.WriteMsg("You are in room " + roomName)
					continue
				}
			} else {
				r = &Room{
					Name:    roomName,
					Members: make(map[string]*User),
				}
				s.rooms[roomName] = r
			}
			s.excludeFromRooms(c.user)
			r.addUser(c.user)
			r.broadcast(c.user, fmt.Sprintf("User %s join the room", c.user.Nick))
			c.user.WriteMsg("Welcome to the room " + r.Name)
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
			s.excludeFromRooms(c.user)
		}
	}
}

func (s *service) excludeFromRooms(u *User) {
	for _, r := range s.rooms {
		if r.haveUser(u) {
			delete(r.Members, u.Nick)
			r.broadcast(u, "User "+u.Nick+" leave the room")
		}
	}
}

func (s *service) WriteMessage(u User, r Room, msg string) error {
	return s.r.WriteMessage(u, r, msg)
}

func (s *service) ExcludeFromRoom(roomName string, u User) error {
	return s.r.ExcludeFromRoom(roomName, u)
}

func (s *service) AddUserToRoom(roomName string, u User) error {
	return s.r.AddUserToRoom(roomName, u)
}

func (s *service) RoomHasUser(roomName string, u User) bool {
	return s.r.RoomHasUser(roomName, u)
}
