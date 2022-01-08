package chatting

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type Repository interface {
	WriteMessage(u Client, r Room, msg string) error
	ExcludeFromRoom(name string, u Client) error
	AddUserToRoom(name string, u Client) error
	RoomHasUser(name string, u Client) bool
}

type Service interface {
	Run()
	NewUser(conn *websocket.Conn, nick string)
	WriteMessage(u Client, r Room, msg string) error
	ExcludeFromRoom(name string, u Client) error
	AddUserToRoom(name string, u Client) error
	RoomHasUser(name string, u Client) bool
}

type service struct {
	rooms    map[string]*Room
	commands chan Command
	rep      Repository
}

func NewService(rep Repository) *service {
	return &service{
		rooms:    make(map[string]*Room),
		commands: make(chan Command),
		rep:      rep,
	}
}

func (s *service) NewUser(conn *websocket.Conn, nick string) {
	u := &Client{
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
				if r.haveUser(c.client) {
					var msg strings.Builder
					msg.WriteString(c.client.Nick)
					msg.WriteString(" ")
					msg.WriteString(strings.Join(c.args, " "))
					fmt.Printf("%s\n", msg.String())

					r.broadcast(c.client, msg.String())
				}
			}
		case CMD_PING:
			c.client.WriteMsg("pong")

		case CMD_JOIN:
			if len(c.args) > 2 {
				c.client.WriteMsg("Please provide only correct room name")
				continue
			}
			roomName := c.args[1]
			room, exists := s.rooms[roomName]
			if exists {
				if room.haveUser(c.client) {
					c.client.WriteMsg("You are in room " + roomName)
					continue
				}
			} else {
				room = &Room{
					Name:    roomName,
					Clients: nil,
				}
				s.rooms[roomName] = room
			}
			s.excludeFromRooms(c.client)
			room.addUser(c.client)
			room.broadcast(c.client, fmt.Sprintf("User %s join the room", c.client.Nick))
			c.client.WriteMsg("Welcome to the room " + room.Name)
		case CMD_ROOMS:
			names := make([]string, 0, len(s.rooms))
			for name := range s.rooms {
				names = append(names, name)
			}
			c.client.WriteMsg("Rooms:\n" + strings.Join(names, "\n"))
		case CMD_USERS:
			for _, room := range s.rooms {
				if room.haveUser(c.client) {
					clients := make([]string, 0, len(room.Clients))
					for client := range room.Clients {
						clients = append(clients, client.Nick)
					}
					var msg strings.Builder
					msg.WriteString("Users in room ")
					msg.WriteString(room.Name)
					msg.WriteString("\n")
					msg.WriteString(strings.Join(clients, "\n"))
					c.client.WriteMsg(msg.String())
				}
			}
		case CMD_QUIT:
			s.excludeFromRooms(c.client)
		}
	}
}

func (s *service) excludeFromRooms(u *Client) {
	for _, r := range s.rooms {
		if r.haveUser(u) {
			delete(r.Clients, u)
			r.broadcast(u, "User "+u.Nick+" leave the room")
		}
	}
}

func (s *service) WriteMessage(u Client, r Room, msg string) error {
	return s.rep.WriteMessage(u, r, msg)
}

func (s *service) ExcludeFromRoom(roomName string, u Client) error {
	return s.rep.ExcludeFromRoom(roomName, u)
}

func (s *service) AddUserToRoom(roomName string, u Client) error {
	return s.rep.AddUserToRoom(roomName, u)
}

func (s *service) RoomHasUser(roomName string, u Client) bool {
	return s.rep.RoomHasUser(roomName, u)
}
