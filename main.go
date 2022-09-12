package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ihopethisisfine/helloworld/internal/pkg/storage/aws"
	"github.com/ihopethisisfine/helloworld/internal/user"
)

func main() {
	// Create a session instance.
	var config *aws.Config
	if len(os.Getenv("DYNAMODB_ENDPOINT")) > 0 {
		config = &aws.Config{
			Address: os.Getenv("DYNAMODB_ENDPOINT"),
			Region:  "eu-west-1",
			Profile: "localdynamo",
			ID:      "test",
			Secret:  "test",
		}
	}
	ses, err := aws.New(config)

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
