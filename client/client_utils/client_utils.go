package client_utils

import (
	"fmt"
	"log"
	"net"
)

// Error checker
func CheckError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}

// Send login message to server
func SendLogin(c net.Conn) {
	fmt.Print("Username: ")
	var username string
	fmt.Scanln(&username)
	fmt.Print("Password: ")
	var password string
	fmt.Scanln(&password)
	msg := "login" + " " + username + " " + password + "\n" // creo mensaje de login
	_, err := c.Write([]byte(msg))
	CheckError(err)
	reply := make([]byte, 2)

	_, err = c.Read(reply)
	CheckError(err)
	if string(reply) == "ok" {
		fmt.Println("Ingresaste correctamente!! Ahora podes empezar a usar twitter")
	}
	// En caso de error puede llamar de nuevo a send_login() para que intente de nuevo
}
