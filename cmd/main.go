package main

import (
	"log"
	"os"
	"os/signal"
	"github.com/egert811/task-server/internal/app"
)

func main() {
	//var wait time.Duration
	//flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	//flag.Parse()

	//Server init
	server, err := server.NewServer()

	if err != nil {
		log.Fatalf("Failed to initialize server: %s", err)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %s", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	server.Shutdown()

	log.Println("shutting down")
	os.Exit(0)
}
