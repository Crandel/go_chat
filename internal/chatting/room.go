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
	if r.haveUser(c) {
		return lg.New(op, lg.Info, "User already in room")
	}
	r.Clients[c] = struct{}{}
	return nil
}

func (r *Room) haveUser(c *Client) bool {
	if _, ok := r.Clients[c]; ok {
		return true
	}
	return false
}

func (r *Room) excludeFromRoom(c *Client) bool {
	if r.haveUser(c) {
		delete(r.Clients, c)
		return true
	}
	return false
}

type roomHandler struct {
	rooms map[string]*Room
}

func (rh *roomHandler) addRoom(roomName string) bool {
	_, exists := rh.getRoom(roomName)
	if exists {
		return false
	}
	room := Room{
		Clients: make(map[*Client]struct{}),
		Name:    roomName,
	}
	rh.rooms[roomName] = &room
	return true
}

func (rh *roomHandler) getRoom(roomName string) (*Room, bool) {
	r, exists := rh.rooms[roomName]
	return r, exists
}

func (rh *roomHandler) roomHasUser(roomName string, c *Client) bool {
	r, exists := rh.getRoom(roomName)
	if !exists {
		return false
	}
	if r.haveUser(c) {
		return true
	}
	return false
}

func (rh *roomHandler) addUser(roomName string, c *Client) bool {
	r, ok := rh.getRoom(roomName)
	if !ok {
		rh.addRoom(roomName)
	}

	if rh.roomHasUser(roomName, c) {
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
		if room.haveUser(c) {
			clients := make([]string, 0, len(room.Clients))
			for client := range room.Clients {
				clients = append(clients, client.Nick)
			}
			return room.Name, clients
		}
	}
	return "", []string{}
}
