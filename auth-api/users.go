package main

import "sync"

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

var (
	users []User

	userMutex sync.Mutex

	nextUserID uint = 1
)