package model

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Following []string `json:"following"`
	Followers []string `json:"followers"`
}

type Users []User
