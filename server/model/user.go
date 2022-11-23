package model

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Following []string `json:"following"`
	Followers []string `json:"followers"`
	Threads []string `json:"threads"` // se guardan los nombres de los threads
}

type Users []User
