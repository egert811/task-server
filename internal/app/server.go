package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/egert811/task-server/internal/pkg"
	"time"
	"fmt"
	"context"
)

// server config, TODO: externalize
var (
	serverPort            int           = 3000
	serverWriteTimeout    time.Duration = 15
	serverReadTimeout     time.Duration = 15
	serverIdleTimeout     time.Duration = 60
	serverShutdownTimeout time.Duration = 15
)

type Task struct {
	*storage.TaskDBItem
	*storage.TaskOutputDBItem
}

type Server struct {
	store *storage.Store
	router *mux.Router
	server *http.Server
}

func NewServer() (*Server, error) {
	//Server init
	store, err := storage.OpenStore()

	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", serverPort),
		WriteTimeout: time.Second * serverWriteTimeout,
		ReadTimeout:  time.Second * serverReadTimeout,
		IdleTimeout:  time.Second * serverIdleTimeout,
		Handler:      r,
	}


	return &Server{
		store: store,
		router: r,
		server: srv,
	}, nil
}

func (s *Server) ListenAndServe() error {

	s.router.HandleFunc("/task", s.handleTaskPost).Methods("POST")
	s.router.HandleFunc("/task", s.handleTaskGet).Methods("GET")
	s.router.HandleFunc("/task/{id}", s.handleTaskDetailsGet).Methods("GET")

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}


// http handlers
func (s *Server) handleTaskPost(w http.ResponseWriter, r *http.Request) {
	var t storage.TaskDBItem
	json.NewDecoder(r.Body).Decode(&t)

	err := s.store.AddTask(&t)

	if err != nil {
		//render error an response here
	}

	json.NewEncoder(w).Encode(t)
}

func (s *Server) handleTaskGet(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.store.GetTasks()

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

	var taskItem *storage.TaskDBItem
	var taskOutput *storage.TaskOutputDBItem
	taskItem, taskOutput, err = s.store.GetTaskDetailsById(id)

	if err != nil {
		//render error an response here
	}

	resp := Task{taskItem, taskOutput}

	json.NewEncoder(w).Encode(resp)
}
