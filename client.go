package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func checkErr(err error) {

	if err != nil {
		log.Fatal(err)
	}
}

func send_login(c net.Conn) {
	fmt.Println("Username:")
	var username string
	fmt.Scanln(&username)
	fmt.Println("Password:")
	var password string
	fmt.Scanln(&password)
	msg := "li" + " " + username + " " + password // creo mensaje de login
	_, err := c.Write([]byte(msg))
	checkErr(err)
	reply := make([]byte, 2)

	_, err = c.Read(reply)
	checkErr(err)
	fmt.Println("Supuesto mensaje de ok:", string(reply))
	if string(reply) == "ok" {
		fmt.Println("Ingresaste correctamente!! Ahora podes empezar a usar twitter")
	}
	// En caso de error puede llamar de nuevo a send_login() para que intente de nuevo
}
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
	checkErr(err)
	send_login(c)
	//Lee una linea y espera la devolucion del server
	//Empieza proceso de ingreso de comandos, (el login es algo aparte)
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

// MENSAJES:

// client->server (mensaje de login): "li username password"
// server->client (respuesta de login): "ok" en caso de éxito
