package main

import (
	"log"
	"os"
	"os/signal"
	"github.com/egert811/task-server/internal/app/server"
	"github.com/egert811/task-server/internal/app/worker"
	"github.com/egert811/task-server/internal/pkg/storage"
)

func main() {
	commChan := make(chan storage.TaskDBItem)

	server, err := server.NewServer(commChan)
	worker := worker.NewWorker(commChan)

	if err != nil {
		log.Fatalf("Failed to initialize server: %s", err)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %s", err)
		}
	}()

	go worker.ExecuteAndPersist()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	server.Shutdown()

	log.Println("shutting down")
	os.Exit(0)
}
