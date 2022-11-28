package model

import "time"

// import (
// 	"primitive"
// )

type Topic struct {
	Idtopic   string    `json:"idtopic"`
	Topicstring   string    `json:"topicstring"`
	Tweets []string     `json:"username"` // ids de los tweets que hablan sobre este topic
	Timestamp time.Time `json:"timestamp"` // timestamp del último tweet sobre el topic
}

type topics []Topic
