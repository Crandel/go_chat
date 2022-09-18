package reading

type UserId string

type User struct {
	Nick       UserId `json:"nick"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
}
