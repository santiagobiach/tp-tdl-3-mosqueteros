package model

import "time"

// import (
// 	"primitive"
// )

type Tweet struct {
	Idtweet   string    `json:"idtweet"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Liked_by  []string  `json:"liked_by"`
	Timestamp time.Time `json:"timestamp"`
}

type Tweets []Tweet
