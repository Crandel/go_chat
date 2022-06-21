package adding

// Repository interface define functions for handling Rooms
type Repository interface {
	AddRoom(string) (string, error)
}

// Service interface proxy functions from Repository for handling Rooms
type Service interface {
	AddRoom(string) (string, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddRoom(rn string) (string, error) {
	return s.r.AddRoom(rn)
}
