package worker

import (
	"bytes"
	"github.com/egert811/task-server/internal/pkg/storage"
	"log"
	"os/exec"
)

type Worker struct {
	store *storage.Store
	in    <-chan storage.TaskDBItem
}

func NewWorker(in <-chan storage.TaskDBItem) *Worker {
	return &Worker{
		store: storage.OpenStore(),
		in:    in,
	}
}

func (w *Worker) ExecuteAndPersist() {
	for {
		select {
		case ti := <-w.in:
			cmd := exec.Command(ti.CMD, ti.Args...)

			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Printf("Failed to execute cmd %s \n", err)
			}

			w.store.AddTaskOutput(&storage.TaskOutputDBItem{
				ID:     ti.ID,
				Output: out.String(),
			})
		}
	}
}
