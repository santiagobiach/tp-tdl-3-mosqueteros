package server_utils

import (
	"context"
	"fmt"
	"log"
	"net"
	"server/database"
	"server/model"
	"strconv"
	"strings"
	"time"

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
		msg := "No puedes iniciar sesión si ya hay alguien logueado " // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
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
			fmt.Fprintf(c, "No existe la cuenta \n")
			return
		}
		log.Fatal(err)
	}
	if user.Password != arguments[2] {
		fmt.Fprintf(c, "El password es incorrecto \n")
		return
	}
	*username = user.Username
	msg := "holà " + user.Username // mensaje de login exitoso
	fmt.Fprintf(c, msg+"\n")
}

// User signup
func HandleSignup(c net.Conn, arguments []string, username *string) {
	if *username != "" {
		msg := "No puedes crear una cuenta si hay alguien logueado" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}
	fmt.Println("Voy a handlear un signup")

	//En los arguments tambien esta el comando
	if arguments[2] != arguments[3] {
		fmt.Fprintf(c, "Las password no son iguales\n")
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
			fmt.Fprintf(c, msg+"\n")
			return
		}
		log.Fatal(err)
	}
	fmt.Fprintf(c, "Ya existe la cuenta \n")

}

func HandleTweet(c net.Conn, arguments []string, username *string) {

	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
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
		//Si la palabra empieza con # agregar a una coleccion de la bdd de topics
		tweet_content += arguments[i] + " "
	}

	tweet.Content = tweet_content
	tweet.Username = *username
	tweet.Timestamp = time.Now()
	insertTweet, err := coll_tweet.InsertOne(ctx, tweet)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("hay un nuevo tweet, on id : ", insertTweet.InsertedID)
	fmt.Println("El nuevo tweet es de ", tweet.Username)

	msg := "tweet enviado" // mensaje de tweet exitoso
	fmt.Fprintf(c, msg+"\n")

}
func HandleFollow(c net.Conn, arguments []string, username *string) {
	fmt.Println("Voy a handlear un follow")
	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
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
			fmt.Fprintf(c, "No existe la cuenta que quieres seguir\n")
			return
		}
		log.Fatal(err)
	}
	userToFollow.Followers = append(userToFollow.Followers, *username)
	_, err = coll.ReplaceOne(ctx, filter, userToFollow)
	var user model.User //yo
	filter = bson.D{
		{"username", *username},
	}

	_ = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe la cuenta que quieres seguir\n")
			return
		}
		log.Fatal(err)
	}
	user.Following = append(user.Following, arguments[1])
	_, err = coll.ReplaceOne(ctx, filter, user)
	fmt.Println(user)
	fmt.Println(userToFollow)
	msg := "¡Has seguido a " + arguments[1] + " !" // mensaje de login exitoso
	fmt.Fprintf(c, msg+"\n")

}
func HandleUnfollow(c net.Conn, arguments []string, username *string) {
	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}

	fmt.Println("Voy a handlear un unf")

	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("users")

	var userToUnfollow model.User //usuario que quiero seguir
	filter := bson.D{
		{"username", arguments[1]},
	}

	err = coll.FindOne(ctx, filter).Decode(&userToUnfollow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe la cuenta a la que quieres hacer unfollow\n")
			return
		}
		log.Fatal(err)
	}

	//aca eliminamos al usuario de la lista de suscriptores del que quiere dejar de seguir
	for i, j := range userToUnfollow.Followers {
		if j == *username {
			userToUnfollow.Followers = append(userToUnfollow.Followers[:i], userToUnfollow.Followers[(i+1):]...)
		}
	}
	_, err = coll.ReplaceOne(ctx, filter, userToUnfollow)
	var user model.User //yo
	filter = bson.D{
		{"username", *username},
	}
	_ = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe la cuenta a la que quieres hacer unfollow\n")
			return
		}
		log.Fatal(err)
	}

	//aca eliminamos al que queremos dejar de seguir de la lista de suscriptores del usuario
	for i, j := range user.Following {
		if j == arguments[1] {
			user.Following = append(user.Following[:i], user.Following[(i+1):]...)
		}
	}
	_, err = coll.ReplaceOne(ctx, filter, user)
	fmt.Println(user)
	fmt.Println(userToUnfollow)
	msg := "¡Has unfollow " + arguments[1] + " !" // mensaje de login exitoso
	fmt.Fprintf(c, msg+"\n")
}

func HandleTweetsFrom(c net.Conn, arguments []string, username *string) {

	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}
	fmt.Println("Voy a handlear un tweet from")

	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("tweets")
	days, err := strconv.Atoi(arguments[2])
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now().UTC()
	aux := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	filter := bson.D{
		{"username", arguments[1]},
		{"timestamp", bson.M{
			"$gte": aux.AddDate(0, 0, -days),
		}},
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No se encontraron tweets\n")
			return
		}
		log.Fatal(err)
	}
	var results model.Tweets
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	msg := ""
	for _, result := range results {
		cursor.Decode(&result)

		msg = result.Idtweet + result.Content + result.Timestamp.String()
		fmt.Fprintf(c, msg+"\n")
	}

}
func HandleTrendingTopic(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un tt")

}

func HandleTrendingTweetsFrom(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un ttfrom")

}
func HandleMyTweets(c net.Conn, arguments []string, username *string) {

	fmt.Println("Voy a handlear un my tweets")

	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("tweets")

	filter := bson.D{
		{"username", *username},
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe la cuenta que queres seguir\n")
			return
		}
		log.Fatal(err)
	}
	var results model.Tweets
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	msg := ""
	for _, result := range results {
		cursor.Decode(&result)

		msg = result.Idtweet + result.Content + result.Timestamp.String()
		fmt.Fprintf(c, msg+"\n")
	}

}
func HandleMyFollowers(c net.Conn, arguments []string, username *string) {

	fmt.Println("Voy a handlear un myfollowers")
	if *username == "" {
		msg := "Tenes que estar logueado " // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("users")

	var user model.User //yo
	filter := bson.D{
		{"username", *username},
	}
	err = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe la cuenta \n")
			return
		}
		log.Fatal(err)
	}
	msg := "Las personas que te siguen son: " // mensaje de login exitoso
	fmt.Println(user.Followers)
	for i := 0; i < len(user.Followers); i++ {
		msg = msg + " " + user.Followers[i]
	}
	fmt.Fprintf(c, msg+"\n")

}
func HandleMyFollowing(c net.Conn, arguments []string, username *string) {
	fmt.Println("Voy a handlear un my following")
	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("users")

	var user model.User //yo
	filter := bson.D{
		{"username", *username},
	}
	err = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe la cuenta \n")
			return
		}
		log.Fatal(err)
	}
	msg := "Las personas que seguis son: " // mensaje de login exitoso
	fmt.Println(user.Following)
	for i := 0; i < len(user.Following); i++ {
		msg = msg + " " + user.Following[i]
	}
	fmt.Fprintf(c, msg+"\n")

}

func HandleFeed(c net.Conn, arguments []string, username *string) {

	fmt.Println("Voy a handlear un feed")

	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("users")

	var user model.User //yo
	filter := bson.D{
		{"username", *username},
	}
	err = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe la cuenta \n")
			return
		}
		log.Fatal(err)
	}
	following := user.Following

	coll = client.Database("tdl-los-tres-mosqueteros").Collection("tweets")
	days, err := strconv.Atoi(arguments[2])
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now().UTC()
	aux := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	filter = bson.D{
		{"username", bson.M{
			"$in": following,
		}},
		{"timestamp", bson.M{
			"$gte": aux.AddDate(0, 0, -days),
		}},
	}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe la cuenta que queres seguir\n")
			return
		}
		log.Fatal(err)
	}
	var results model.Tweets
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	msg := ""
	for _, result := range results {
		cursor.Decode(&result)

		msg = result.Idtweet + result.Content + result.Timestamp.String()
		fmt.Fprintf(c, msg+"\n")
	}

}
func HandleReply(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un reply")
}

func HandleAddTweetToThread(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un addtweedtothread")

}

func HandleNewThread(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un newthread")

}
func HandleThread(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un thread")

}

func HandleLike(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un like")

}
func HandleMostLiked(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un mostliked")

}

func HandleMostFollowed(c net.Conn, arguments []string) {

	fmt.Println("Voy a handlear un mostfollowed")

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
		HandleMyTweets(c, split_message, username)
	case MyFollowers:
		HandleMyFollowers(c, split_message, username)
	case MyFollowing:
		HandleMyFollowing(c, split_message, username)
	case Feed:
		HandleFeed(c, split_message, username)
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
		fmt.Fprintf(c, msg+"\n")
	}
	msg := "ok"
	fmt.Fprintf(c, msg+"\n")
}
