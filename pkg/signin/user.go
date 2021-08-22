package signin

type User struct {
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type SigninResponse struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}
