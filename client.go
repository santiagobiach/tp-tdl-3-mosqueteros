package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	//Se provee host:port como argumento del programa
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}
	//Se intenta conectar
	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Lee una linea y espera la devolucion del server
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}

