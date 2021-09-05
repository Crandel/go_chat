package login

import "log"

type Repository interface {
	LoginUser(User) (string, error)
}

type Service interface {
	LoginUser(User) (string, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) LoginUser(u User) (string, error) {
	token, err := s.r.LoginUser(u)
	if err != nil {
		log.Println("Error while login:", err)
		return "", err
	}
	return token, nil
}
