package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Chronostasys/centralog/centralog"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	var (
		conn string
	)
	flag.StringVar(&conn, "s", "127.0.0.1:8001", "server address")
	flag.Parse()
	err := centralog.InitLoggerWithOpt(zap.NewProductionConfig(), &centralog.LogOptions{
		Server:     conn,
		Db:         "logtest",
		Collection: "logtest",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	writer := bufio.NewWriter(os.Stdout)
	id := uuid.New()
	ctx := context.Background()
	ctx = context.WithValue(ctx, centralog.IDKey, id)
	_, err = writer.WriteString("using ctx id: " + id.String() + "\n")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		_, err := writer.WriteString("centralog> ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Read the keyboad input.
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		input = strings.Trim(input, "\n\r ")
		if len(input) == 0 {
			continue
		}
		centralog.Info(input).CtxID(ctx).Log()
		centralog.Sync()
	}
}
