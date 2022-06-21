package auth

// LoginUser is for logging
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SigninUser is for signin
type SigninUser struct {
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
