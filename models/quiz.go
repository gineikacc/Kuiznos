package models

type Quiz struct {
	Id        int               `form:"id" json:"id"`
	Title     string            `form:"title" json:"title"`
	Author    string            `form:"author" json:"author"`
	Questions []*Question_entry `form:"questions" json:"questions"`
}

type Question_entry struct {
	Id      int                `form:"id" json:"id"`
	Title   string             `form:"title" json:"title"`
	Choices []*Question_choice `form:"choices" json:"choices"`
}

type Question_choice struct {
	Id         int    `form:"id" json:"id"`
	Title      string `form:"title" json:"title"`
	Is_correct bool   `form:"correct" json:"correct"`
}
