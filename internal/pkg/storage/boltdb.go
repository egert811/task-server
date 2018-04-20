package storage

import (
	"encoding/binary"
	"encoding/json"
	"github.com/coreos/bbolt"
	"log"
	"time"
)

// TODO: externalize
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

// don't do this in production!
var singletonStore *Store

func init() {
	var err error
	singletonStore, err = createStore()

	if err != nil {
		log.Fatal("Failed to initilize BoltDB")
	}

}

func createStore() (*Store, error) {
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

func OpenStore() *Store {
	return singletonStore
}

// Task bucket handlers
func (s *Store) AddTask(t *TaskDBItem) error {
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

func (s *Store) GetTasks() ([]TaskDBItem, error) {
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
func (s *Store) AddTaskOutput(t *TaskOutputDBItem) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbTaskOutputBucket)

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
func (s *Store) GetTaskDetailsById(id int) (ti *TaskDBItem, toi *TaskOutputDBItem, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {

		taskItemBytes := tx.Bucket(dbTaskBucket).Get(itob(id))

		err := json.Unmarshal(taskItemBytes, &ti)

		if err != nil {
			return err
		}

		taskOutputBytes := tx.Bucket(dbTaskOutputBucket).Get(itob(id))

		if len(taskOutputBytes) > 0 {
			err = json.Unmarshal(taskOutputBytes, &toi)

			if err != nil {
				return err
			}
		} else {
			toi = &TaskOutputDBItem{
				ID:     id,
				Output: "",
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return
}
