package main

import (
	"log"

	"github.com/Pivot-Studio/centralog/logserver"
)

func main() {
	ls, err := logserver.CreateLogListener(&logserver.LogServerOptions{
		MongoUrl:           "mongodb://localhost:27018",
		Database:           "testlog",
		Collection:         "zaplog",
		ExpireAfterSeconds: 3600,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = ls.Listen("0.0.0.0:8001")
	if err != nil {
		log.Fatal(err)
	}
}
