package server_utils

import (
	"fmt"
	"log"
	"net"
	"server/database"
	"server/model"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
func HandleLogin(c net.Conn, arguments []string, username *string) {

	if *username != "" {
		msg := "No puedes iniciar sesión si ya hay alguien logueado" // mensaje de error
		_, _ = c.Write([]byte(msg))
		return;
	}

	fmt.Println("Voy a handlear un login")
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("users")

	var user model.User
	filter := bson.D{
		{"username", arguments[1]},
	}

	err = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, _ = c.Write([]byte("No existe la cuenta"))
			return
		}
		log.Fatal(err)
	}
	if user.Password != arguments[2] {
		_, _ = c.Write([]byte("El password es incorrecto"))
		return
	}
	*username = user.Username
	msg := "holà " + user.Username // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))
}

// User signup
func HandleSignup(c net.Conn, arguments []string, username *string) {
	if *username != "" {
		msg := "No puedes crear una cuenta si hay alguien logueado" // mensaje de error
		_, _ = c.Write([]byte(msg))
		return;
	}
	fmt.Println("Voy a handlear un signup")

	//En los arguments tambien esta el comando
	if arguments[2] != arguments[3] {
		_, _ = c.Write([]byte("Los password no son iguales"))
		return
	}
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("users")

	var user model.User
	user.Username = arguments[1]
	user.Password = arguments[2]
	user.Following = []string{}
	user.Followers = []string{}
	//user.Isonline = true
	//Primero fijarse si ya existe
	filter := bson.D{
		{"username", arguments[1]},
	}
	err = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			insertResult, err := coll.InsertOne(ctx, user)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Inserted a single document: ", insertResult.InsertedID)

			msg := "cuenta creada" // mensaje de login exitoso
			_, _ = c.Write([]byte(msg))
			return
		}
		log.Fatal(err)
	}
	_, _ = c.Write([]byte("Ya existe la cuenta"))

}

func HandleTweet(c net.Conn, arguments []string, username *string) {

	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		_, _ = c.Write([]byte(msg))
		return;
	}


	fmt.Println("Voy a handlear un tweet")

	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	//coll_user := client.Database("tdl-los-tres-mosqueteros").Collection("users")
	coll_tweet := client.Database("tdl-los-tres-mosqueteros").Collection("tweets")

	var tweet model.Tweet
	var tweet_content string

	for i := 1; i < len(arguments); i++ {
		tweet_content += arguments[i] + " "
	}

	tweet.Content = tweet_content
	tweet.Username = *username

	insertTweet, err := coll_tweet.InsertOne(ctx, tweet)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("hay un nuevo tweet, on id : ", insertTweet.InsertedID)
	fmt.Println("El nuevo tweet es de ", tweet.Username)

	msg := "tweet enviado" // mensaje de tweet exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleFollow(c net.Conn, arguments []string, username *string) {
	fmt.Println("Voy a handlear un follow")
	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		_, _ = c.Write([]byte(msg))
		return;
	}

	// Primero, busco al usuario que quiero seguir solo para saber si existe (no podes seguir a alguien que no existe)
	// Despues, busco mi usuario (username)

	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("users")

	var userToFollow model.User //usuario que quiero seguir
	filter := bson.D{
		{"username", arguments[1]},
	}

	err = coll.FindOne(ctx, filter).Decode(&userToFollow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, _ = c.Write([]byte("No existe la cuenta que quieres seguir"))
			return
		}
		log.Fatal(err)
	}
	userToFollow.Followers = append(userToFollow.Followers, *username)

	var user model.User //yo
	// filter := bson.D{
	// 	{"username", username},
	// }

	_ = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, _ = c.Write([]byte("No existe la cuenta que quieres seguir"))
			return
		}
		log.Fatal(err)
	}
	user.Following = append(user.Following, arguments[1])
	fmt.Println(user)
	fmt.Println(userToFollow)
	msg := "¡Has seguido a " + arguments[1] + " !" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
func HandleUnfollow(c net.Conn, arguments []string, username *string) {
	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		_, _ = c.Write([]byte(msg))
		return;
	}
	fmt.Println("Voy a handlear un unf")
	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))
}
func HandleTweetsFrom(c net.Conn, arguments []string, username *string) {

	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		_, _ = c.Write([]byte(msg))
		return;
	}
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

func ParseMessage(c net.Conn, message string, username *string) {
	split_message := strings.SplitAfter(message, " ")
	for i, v := range split_message {
		split_message[i] = strings.TrimSpace(v)
	}
	fmt.Println(split_message)

	// TODAS DEBERIAN RECIBIR EL USERNAME

	switch strings.TrimSuffix(split_message[0], " ") { // Hay que hacer eso sí o sí porque deja un " " de más
	case Login:
		HandleLogin(c, split_message, username)
	case Signup:
		HandleSignup(c, split_message, username)
	case Tweet:
		HandleTweet(c, split_message, username)
	case Follow:
		HandleFollow(c, split_message, username)
	case Unfollow:
		HandleUnfollow(c, split_message, username)
	case TweetsFrom:
		HandleTweetsFrom(c, split_message, username)
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
