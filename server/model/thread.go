package model


type Thread struct {
	Threadname   string    `json:"threadname"`
	Tweets  []string    `json:"username"` // se guardan los ids de los tweets
}

type Threads []Thread