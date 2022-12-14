package server_utils

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"server/database"
	"server/model"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	NewThread          string = "newThread"
	AddTweetToThread   string = "addTweetToThread"
	Thread             string = "thread"
	Like               string = "like"
	MostLiked          string = "mostLiked" //puede ser de 2 tipos
	MostFollowed       string = "mostFollowed"
	Help string = "help"
)

func HandleHelp(c net.Conn) {
	fmt.Fprintf(c, "BIENVENIDO A TWITTER\n")

	fmt.Fprintf(c, "COMANDOS: \n")

	fmt.Fprintf(c, "	creá tu cuenta: \n")
	fmt.Fprintf(c, "		signup <usuario> <contraseña> <contraseña> \n")

	fmt.Fprintf(c, "	ingresá en tu cuenta: \n")
	fmt.Fprintf(c, "		login <usuario> <contraseña> \n")

	fmt.Fprintf(c, "	tweetea algo: \n")
	fmt.Fprintf(c, "		tweet <escribir algo> \n")

	fmt.Fprintf(c, "	seguir a otro usuario: \n")
	fmt.Fprintf(c, "		follow <usuario> \n")

	fmt.Fprintf(c, "	dejá de seguir a otro usuario: \n")
	fmt.Fprintf(c, "		unfollow <usuario> \n")

	fmt.Fprintf(c, "	ver los tweets de otro usuario en los últimos días: \n")
	fmt.Fprintf(c, "		tweetsFrom <usuario> <cantidad de días> \n")

	fmt.Fprintf(c, "	ver cierta cantidad de top tendencias en los últimos días: \n")
	fmt.Fprintf(c, "		trendingTopic <cantidad de tendencias> <cantidad de días> \n")

	fmt.Fprintf(c, "	ver cierta cantidad de tweets de cierta tendencia: \n")
	fmt.Fprintf(c, "		trendingTweetsFrom <tendencia> <cantidad de tweets> \n")

	fmt.Fprintf(c, "	ver mis tweets: \n")
	fmt.Fprintf(c, "		myTweets \n")

	fmt.Fprintf(c, "	ver los usuarios que me siguen: \n")
	fmt.Fprintf(c, "		myFollowers \n")

	fmt.Fprintf(c, "	ver tweets de los usuarios que sigo en los últimos días: \n")
	fmt.Fprintf(c, "		feed <cantidad de días> \n")

	fmt.Fprintf(c, "	crear un hilo: \n")
	fmt.Fprintf(c, "		newThread <nombre del thread> <primer tweet> \n")

	fmt.Fprintf(c, "	agregar un tweet a un thread ya existente: \n")
	fmt.Fprintf(c, "		addTweetToThread <nombre del thread> <tweet> \n")

	fmt.Fprintf(c, "	ver todos los tweets de alguno de mis threads: \n")
	fmt.Fprintf(c, "		thread <nombre del thread> \n")

	fmt.Fprintf(c, "	ver los usuarios a los que sigo: \n")
	fmt.Fprintf(c, "		myFollowing \n")

	fmt.Fprintf(c, "	ver el usuario más seguido de twitter: \n")
	fmt.Fprintf(c, "		mostFollowed \n")
}
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
	msg := "¡Bienvenido, " + user.Username + "!"// mensaje de login exitoso
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

			msg := "¡Tu cuenta ha sido creada!" // mensaje de login exitoso
			fmt.Fprintf(c, msg+"\n")
			return
		}
		log.Fatal(err)
	}
	fmt.Fprintf(c, "Ya existe la cuenta \n")

}

func HandleTweet(c net.Conn, arguments []string, username *string) {

	if *username == "" {
		msg := "Tenes que estar logueado para ejecutar este comando" // mensaje de error
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

	coll_tweet := client.Database("tdl-los-tres-mosqueteros").Collection("tweets")
	topics := client.Database("tdl-los-tres-mosqueteros").Collection("topics")
	var tweet model.Tweet
	var tweet_content string
	topics_content := make([]string, 5)
	for i := 1; i < len(arguments); i++ {
		//Si la palabra empieza con # agregar a una coleccion de la bdd de topics
		if strings.HasPrefix(arguments[i], "#") {
			topics_content = append(topics_content, arguments[i])
		}
		tweet_content += arguments[i] + " "
	}

	tweet.Content = tweet_content
	tweet.Username = *username
	tweet.Timestamp = time.Now()
	insertTweet, err := coll_tweet.InsertOne(ctx, tweet)
	for t := range topics_content {
		opts := options.Update().SetUpsert(true)
		filter := bson.D{
			{"Topicstring", t},
		}
		update := bson.D{{"$push", bson.D{{"Tweets", tweet}}}}
		_, err := topics.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("hay un nuevo tweet, on id : ", insertTweet.InsertedID)
	fmt.Println("El nuevo tweet es de ", tweet.Username)

	msg := "Tweet enviado satisfactoriamente" // mensaje de tweet exitoso
	fmt.Fprintf(c, msg+"\n")

}
func HandleFollow(c net.Conn, arguments []string, username *string) {
	fmt.Println("Voy a handlear un follow")
	if *username == "" {
		msg := "Tenes que estar logueado para ejecutar este comando" // mensaje de error
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
		msg := "Tenes que estar logueado para ejecutar este comando" // mensaje de error
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
		msg := "Tenes que estar logueado para ejecutar este comando" // mensaje de error
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

		msg = result.Idtweet + result.Content + " - " + result.Timestamp.Format(time.RFC822)
		fmt.Fprintf(c, msg+"\n")
	}

}
func HandleTrendingTopic(c net.Conn, arguments []string, username *string) {
	fmt.Println("Voy a handlear un tt")
	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}
	if len(arguments) != 2 {
		// arguments[1] = nombre de la tendencia
		// arguments[2] = numero de tweets para ver de la tendencia
		msg := "Comando incompleto"
		fmt.Fprintf(c, msg+"\n")
		return
	}
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)
	// por ahora lo hago sin los d dias y los n tweets, busco una tendencia y devuelvo sus tweets
	coll := client.Database("tdl-los-tres-mosqueteros").Collection("topics") //cambiar a tt

	//var limit int
	var requested_n, errr = strconv.Atoi(arguments[1])
	if errr != nil {
		log.Fatal(err)
	}

	filter := bson.D{
		{},
	}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe la cuenta que queres seguir\n")
			return
		}
		log.Fatal(err)
	}
	var results model.Topics
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	println(results)
	sort.Slice(results[:], func(i, j int) bool {
		return len(results[i].Tweets) < len(results[i].Tweets)
	})

	var limit int
	if requested_n >= len(results) {
		limit = len(results)
	} else {
		limit = requested_n
	}

	for i := 0; i < limit; i++ {
		// busco el tweet y lo muestro.
		fmt.Fprintf(c, strconv.Itoa(i)+results[i].Topicstring+"\n")
	}

}

func HandleTrendingTweetsFrom(c net.Conn, arguments []string, username *string) {
	//Tweets from topic
	fmt.Println("Voy a handlear un ttfrom")
	if *username == "" {
		msg := "Tenes que estar logueado" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}
	if len(arguments) != 3 {
		// arguments[1] = nombre de la tendencia
		// arguments[2] = numero de tweets para ver de la tendencia
		msg := "Comando incompleto"
		fmt.Fprintf(c, msg+"\n")
		return
	}
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)
	// por ahora lo hago sin los d dias y los n tweets, busco una tendencia y devuelvo sus tweets
	coll := client.Database("tdl-los-tres-mosqueteros").Collection("topics")

	//var limit int
	var requested_n, errr = strconv.Atoi(arguments[2])
	if errr != nil {
		log.Fatal(err)
	}

	var topic model.Topic
	filter := bson.D{
		{"Topicstring", arguments[1]},
	}
	err = coll.FindOne(ctx, filter).Decode(&topic)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe el topic \n")
			return
		}
		log.Fatal(err)
	}

	var limit int
	if requested_n >= len(topic.Tweets) {
		limit = len(topic.Tweets)
	} else {
		limit = requested_n
	}
	tope := len(topic.Tweets) - 1
	for i := 0; i < limit; i++ {
		// busco el tweet y lo muestro.
		tweet := topic.Tweets[tope-i]
		msg := ""

		msg = tweet.Username + ":" + tweet.Idtweet + tweet.Content + tweet.Timestamp.Format(time.RFC822)
		fmt.Fprintf(c, msg+"\n")
	}

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
		msg = result.Idtweet + result.Content + " - " + result.Timestamp.Format(time.RFC822)
		fmt.Fprintf(c, msg+"\n")
	}

}
func HandleMyFollowers(c net.Conn, arguments []string, username *string) {

	fmt.Println("Voy a handlear un myfollowers")
	if *username == "" {
		msg := "Tenes que estar logueado" 
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
	days, err := strconv.Atoi(arguments[1])
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

		msg = result.Username + result.Idtweet + result.Content + " - " + result.Timestamp.Format(time.RFC822)
		fmt.Fprintf(c, msg+"\n")
	}

}

func HandleAddTweetToThread(c net.Conn, arguments []string, username *string) {

	fmt.Println("Voy a handlear un addtweedtothread")
	if *username == "" {
		msg := "Tenes que estar logueado para ejecutar este comando" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	// creo el tweet
	coll_tweet := client.Database("tdl-los-tres-mosqueteros").Collection("tweets")

	var tweet model.Tweet
	var tweet_content string

	for i := 2; i < len(arguments); i++ {
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
	// creo el thread con el ID del primer tweet en su lista de tweets

	// agrego el tweet al thread
	coll := client.Database("tdl-los-tres-mosqueteros").Collection("threads")

	var thread model.Thread // yo
	filter := bson.D{
		{"threadname", arguments[1]},
	}
	err = coll.FindOne(ctx, filter).Decode(&thread)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe el thread \n")
			return
		}
		log.Fatal(err)
	}

	thread.Tweets = append(thread.Tweets, insertTweet.InsertedID.(primitive.ObjectID).Hex())
	_, err = coll.ReplaceOne(ctx, filter, thread)
	msg := "¡Has agregado un tweet a tu thread!" // mensaje de login exitoso
	fmt.Fprintf(c, msg+"\n")
}

func HandleNewThread(c net.Conn, arguments []string, username *string) {

	// newThread <name> <firstTweet>
	fmt.Println("Voy a handlear un newThread")
	if *username == "" {
		msg := "Tenes que estar logueado para ejecutar este comando" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}

	// creo el tweet y lo guardo en la BDD, me quedo con su ID
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll_tweet := client.Database("tdl-los-tres-mosqueteros").Collection("tweets")

	var tweet model.Tweet
	var tweet_content string

	for i := 2; i < len(arguments); i++ {
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
	// var tweetId string = insertTweet.InsertedID

	// creo el thread con el ID del primer tweet en su lista de tweets

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("threads")

	var thread model.Thread
	thread.Threadname = arguments[1]
	thread.Tweets = []string{}
	thread.Tweets = append(thread.Tweets, insertTweet.InsertedID.(primitive.ObjectID).Hex())
	fmt.Println(thread.Tweets)
	_, err = coll.InsertOne(ctx, thread)
	if err != nil {
		log.Fatal(err)
	}

	// le agrego el nombre del thread al usuario en su lista de threads
	coll = client.Database("tdl-los-tres-mosqueteros").Collection("users")

	var user model.User //yo
	filter := bson.D{
		{"username", *username},
	}

	_ = coll.FindOne(ctx, filter).Decode(&user) // no deberia tirar error xq estoy logueado ok
	user.Threads = append(user.Threads, thread.Threadname)
	_, err = coll.ReplaceOne(ctx, filter, user)
	if err != nil {
		log.Fatal(err)
	}

	// mensaje de éxito

	msg := "¡Creaste un thread correctamente!"
	fmt.Fprintf(c, msg+"\n")
}

func HandleThread(c net.Conn, arguments []string, username *string) {

	fmt.Println("Voy a handlear un thread")

	if *username == "" {
		msg := "Tenes que estar logueado para ejecutar este comando" // mensaje de error
		fmt.Fprintf(c, msg+"\n")
		return
	}
	// solo busco mi thread y lo veo
	// como los threads se guardan en una base de datos todos juntos,
	// no puden tener nombre repetido. entonces no hace falta que busque al usuario y de ahí al thread
	//, puedo buscar directamente el thread

	// dejo todo comentado para no olvidarme jaja

	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("threads")

	var thread model.Thread // yo
	filter := bson.D{
		{"threadname", arguments[1]},
	}
	err = coll.FindOne(ctx, filter).Decode(&thread)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(c, "No existe el thread \n")
			return
		}
		log.Fatal(err)
	}
	coll_tweets := client.Database("tdl-los-tres-mosqueteros").Collection("tweets")

	for i := 0; i < len(thread.Tweets); i++ {
		// busco el tweet y lo muestro.
		var tweet model.Tweet
		id, error := primitive.ObjectIDFromHex(thread.Tweets[i])
		if error != nil {
			fmt.Fprintf(c, "Error")
		}
		filter = bson.D{
			{"_id", id}, // problema: no lo está encontrando, lo estoy buscando mal
		}
		_ = coll_tweets.FindOne(ctx, filter).Decode(&tweet)
		fmt.Fprintf(c, tweet.Content+"\n")
	}

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
		HandleTrendingTopic(c, split_message, username)
	case TrendingTweetsFrom:
		HandleTrendingTweetsFrom(c, split_message, username)
	case MyTweets:
		HandleMyTweets(c, split_message, username)
	case MyFollowers:
		HandleMyFollowers(c, split_message, username)
	case MyFollowing:
		HandleMyFollowing(c, split_message, username)
	case Feed:
		HandleFeed(c, split_message, username)
	// case Reply:
	// 	HandleReply(c, split_message)
	case NewThread:
		HandleNewThread(c, split_message, username)
	case AddTweetToThread:
		HandleAddTweetToThread(c, split_message, username)
	case Thread:
		HandleThread(c, split_message, username)
	case Like:
		HandleLike(c, split_message)
	case MostLiked:
		HandleMostLiked(c, split_message)
	case MostFollowed:
		HandleMostFollowed(c, split_message)
	case Help:
		HandleHelp(c)
	default:
		msg := "El comando que has ingresado no existe" // comando no existe
		fmt.Fprintf(c, msg+"\n")
	}
	msg := "ok"
	fmt.Fprintf(c, msg+"\n")
}
func UpdateTrendingTopics() {
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	coll := client.Database("tdl-los-tres-mosqueteros").Collection("tweets")
	tts := client.Database("tdl-los-tres-mosqueteros").Collection("trendingtopics")

	tts.DeleteMany(ctx, bson.M{})
	//Tambien borrar lo de trending topics
	now := time.Now().UTC()
	aux := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	filter := bson.D{
		{"timestamp", bson.M{
			"$gte": aux.Add(-5 * time.Minute),
		}},
	}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	var results model.Tweets
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	cant_goroutines := 4
	var channels [4]chan model.Tweet
	for i := 0; i < cant_goroutines; i++ {
		channels[i] = make(chan model.Tweet)
		go ProcessTweets(channels[i])
	}
	for i, result := range results {
		cursor.Decode(&result)
		aux := i % cant_goroutines
		channels[aux] <- result
	}
	for i := range channels {
		close(channels[i])
	}
	println("TERMINANDO HILO PRINCIPAL")
}
func ProcessTweets(c chan model.Tweet) {
	client, ctx, cancel, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(client, ctx, cancel)
	database.Ping(client, ctx)

	tts := client.Database("tdl-los-tres-mosqueteros").Collection("trendingtopics")

	for {
		tweet, ok := <-c
		if !ok {
			println("TERMINANDO HILO SECUNDARIO")
			return
		}
		split_message := strings.SplitAfter(tweet.Content, " ")
		for i, v := range split_message {
			split_message[i] = strings.TrimSpace(v)
			if strings.HasPrefix(split_message[i], "#") {
				opts := options.Update().SetUpsert(true)
				filter := bson.D{
					{"Topicstring", split_message[i]},
				}
				update := bson.D{{"$push", bson.D{{"Tweets", tweet}}}}
				_, err := tts.UpdateOne(ctx, filter, update, opts)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}
