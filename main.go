package main

import (
	"encoding/json"
	"fmt"

	"github.com/Crandel/go_chat/app"
)

var config = &configuration{}

// config struct
type configuration struct {
	Version  string       `json:"Version"`
	Database app.Database `json:"Database"`
	Server   app.Server   `json:"Server"`
	// Template Templates `json:"Template"`
	// Session  Session       `json:"Session"`
}

// ParseJSON ...
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

func main() {
	fmt.Println("Hi")
}
