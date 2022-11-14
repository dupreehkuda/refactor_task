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
	mtx  sync.Mutex
	file string
}

// New creates new instance of storage
func New(filepath string) *storage {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return &storage{
		file: filepath,
	}
}

// readList reads JSON and returns UserStore instance
func (s *storage) readList() (UserStore, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	f, err := os.ReadFile(s.file)
	if err != nil {
		log.Printf("Error occurred reading file: %v", err)
		return UserStore{}, err
	}

	store := UserStore{}

	err = json.Unmarshal(f, &store)
	if err != nil {
		log.Printf("Error occurred unmarshaling: %v", err)
		return UserStore{}, err
	}

	return store, nil
}

// writeList rewrites JSON by passed UserStore
func (s *storage) writeList(store *UserStore) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	body, err := json.Marshal(&store)
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
