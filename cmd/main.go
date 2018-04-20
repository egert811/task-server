package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/coreos/bbolt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Task struct {
	*TaskDBItem
	*TaskOutputDBItem
}

type TaskDBItem struct {
	ID  int    `json:"id"`
	CMD string `json:"cmd"`
}

type TaskOutputDBItem struct {
	ID     int    `json:"-"`
	Output string `json:"output"`
}

var (
	dbPath             string = "task.db"
	dbTaskBucket       []byte = []byte("Tasks")
	dbTaskOutputBucket []byte = []byte("TaskOutputs")
)

// Bolt helpers
// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

type Store struct {
	db *bolt.DB
}

func openStore() (*Store, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return nil, err
	}

	// Ensure that the buckets are exists
	err = db.Update(func(tx *bolt.Tx) error {

		tx.CreateBucketIfNotExists(dbTaskBucket)
		tx.CreateBucketIfNotExists(dbTaskOutputBucket)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

// Task bucket handlers
func (s *Store) addTask(t *TaskDBItem) error {
	return s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(dbTaskBucket)

		//grap the next id
		id, _ := b.NextSequence()
		t.ID = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(t.ID), buf)
	})

}

func (s *Store) getTasks() ([]TaskDBItem, error) {
	resp := make([]TaskDBItem, 0)

	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(dbTaskBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var t TaskDBItem
			err := json.Unmarshal(v, &t)
			if err != nil {
				return err
			}

			resp = append(resp, t)

		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil

}

// Task output handlers
func (s *Store) addTaskOutput(t *TaskDBItem) error {
	return s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(dbTaskOutputBucket)

		//grap the next id
		id, _ := b.NextSequence()
		t.ID = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(t.ID), buf)
	})

}

// fetch all
func (s *Store) getTaskDetailsById(id int) (*Task, error) {
	var resp Task

	err := s.db.View(func(tx *bolt.Tx) error {

		taskItemBytes := tx.Bucket(dbTaskBucket).Get(itob(id))

		var t TaskDBItem
		err := json.Unmarshal(taskItemBytes, &t)

		if err != nil {
			return err
		}

		taskOutputBytes := tx.Bucket(dbTaskOutputBucket).Get(itob(id))
		var to TaskOutputDBItem

		if len(taskOutputBytes) > 0 {
			err = json.Unmarshal(taskOutputBytes, &to)

			if err != nil {
				return err
			}
		} else {
			to = TaskOutputDBItem{
				ID:     id,
				Output: "",
			}
		}

		resp = Task{&t, &to}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

//app
type Server struct {
	store *Store
}

// http handlers
func (s *Server) handleTaskPost(w http.ResponseWriter, r *http.Request) {
	var t TaskDBItem
	json.NewDecoder(r.Body).Decode(&t)

	err := s.store.addTask(&t)

	if err != nil {
		//render error an response here
	}

	json.NewEncoder(w).Encode(t)
}

func (s *Server) handleTaskGet(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.store.getTasks()

	if err != nil {
		//render error an response here
	}

	json.NewEncoder(w).Encode(tasks)
}

func (s *Server) handleTaskDetailsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		//render error an response here
	}

	var t *Task
	t, err = s.store.getTaskDetailsById(id)

	if err != nil {
		//render error an response here
	}

	json.NewEncoder(w).Encode(t)
}

// server config, TODO: externalize
var serverPort int = 3000
var serverWriteTimeout time.Duration = 15
var serverReadTimeout time.Duration = 15
var serverIdleTimeout time.Duration = 60
var serverShutdownTimeout time.Duration = 15

func main() {
	//var wait time.Duration
	//flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	//flag.Parse()

	//Server init
	store, err := openStore()

	if err != nil {
		log.Fatalf("Failed to initialize the store: %s", err)
	}

	s := Server{
		store: store,
	}

	r := mux.NewRouter()
	// Add your routes as needed
	r.HandleFunc("/task", s.handleTaskPost).Methods("POST")
	r.HandleFunc("/task", s.handleTaskGet).Methods("GET")
	r.HandleFunc("/task/{id}", s.handleTaskDetailsGet).Methods("GET")

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
