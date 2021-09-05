package reading

type UserId string

type User struct {
	ID         UserId `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
}
