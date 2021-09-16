package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Chronostasys/centralog/logserver"
)

func main() {
	var (
		db   string
		col  string
		port string
		exp  int
		conn string
	)
	flag.StringVar(&db, "db", "logdb", "database to store logs")
	flag.StringVar(&col, "col", "logcol", "collection to store logs")
	flag.StringVar(&port, "p", "8001", "port to listen")
	flag.StringVar(&conn, "c", "mongodb://localhost:27017", "mongodb connection string")
	flag.IntVar(&exp, "e", 360000, "how long a log will store")
	flag.Parse()
	server, err := logserver.CreateLogListener(&logserver.LogServerOptions{
		Database:           db,
		Collection:         col,
		ExpireAfterSeconds: int32(exp),
		MongoUrl:           conn,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server started listen on "+port)
	server.Listen("0.0.0.0:" + port)
}
