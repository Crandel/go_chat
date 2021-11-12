package chatting

type service struct {
	rooms    map[string]*Room
	commands chan Command
}

func NewService() *service {
	return &service{
		rooms:    make(map[string]*Room),
		commands: make(chan Command),
	}
}

func (s *service) NewUser(email string) {
	u := &User{
		Email:    email,
		commands: s.commands,
	}
	u.readInput()
}
