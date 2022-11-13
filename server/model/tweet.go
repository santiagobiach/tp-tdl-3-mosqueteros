package model

type Tweet struct {
	Iduser  string `json:"iduser"`
	Content string `json:"content"`
}

type Tweets []Tweet
