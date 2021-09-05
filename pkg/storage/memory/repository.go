package memory

type Storage struct {
	Users map[UserId]User
	Rooms map[string]Room
}

func NewStorage() Storage {
	return Storage{}
}
