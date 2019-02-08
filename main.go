package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	badsecEndpoint := os.Getenv("BADSEC_ENDPOINT")

	// If env badsec endpoint is not defined, init as default value localhost
	if badsecEndpoint == "" {
		badsecEndpoint = "http://localhost:8888"
	}

	BADSECService := NewService(badsecEndpoint)

	numberOfRetryAttemps := 2
	users, err := BADSECService.getUsers(numberOfRetryAttemps)

	// Exits the application with non-zero status
	if err != nil {
		log.Fatal(err)
	}

	// Prints on Stdout
	fmt.Println(users)

}
