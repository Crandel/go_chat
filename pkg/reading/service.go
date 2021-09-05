package reading

type Repository interface {
	ReadUsers() ([]User, error)
	ReadUser(UserId) (User, error)
	ReadRooms() ([]Room, error)
	ReadRoom(string) (Room, error)
}

type Service interface {
	ReadUsers() ([]User, error)
	ReadUser(UserId) (User, error)
	ReadRooms() ([]Room, error)
	ReadRoom(string) (Room, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) ReadUsers() ([]User, error) {
	return s.r.ReadUsers()
}

func (s *service) ReadUser(ui UserId) (User, error) {
	return s.r.ReadUser(ui)
}

func (s *service) ReadRooms() ([]Room, error) {
	return s.r.ReadRooms()
}

func (s *service) ReadRoom(ri string) (Room, error) {
	return s.r.ReadRoom(ri)
}
