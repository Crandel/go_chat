package adding

type Repository interface {
	AddRoom(string) (string, error)
}

type Service interface {
	AddRoom(string) (string, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) AddRoom(rn string) (string, error) {
	return s.r.AddRoom(rn)
}
