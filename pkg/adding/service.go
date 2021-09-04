package adding

type Repository interface {
	AddMessage(Message) (string, error)
}

type Service interface {
	AddMessage(Message) (string, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) AddMessage(m Message) (string, error) {
	return s.r.AddMessage(m)
}
