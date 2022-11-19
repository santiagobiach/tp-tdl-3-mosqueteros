package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"client/client_utils"
	"regexp"
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
	client_utils.CheckError(err)
	//client_utils.SendLogin(c)
	//Lee una linea y espera la devolucion del server
	//Empieza proceso de ingreso de comandos, (el login es algo aparte)
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		//Esto definitivamente no va.
		fmt.Println("->: " + regexp.MustCompile(`[^a-zA-Z ]+`).ReplaceAllString(message, ""))
		fmt.Println()
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}

// MENSAJES:

// client->server (mensaje de login): "li username password"
// server->client (respuesta de login): "ok" en caso de éxito
