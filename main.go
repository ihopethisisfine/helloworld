package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ihopethisisfine/helloworld/internal/pkg/storage/aws"
	"github.com/ihopethisisfine/helloworld/internal/user"
)

func main() {
	// Create a session instance.
	ses, err := aws.New(aws.Config{
		Address: "http://localhost:4566",
		Region:  "eu-west-1",
		Profile: "localstack",
		ID:      "test",
		Secret:  "test",
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Instantiate HTTP app
	usr := user.Controller{
		Storage: aws.NewUserStorage(ses, time.Second*5),
	}

	// Instantiate HTTP router
	rtr := http.NewServeMux()
	rtr.HandleFunc("/hello/", usr.Hello)

	// Start HTTP server
	addr := ":8080"
	log.Println("listen on", addr)
	log.Fatalln(http.ListenAndServe(addr, rtr))
}
