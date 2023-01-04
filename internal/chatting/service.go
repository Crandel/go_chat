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
	NewClient(conn *websocket.Conn, nick string)
	Repository
}

type service struct {
	*roomHandler
	commands chan Command
	rep      Repository
}

func NewService(rep Repository) Service {
	roomHandler := NewRoomHandler()
	return &service{
		roomHandler: roomHandler,
		commands:    make(chan Command),
		rep:         rep,
	}
}

func (s *service) NewClient(conn *websocket.Conn, nick string) {
	const op = "chatting#NewClient:"
	u := &Client{
		Nick:     nick,
		conn:     conn,
		commands: s.commands,
	}
	log.Log(lg.Info, op, "New User was successfuly connected")
	u.ReadCommands()
}

func (s *service) Run() {
	const op = "chatting#Run "
	log.Log(lg.Debug, op, "Before loop")
	for command := range s.commands {
		log.Log(lg.Debug, op, "command ", command.id)
		switch command.id {
		case CmdMsg:
			log.Log(lg.Debug, op, "MSG ", s.rooms)
			for _, r := range s.rooms {
				if r.hasUser(command.client) {
					var msg strings.Builder
					finalMsg := strings.Join(command.args, " ")
					err := s.WriteMessage(command.client, r, finalMsg)
					if err == nil {
						msg.WriteString(command.client.Nick)
						msg.WriteString(" ")
						msg.WriteString(finalMsg)
						r.broadcast(command.client, msg.String())
					}
				}
			}
		case CmdPing:
			command.client.WriteMsg("pong")
		case CmdJoin:
			if len(command.args) != 2 {
				command.client.WriteMsg("Please provide correct room name")
				continue
			}

			roomName := command.args[1]
			exists, _ := s.RoomHasUser(roomName, command.client)
			if exists {
				continue
			} else {
				s.excludeFromRooms(command.client)
			}
			log.Log(lg.Debug, "Inside CmdJoin, user exists: ", exists)
			if err := s.AddUserToRoom(roomName, command.client); err != nil {
				command.client.WriteMsg("Something went wrong, err: " + err.Error())
				continue
			}
			command.client.WriteMsg("Welcome to the room " + roomName)
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
				r.broadcast(c, "User "+c.Nick+" leave the room")
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
	message := fmt.Sprintf("User %s join the room", c.Nick)

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
