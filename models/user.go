package models

type User struct {
	Name            string `form:"name" json:"name"`
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	Quiz_done_count int    `form:"quiz_done_count" json:"quiz_done_count"`
}
