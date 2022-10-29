package server_utils

import (
	"fmt"
	"net"
)

// User login
func HandleLogin(c net.Conn) {
	fmt.Println("Voy a handlear un login")
	reply := make([]byte, 1024)

	_, _ = c.Read(reply)

	fmt.Println("Mensaje enviado por el cliente:", string(reply)) // aca deberia fijarse en bdd para chequear q este OK

	// si entro correctamente:

	msg := "ok" // mensaje de login exitoso
	_, _ = c.Write([]byte(msg))

}
