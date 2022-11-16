package model
// import (
// 	"primitive"
// )

type Tweet struct {
	Idtweet  string `json:"idtweet"`
	Username string `json:"username"`
	Content string `json:"content"`
}

type Tweets []Tweet
