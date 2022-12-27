package chatting

import (
	"fmt"
	"strings"

	lg "github.com/Crandel/go_chat/internal/logging"
	"github.com/gorilla/websocket"
)

var log = lg.InitLogger()

type Repository interface {
	WriteMessage(c *Client, r *Room, msg string) error
	ExcludeFromRoom(name string, c *Client) error
	AddUserToRoom(name string, c *Client) error
	RoomHasUser(name string, c *Client) (bool, int)
}

type Service interface {
	Run()
	NewClient(conn *websocket.Conn)
	Repository
}

type service struct {
	roomHandler
	commands chan Command
	rep      Repository
}

func NewService(rep Repository) Service {
	return &service{
		roomHandler: roomHandler{
			rooms: make(map[string]*Room),
		},
		commands: make(chan Command),
		rep:      rep,
	}
}

func (s *service) NewClient(conn *websocket.Conn) {
	log.SetPrefix("chatting#NewClient: ")
	u := &Client{
		Nick:     nil,
		conn:     conn,
		commands: s.commands,
	}
	log.Debugln("New User was successfuly connected")
	u.ReadCommands()
}

func (s *service) Run() {
	log.SetPrefix("chatting#Run ")
	log.Debugln("Before loop")
	for command := range s.commands {
		if command.client.Nick == nil && command.id != CmdJoin {
			command.client.WriteMsg("Please provide join room and specify your name")
			continue
		}
		log.Debugln("command " + command.id)
		switch command.id {
		case CmdMsg:
			log.Debugln("MSG ")
			for _, r := range s.rooms {
				if r.haveUser(command.client) {
					var msg strings.Builder
					finalMsg := strings.Join(command.args, " ")
					err := s.WriteMessage(command.client, r, finalMsg)
					if err == nil {
						msg.WriteString(*command.client.Nick)
						msg.WriteString(" ")
						msg.WriteString(finalMsg)
						r.broadcast(command.client, msg.String())
					}
				}
			}
		case CmdPing:
			command.client.WriteMsg("pong")
		case CmdJoin:
			if len(command.args) != 3 {
				command.client.WriteMsg("Please provide correct room and user name")
				continue
			}

			roomName := command.args[1]
			userName := command.args[2]
			command.client.Nick = &userName
			exists, _ := s.RoomHasUser(roomName, command.client)
			if exists {
				continue
			} else {
				s.excludeFromRooms(command.client)
			}
			if err := s.AddUserToRoom(roomName, command.client); err == nil {
				command.client.WriteMsg("Welcome to the room " + roomName)
			}
		case CmdRooms:
			names := s.roomHandler.listRooms()
			command.client.WriteMsg("Rooms:\n" + strings.Join(names, "\n"))
		case CmdUsers:
			roomName, clients := s.roomHandler.listUsers(command.client)
			var msg strings.Builder
			msg.WriteString("Users in room ")
			msg.WriteString(roomName)
			msg.WriteString("\n")
			msg.WriteString(strings.Join(clients, "\n"))
			command.client.WriteMsg(msg.String())
		case CmdQuit:
			s.excludeFromRooms(command.client)
		}
	}
}

func (s *service) excludeFromRooms(c *Client) {
	for _, r := range s.rooms {
		done := s.roomHandler.excludeFromRoom(r.Name, c)
		if done {
			err := s.ExcludeFromRoom(r.Name, c)
			if err != nil {
				r.broadcast(c, "User "+*c.Nick+" leave the room")
			}
		}
	}
}

func (s *service) WriteMessage(c *Client, r *Room, msg string) error {
	err := s.rep.WriteMessage(c, r, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ExcludeFromRoom(roomName string, c *Client) error {
	err := s.rep.ExcludeFromRoom(roomName, c)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddUserToRoom(roomName string, c *Client) error {
	err := s.rep.AddUserToRoom(roomName, c)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("User %s join the room", *c.Nick)

	if done := s.roomHandler.addUser(roomName, c); done {
		s.roomHandler.broadcast(roomName, c, message)
	}

	return nil
}

func (s *service) RoomHasUser(roomName string, c *Client) (bool, int) {
	if done := s.roomHandler.roomHasUser(roomName, c); done {
		exists, id := s.rep.RoomHasUser(roomName, c)
		if exists {
			return true, id
		}
	}
	return false, 0
}
