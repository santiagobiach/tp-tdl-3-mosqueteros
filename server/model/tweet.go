package model

import "time"

// import (
// 	"primitive"
// )

type Tweet struct {
	Idtweet   string    `json:"idtweet"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type Tweets []Tweet
