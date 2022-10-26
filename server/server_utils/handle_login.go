package server_utils

import (
	"fmt"
)

// User login
func Handle_login() {
	fmt.Println("Enter Your First Name: ")

	// var then variable name then variable type
	var first string

	// Taking input from user
	fmt.Scanln(&first)
	fmt.Println("Enter Second Last Name: ")
	var second string
	fmt.Scanln(&second)

	fmt.Print("Your Full Name is: ")

	fmt.Print(first + " " + second)
}
