package models

type Quiz struct {
	id        int
	Title     string
	Author    User
	Questions []Question
}

type Question struct {
	Title          string
	Choices        []string
	Correct_choice int
}
