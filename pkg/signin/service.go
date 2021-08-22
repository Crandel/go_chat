package signin

type Repository interface {
	SigninUser(User) (SigninResponse, error)
}

type Service interface {
	SigninUser(User) (SigninResponse, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) SigninUser(u User) (SigninResponse, error) {
	response, err := s.r.SigninUser(u)
	return response, err
}
