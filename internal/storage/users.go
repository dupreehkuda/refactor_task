package storage

import (
	"encoding/json"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	i "github.com/dupreehkuda/refactor_task/internal"
)

// SearchUsers returns list of all users
func (s *storage) SearchUsers() ([]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	err := s.readList()
	if err != nil {
		log.Printf("Error occurred getting file: %v", err)
		return nil, err
	}

	resp, err := json.Marshal(s.store.List)
	if err != nil {
		log.Printf("Error occurred marshaling: %v", err)
		return nil, err
	}

	return resp, nil
}

// CreateUser creates new user if name doesn't exist
func (s *storage) CreateUser(name, email string) (string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	err := s.readList()
	if err != nil {
		log.Printf("Error occurred getting file: %v", err)
		return "", err
	}

	if err = s.checkUser(name, s.store); err == i.UserExists {
		return "", i.UserExists
	}

	atomic.AddInt64(&s.store.Increment, 1)
	newUser := User{
		CreatedAt:   time.Now(),
		DisplayName: name,
		Email:       email,
	}

	id := strconv.Itoa(int(s.store.Increment))
	s.store.List[id] = newUser

	if err = s.writeList(); err != nil {
		log.Printf("Error occurred writing to file: %v", err)
		return "", err
	}

	return id, nil
}

// GetUser returns a specific user
func (s *storage) GetUser(id string) ([]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	err := s.readList()
	if err != nil {
		log.Printf("Error occurred getting file: %v", err)
		return nil, err
	}

	if s.store.List[id].DisplayName == "" {
		return nil, i.UserNotFound
	}

	resp, err := json.Marshal(s.store.List[id])
	if err != nil {
		log.Printf("Error occurred marshaling: %v", err)
		return nil, err
	}

	return resp, nil
}

// UpdateUser updates user data if exists
func (s *storage) UpdateUser(id, name, email string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	err := s.readList()
	if err != nil {
		log.Printf("Error occurred getting file: %v", err)
		return err
	}

	if _, ok := s.store.List[id]; !ok {
		return i.UserNotFound
	}

	user := s.store.List[id]
	user.DisplayName, user.Email = name, email

	s.store.List[id] = user

	if err = s.writeList(); err != nil {
		log.Printf("Error occurred writing to file: %v", err)
		return err
	}

	return nil
}

// DeleteUser removes specific user from the list
func (s *storage) DeleteUser(id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	err := s.readList()
	if err != nil {
		log.Printf("Error occurred getting file: %v", err)
		return err
	}

	if _, ok := s.store.List[id]; !ok {
		return i.UserNotFound
	}

	delete(s.store.List, id)

	if err = s.writeList(); err != nil {
		log.Printf("Error occurred writing to file: %v", err)
		return err
	}

	return nil
}
