package main

import (
	"fmt"
	"log"
)

func main() {

	BADSECService := NewService("http://localhost:8888")

	numberOfRetryAttemps := 2

	users, err := BADSECService.getUsers(numberOfRetryAttemps)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)

}
