package domain

import "time"

type Task struct {
	ID          string
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}

type User struct {
	ID              string
	Firstname       string
	Lastname        string
	Username        string
	Password        string
	Role            string
}
