package memory

type Storage struct {
	Users    []User
	Rooms    []Room
	Messages []Message
}

func NewStorage() Storage {
	return Storage{}
}
