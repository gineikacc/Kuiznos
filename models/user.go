package models

type User struct {
	Name            string `form:"name"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	Quiz_done_count int    `form:"quiz_done_count"`
}
