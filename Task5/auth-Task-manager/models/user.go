package models

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"-" bson:"password"`
	Role     string `json:"role" bson:"role"`
}