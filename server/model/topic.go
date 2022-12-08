package model

// import (
// 	"primitive"
// )

type Topic struct {
	Topicstring string  `json:"topicstring"`
	Tweets      []Tweet `json:"tweets"` // ids de los tweets que hablan sobre este topic
}

type Topics []Topic
