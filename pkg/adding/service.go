package adding

type Repository interface {
	AddRoom(Room) (string, error)
}

type Service interface {
	AddRoom(Room) (string, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) AddRoom(r Room) (string, error) {
	return s.r.AddRoom(r)
}
