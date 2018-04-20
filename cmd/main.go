package main

import (
	"context"
	"fmt"
	"github.com/coreos/bbolt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"encoding/json"
)

type Task struct {
	ID  int `json:"id"`
	CMD string `json:"cmd"`
	Output string `json:"output"`
}

type Store struct {
	db bolt.DB
}

func openStore() (*Store, error) {

	return nil, nil
}

func (s *Store) addTask(t *Task) error {
	return nil
}

func (s *Store) getTasks() ([]Task, error) {
	return nil, nil
}

func (s *Store) getTaskById(id int) (*Task, error) {
	return nil, nil
}

// http handlers
func handleTaskPost(w http.ResponseWriter, r *http.Request) {
	t := Task{
		ID: 1,
		CMD: "ls -alh",
	}

	json.NewEncoder(w).Encode(t)
}

func handleTaskGet(w http.ResponseWriter, r *http.Request) {
	t := Task{
		ID: 1,
		CMD: "ls -alh",
	}

	json.NewEncoder(w).Encode(t)
}

func handleTaskDetailsGet(w http.ResponseWriter, r *http.Request) {
	t := Task{
		ID: 1,
		CMD: "ls -alh",
	}

	json.NewEncoder(w).Encode(t)
}

// server config, TODO: externalize
var serverPort 				int = 3000
var serverWriteTimeout  	time.Duration = 15
var serverReadTimeout 		time.Duration = 15
var serverIdleTimeout 	 	time.Duration = 60
var serverShutdownTimeout 	time.Duration = 15

func main() {
	//var wait time.Duration
	//flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	//flag.Parse()

	r := mux.NewRouter()
	// Add your routes as needed
	r.HandleFunc("/task", handleTaskPost).Methods("POST")
	r.HandleFunc("/task", handleTaskGet).Methods("GET")
	r.HandleFunc("/task/{id}", handleTaskDetailsGet).Methods("GET")

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", serverPort),
		WriteTimeout: time.Second * serverWriteTimeout,
		ReadTimeout:  time.Second * serverReadTimeout,
		IdleTimeout:  time.Second * serverIdleTimeout,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}
