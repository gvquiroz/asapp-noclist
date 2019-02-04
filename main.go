package main

import (
	"fmt"
	"log"
)

func main() {

	BADSECService := NewService("http://localhost:8888")

	numberOfRetryAttemps := 3

	users, err := BADSECService.getUsers(numberOfRetryAttemps)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)

}
