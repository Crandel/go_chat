package chatting

import (
	"fmt"
	"log"
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
	NewClient(conn *websocket.Conn)
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

func NewService(rep Repository) Service {
	return &service{
		rooms:    make(map[string]*Room),
		commands: make(chan Command),
		rep:      rep,
	}
}

func (s *service) NewClient(conn *websocket.Conn) {
	u := &Client{
		Nick:     nil,
		conn:     conn,
		commands: s.commands,
	}
	log.Printf("chatting#NewClient: New User was successfuly connected\n")
	u.ReadCommands()
}

func (s *service) Run() {
	log.SetPrefix("chatting#Run ")
	log.Println("Before loop")
	for c := range s.commands {
		if c.client.Nick == nil && c.id != CmdJoin {
			c.client.WriteMsg("Please provide join room and specify your name")
			continue
		}
		log.Println("command " + c.id)
		switch c.id {
		case CmdMsg:
			log.Println("MSG ")
			for _, r := range s.rooms {
				if r.haveUser(c.client) {
					var msg strings.Builder
					msg.WriteString(*c.client.Nick)
					msg.WriteString(" ")
					msg.WriteString(strings.Join(c.args, " "))
					r.broadcast(c.client, msg.String())
				}
			}
		case CmdPing:
			c.client.WriteMsg("pong")

		case CmdJoin:

			if len(c.args) != 3 {
				c.client.WriteMsg("Please provide correct room and user name")
				continue
			}

			roomName := c.args[1]
			userName := c.args[2]
			c.client.Nick = &userName
			room, exists := s.rooms[roomName]
			if exists {
				if room.haveUser(c.client) {
					c.client.WriteMsg("You are in room " + roomName)
					continue
				}
			} else {
				room = &Room{
					Name:    roomName,
					Clients: make(map[*Client]struct{}),
				}
				s.rooms[roomName] = room
			}
			s.excludeFromRooms(c.client)
			room.addUser(c.client)
			room.broadcast(c.client, fmt.Sprintf("User %s join the room", *c.client.Nick))
			c.client.WriteMsg("Welcome to the room " + room.Name)
		case CmdRooms:
			names := make([]string, 0, len(s.rooms))
			for name := range s.rooms {
				names = append(names, name)
			}
			c.client.WriteMsg("Rooms:\n" + strings.Join(names, "\n"))
		case CmdUsers:
			for _, room := range s.rooms {
				if room.haveUser(c.client) {
					clients := make([]string, 0, len(room.Clients))
					for client := range room.Clients {
						clients = append(clients, *client.Nick)
					}
					var msg strings.Builder
					msg.WriteString("Users in room ")
					msg.WriteString(room.Name)
					msg.WriteString("\n")
					msg.WriteString(strings.Join(clients, "\n"))
					c.client.WriteMsg(msg.String())
				}
			}
		case CmdQuit:
			s.excludeFromRooms(c.client)
		}
	}
}

func (s *service) excludeFromRooms(u *Client) {
	for _, r := range s.rooms {
		if r.haveUser(u) {
			delete(r.Clients, u)
			r.broadcast(u, "User "+*u.Nick+" leave the room")
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
