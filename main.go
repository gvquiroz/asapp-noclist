package main

import (
	"fmt"
	"log"
)

func main() {

	BADSECService := NewService("http://localhost:8888")

	users, err := BADSECService.getUsers()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)

}
