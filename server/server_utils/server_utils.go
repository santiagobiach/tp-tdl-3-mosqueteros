package server_utils

import (
	"fmt"
	"net"
	"strings"
)

const (
	Login              string = "login"
	Signup             string = "signup"
	Tweet              string = "tweet"
	Follow             string = "follow"
	Unfollow           string = "unfollow"
	TweetsFrom         string = "tweetsFrom"
	TrendingTopic      string = "trendingTopic"
	TrendingTweetsFrom string = "trendingTweetsFrom"
	MyTweets           string = "myTweets"
	MyFollowers        string = "myFollowers"
	MyFollowing        string = "myFollowing"
	Feed               string = "feed"
	Reply              string = "reply"
	NewThread          string = "newThread"
	AddTweetToThread   string = "addTweetToThread"
	Thread             string = "thread"
	Like               string = "like"
	MostLiked          string = "mostLiked" //puede ser de 2 tipos
	MostFollowed       string = "mostFollowed"
)

// User login
func HandleLogin(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un login")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}

// User signup
func HandleSignup(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un signup")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}

func HandleTweet(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un tweet")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleFollow(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un follow")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleUnfollow(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un unf")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleTweetsFrom(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un tweet from")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleTrendingTopic(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un tt")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleTrendingTweetsFrom(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un ttfrom")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleMyTweets(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un my tweets")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleMyFollowers(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un myfollowers")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleMyFollowing(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un myfollowing")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}

func HandleFeed(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un feed")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleReply(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un reply")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}

func HandleAddTweetToThread(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un addtweedtothread")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}

func HandleNewThread(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un newthread")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleThread(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un thread")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}

func HandleLike(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un like")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleMostLiked(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un mostliked")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}

func HandleMostFollowed(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un mostfollowed")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}

func ParseMessage(c net.Conn, message string) {
	split_message := strings.SplitAfter(message, " ")
	//fmt.Println(strings.TrimSuffix("Foo++", "+"))
	switch strings.TrimSuffix(split_message[0], " ") { // Hay que hacer eso sí o sí porque deja un " " de más
	case Login:
		HandleLogin(c, split_message)
	case Signup:
		HandleSignup(c, split_message)
	case Tweet:
		HandleTweet(c, split_message)
	case Follow:
		HandleFollow(c, split_message)
	case Unfollow:
		HandleUnfollow(c, split_message)
	case TweetsFrom:
		HandleTweetsFrom(c, split_message)
	case TrendingTopic:
		HandleTrendingTopic(c, split_message)
	case TrendingTweetsFrom:
		HandleTrendingTweetsFrom(c, split_message)
	case MyTweets:
		HandleMyTweets(c, split_message)
	case MyFollowers:
		HandleMyFollowers(c, split_message)
	case MyFollowing:
		HandleMyFollowing(c, split_message)
	case Feed:
		HandleFeed(c, split_message)
	case Reply:
		HandleReply(c, split_message)
	case NewThread:
		HandleNewThread(c, split_message)
	case AddTweetToThread:
		HandleAddTweetToThread(c, split_message)
	case Thread:
		HandleThread(c, split_message)
	case Like:
		HandleLike(c, split_message)
	case MostLiked:
		HandleMostLiked(c, split_message)
	case MostFollowed:
		HandleMostFollowed(c, split_message)
	default:
		msg := "Error" // comando no existe
		_, _ = c.Write([]byte(msg))
	}
}
