
package main

import (
    "sync"
    "errors"
)

type User struct {
    ID       uint   `json:"id"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password,omitempty" binding:"required,min=8"`
    Role     string `json:"role"`
}

type UserStore struct {
    users map[string]*User
    mu    sync.RWMutex
    nextID uint
}

func NewUserStore() *UserStore {
    return &UserStore{
        users: make(map[string]*User),
        nextID: 1,
    }
}


func (s *UserStore) Create(user *User) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.users[user.Email]; exists {
        return errors.New("user already exists")
    }

    user.ID = s.nextID
    s.nextID++
    s.users[user.Email] = user
    return nil
}

func (s *UserStore) GetByEmail(email string) (*User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    user, exists := s.users[email]
    if !exists {
        return nil, errors.New("user not found")
    }
    return user, nil
}

var userStore = NewUserStore()
