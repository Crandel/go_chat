package reading

type UserId string

type User struct {
	Email      UserId `json:"email"`
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
}
