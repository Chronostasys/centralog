package main

import (
	"fmt"
	"os"
	"sync"

	log "github.com/Pivot-Studio/centralog/centralog"
	"go.uber.org/zap"
)

func main() {
	err := log.InitLogger(zap.NewProductionConfig(), "127.0.0.1:8001")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(no int) {
			writeLogs()
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func writeLogs() {
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(n int) {
			defer log.Sync() // flushes buffer, if any
			log.Info("test").Any("testdata", map[string]interface{}{
				"hello": "world",
				"age":   18,
			}).Log()
			log.Warn("test").Any("testdata", map[string]interface{}{
				"hello": "world",
				"age":   18,
			}).Log()
			log.Error("test").Any("testdata", map[string]interface{}{
				"hello": "world",
				"age":   18,
			}).Log()
			wg.Done()
		}(i)
	}
	wg.Wait()
}
