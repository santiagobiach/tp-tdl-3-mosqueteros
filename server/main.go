package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"server/server_utils"
	"strings"
	"time"
)

func handleConnection(c net.Conn) {
	fmt.Println("New client connected")
	// server_utils.HandleLogin(c)

	// Nuevos mensajes. Hay que determinar cuál es el pedido que llega y procesarlo
	var username string
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
		fmt.Println(temp)
		server_utils.ParseMessage(c, temp, &username)
	}
	c.Close()
}

func createTrendingTopics() {
	for {
		server_utils.UpdateTrendingTopics()
		time.Sleep(10 * time.Minute)
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	go createTrendingTopics()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
