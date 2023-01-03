package chatting

import (
	lg "github.com/Crandel/go_chat/internal/logging"
)

// Room is the place for clients
type Room struct {
	Clients map[*Client]struct{}
	Name    string
}

func (r *Room) broadcast(sender *Client, message string) {
	for member := range r.Clients {
		if member != sender {
			log.Printf("Broadcast message %s \n", message)
			member.WriteMsg(message)
		}
	}
}

func (r *Room) addUser(c *Client) error {
	const op lg.Op = "chatting.Room.addUser"
	if r.hasUser(c) {
		return lg.New(op, "User already in room")
	}
	r.Clients[c] = struct{}{}
	return nil
}

func (r *Room) hasUser(c *Client) bool {
	if _, ok := r.Clients[c]; ok {
		return true
	}
	return false
}

func (r *Room) excludeFromRoom(c *Client) bool {
	if r.hasUser(c) {
		delete(r.Clients, c)
		return true
	}
	return false
}

type roomHandler struct {
	rooms map[string]*Room
}

func NewRoomHandler() *roomHandler {
	return &roomHandler{
		rooms: make(map[string]*Room, 0),
	}
}

func (rh *roomHandler) addRoom(roomName string) *Room {
	_, exists := rh.getRoom(roomName)
	if exists {
		return nil
	}
	room := &Room{
		Clients: make(map[*Client]struct{}),
		Name:    roomName,
	}
	rh.rooms[roomName] = room
	return room
}

func (rh *roomHandler) getRoom(roomName string) (*Room, bool) {
	if rh.rooms == nil {
		return nil, false
	}
	r, exists := rh.rooms[roomName]
	return r, exists
}

func (rh *roomHandler) roomHasUser(roomName string, c *Client) bool {
	r, exists := rh.getRoom(roomName)
	if !exists {
		return false
	}
	if r.hasUser(c) {
		return true
	}
	return false
}

func (rh *roomHandler) addUser(roomName string, c *Client) bool {
	r, ok := rh.getRoom(roomName)
	if !ok {
		r = rh.addRoom(roomName)
	}

	if r.hasUser(c) {
		return false
	}
	err := r.addUser(c)
	return err == nil
}

func (rh *roomHandler) excludeFromRoom(roomName string, c *Client) bool {
	r, exists := rh.getRoom(roomName)
	if !exists {
		return false
	}
	return r.excludeFromRoom(c)
}

func (rh *roomHandler) broadcast(roomName string, sender *Client, message string) {
	r, _ := rh.getRoom(roomName)
	r.broadcast(sender, message)
}

func (rh *roomHandler) listRooms() []string {
	names := make([]string, 0, len(rh.rooms))
	for name := range rh.rooms {
		names = append(names, name)
	}
	return names
}

func (rh *roomHandler) listUsers(c *Client) (string, []string) {
	for _, room := range rh.rooms {
		if room.hasUser(c) {
			clients := make([]string, 0, len(room.Clients))
			for client := range room.Clients {
				clients = append(clients, client.Nick)
			}
			return room.Name, clients
		}
	}
	return "", []string{}
}
