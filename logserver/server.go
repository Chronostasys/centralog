package logserver

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogListener interface {
	Listen(address string) error
}
type logServer struct {
	opts   *LogServerOptions
	client *mongo.Client
}

type LogServerOptions struct {
	MongoUrl           string
	Database           string
	Collection         string
	ExtraIndexes       []mongo.IndexModel // ts, level and id is automatically indexed
	ExpireAfterSeconds int32
}

func (ls *logServer) Listen(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Printf("received connection from %s", conn.RemoteAddr().String())
		go ls.handleConn(conn)
	}
}

func CreateLogListener(opts *LogServerOptions) (LogListener, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(opts.MongoUrl))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &logServer{
		opts:   opts,
		client: client,
	}, nil
}

func (ls *logServer) getIndexedCol(dbname, colname string) (*mongo.Collection, error) {
	db := ls.client.Database(dbname)
	col := db.Collection(colname)
	_, err := col.Indexes().CreateMany(
		context.Background(),
		append([]mongo.IndexModel{
			{
				Keys:    bson.M{"ts": -1},
				Options: &options.IndexOptions{ExpireAfterSeconds: &ls.opts.ExpireAfterSeconds},
			},
			{
				Keys: bson.M{"level": 1},
			},
			{
				Keys: bson.M{"id": 1},
			},
		}, ls.opts.ExtraIndexes...),
	)
	if err != nil {
		return nil, err
	}
	return col, nil
}

func (ls *logServer) handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	dec := json.NewDecoder(reader)
	logs := make([]interface{}, 10)
	i := 0
	errcount := 0
	first := true
	var col *mongo.Collection
	for {
		i = 0
		for {
			var doc map[string]interface{}
			err := dec.Decode(&doc)
			size := reader.Buffered()
			if err == io.EOF {
				log.Printf("connection closed from %s", conn.RemoteAddr().String())
				// all done
				return
			}
			if err != nil {
				log.Println(err)
				errcount++
				if strings.Contains(err.Error(), "closed") {
					// conn closed by remote, return
					return
				}
				if errcount > 10 {
					// something isn't right, reset connection
					return
				}
				continue
			}
			if first {
				msg := doc["msg"].(string)
				sps := strings.Split(msg, ";")
				var dbname, colname string
				if len(sps) == 2 && len(sps[0]) > 0 && len(sps[1]) > 0 {
					dbname, colname = sps[0], sps[1]
				} else {
					dbname, colname = ls.opts.Database, ls.opts.Collection
				}
				col, err = ls.getIndexedCol(dbname, colname)
				if err != nil {
					log.Println("create index failed", "dbname: "+dbname, "colname: "+colname, err)
					return
				}
				first = false
				log.Println("get indexed col success", conn.RemoteAddr(), "dbname: "+dbname, "colname: "+colname)
				if len(sps) == 2 && len(sps[0]) > 0 && len(sps[1]) > 0 {
					// skip proto log
					continue
				}
			}
			doc["ts"], _ = time.Parse("2006-01-02T15:04:05Z07:00", doc["ts"].(string))
			if i < len(logs) {
				logs[i] = doc
			} else {
				logs = append(logs, doc)
			}
			i++
			if size == 0 || i > 100 {
				break
			}
		}
		if i > 0 {
			_, err := col.InsertMany(context.Background(), logs[:i])
			if err != nil {
				log.Println(err)
			}
		}
	}
}
