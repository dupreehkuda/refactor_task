package storage

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"sync"

	i "github.com/dupreehkuda/refactor_task/internal"
)

type storage struct {
	mtx   sync.Mutex
	file  string
	store UserStore
}

// New creates new instance of storage
func New(filepath string) *storage {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return &storage{
		file:  filepath,
		store: UserStore{},
	}
}

// readList reads JSON and returns UserStore instance
func (s *storage) readList() error {
	f, err := os.ReadFile(s.file)
	if err != nil {
		log.Printf("Error occurred reading file: %v", err)
		return err
	}

	err = json.Unmarshal(f, &s.store)
	if err != nil {
		log.Printf("Error occurred unmarshaling: %v", err)
		return err
	}

	return nil
}

// writeList rewrites JSON by passed UserStore
func (s *storage) writeList() error {
	body, err := json.Marshal(&s.store)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.file, body, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// checkUser throws an error if user already exist
func (s *storage) checkUser(name string, store UserStore) error {
	for _, val := range store.List {
		if val.DisplayName == name {
			log.Print("User already exist")
			return i.UserExists
		}
	}

	return nil
}
