package dbconnect

import (
	"encoding/binary"
	"errors"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

//Task struct has same structure as tasks Bucket.
type Task struct {
	Key   int
	Value string
}

//Init take a dbPath as argument and establish the connection to it
//If connection successful then creates a bucket of value taskBucket.
//If any error occur during connection or bucket creation it will return the error.
func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println("DB connection error: ", err)
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		if err != nil {
			log.Printf("Failed to create Bucket: %s, error: %v", taskBucket, err)
		}
		return err
	})
}

//CloseDB only close the connection to boltDB.
func CloseDB() {
	db.Close()
}

//CreateTask take a task as parameter
//and add that task to the Bucket.
//If successfully added then it will return the id of that task
//otherwise return -1 with the error
func CreateTask(task string) (int, error) {
	var taskId int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		if b == nil {
			return errors.New("Bucket not found")
		}
		taskId64, _ := b.NextSequence()
		taskId = int(taskId64)
		key := itob(taskId)
		return b.Put(key, []byte(task))
	})
	if err != nil {
		log.Println("Read-write transaction Failed: ", err)
		return -1, err
	}
	return taskId, nil
}

//AllTasks return the list of tasks present in the Bucket
//return type is slice of Task struct
//If any error occurs it will return nil with the error.
func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		if b == nil {
			return errors.New("Bucket not found")
		}
		b.ForEach(func(k, v []byte) error {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
			return nil
		})
		return nil
	})
	if err != nil {
		log.Println("Read-only transaction Failed: ", err)
		return nil, err
	}
	return tasks, nil

}

//DeleteTask takes a key as argument and
//delete the same task from the Bucket.
func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		if b == nil {
			return errors.New("Bucket not found")
		}
		return b.Delete(itob(key))
	})
}

//itob convert the integer to byte slice.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

//btoi convert the byte slice to integer.
func btoi(b []byte) int {
	i := int(binary.BigEndian.Uint64(b))
	return i
}
