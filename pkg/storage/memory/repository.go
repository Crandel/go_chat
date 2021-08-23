package memory

type Storage struct {
	Users    map[string]User
	Rooms    map[string]Room
	Messages map[string]Message
}

func NewStorage() Storage {
	return Storage{}
}
