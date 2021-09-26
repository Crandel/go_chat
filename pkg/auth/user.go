package auth

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninUser struct {
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
