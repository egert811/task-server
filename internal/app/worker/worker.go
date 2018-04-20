package worker

import (
	"github.com/egert811/task-server/internal/pkg/storage"
	"os/exec"
	"bytes"
	"log"
)

type Worker struct {
	store *storage.Store
	in <- chan storage.TaskDBItem
}


func NewWorker(in <- chan storage.TaskDBItem) (*Worker) {
	return &Worker{
		store: storage.OpenStore(),
		in : in,
	}
}

func (w* Worker) ExecuteAndPersist( in <- chan storage.TaskDBItem ){
	for {
		select {
		case ti := <- w.in:
			_ := ti

			cmd := exec.Command("ls", "-alh")

			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}

			w.store.AddTaskOutput(&storage.TaskOutputDBItem{
				ID: ti.ID,
				Output: out.String(),
			})
		}
	}
}