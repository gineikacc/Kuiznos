package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"myapp/models"
	"os"

	"github.com/go-sql-driver/mysql"
)

type MysqlStore struct {
	d   *sql.DB
	cfg mysql.Config
}

func New() MysqlStore {
	fmt.Printf("%v arsrs", os.Getenv("DB_PASS"))
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	q := `
CREATE DATABASE IF NOT EXISTS temp_db
DEFAULT CHARACTER SET = 'utf8mb4';`
	db.Exec(q)

	db.Exec(fmt.Sprintf("use %v;", os.Getenv("DB_NAME")))

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return MysqlStore{d: db, cfg: cfg}
}
func (store MysqlStore) Setup() error {

	q := `
CREATE TABLE IF NOT EXISTS customer (
id int PRIMARY KEY AUTO_INCREMENT,
username varchar(64) NOT NULL,
auth_password char(128) NOT NULL,
email varchar(64) NOT NULL,
quiz_done_count int NOT NULL
);`
	store.d.Exec(q)
	q = `
CREATE TABLE IF NOT EXISTS question (
id int PRIMARY KEY AUTO_INCREMENT,
title varchar(64) NOT NULL,
quiz_id int NOT NULL,
FOREIGN KEY (quiz_id) REFERENCES quiz (id) 
);`
	store.d.Exec(q)
	q = `
CREATE TABLE IF NOT EXISTS question_choice  (
id int PRIMARY KEY AUTO_INCREMENT,
choice varchar(128) NOT NULL,
is_correct BOOLEAN NOT NULL DEFAULT FALSE,
question_id int,
FOREIGN KEY (question_id) REFERENCES question (id) 
) ;`
	store.d.Exec(q)
	q = `
CREATE TABLE IF NOT EXISTS quiz (
id int PRIMARY KEY AUTO_INCREMENT,
title varchar(64) NOT NULL,
author int,
FOREIGN KEY (author) REFERENCES customer (id)
)`
	store.d.Exec(q)

	return nil
}
func (store MysqlStore) Create_user(u models.User) error {
	//Check if user doesnt already exist in db
	exists := store.User_exists(u.Name)
	if exists {
		return errors.New("user already exists")
	}
	//Add user to db
	q := `INSERT INTO customer (username, auth_password, email, quiz_done_count) VALUES (?, ?, ?, 0);`
	_, err := store.d.Exec(q, u.Name, u.Password, u.Email)
	if err != nil {

		log.Fatal(err)
		return err
	}
	return nil
}
func (store MysqlStore) Read_user(name string) (models.User, error) {
	var u models.User
	q := `SELECT username, auth_password, email, quiz_done_count FROM customer WHERE username = ? `
	rows, err := store.d.Query(q, name)

	if err != nil {
		log.Println("Couldnt get user")
		return models.User{}, err
	}
	defer rows.Close()
	for rows.Next() {

		rows.Scan(&u.Name, &u.Password, &u.Email, &u.Quiz_done_count)
	}
	return u, nil
}
func (store MysqlStore) Increment_user_quiz_done_count(name string) error {
	//Add user to db
	q := `UPDATE customer SET quiz_done_count = quiz_done_count + 1 WHERE customer.username = "?";`
	_, err := store.d.Exec(q, name)
	if err != nil {

		log.Fatal(err)
		return err
	}
	return nil
}
func (store MysqlStore) User_exists(name string) bool {
	rows, err := store.d.Query("SELECT COUNT(*) as count FROM  customer WHERE username = ?", name)
	if err != nil {
		log.Fatal(err.Error())
	}
	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	return count != 0
}
func (store MysqlStore) Create_quiz(quiz models.Quiz, author string) error {
	//Check if user doesnt already exist in db
	exists := store.Quiz_exists(quiz.Title)
	if exists {
		return errors.New("title already taken")
	}
	exists = store.User_exists(author)
	if !exists {
		return errors.New("author doesnt exist")
	}
	//Create quiz

	q := `INSERT INTO quiz (title, author) VALUES (?,
	(SELECT id FROM customer WHERE username = ?));`
	result, err := store.d.Exec(q, quiz.Title, author)
	if err != nil {
		log.Fatal(err)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	quiz.Id = int(id)

	//Create questions
	for _, question := range quiz.Questions {

		q := `INSERT INTO question (title, quiz_id) VALUES (?, ?);`
		result, err := store.d.Exec(q, question.Title, quiz.Id)
		if err != nil {
			log.Fatal(err)
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		question.Id = int(id)

		//Make choices
		for _, choice := range question.Choices {
			q := `INSERT INTO question_choice (choice, is_correct, question_id) VALUES (?, ?, ?);`
			_, err := store.d.Exec(q, choice.Title, choice.Is_correct, question.Id)
			if err != nil {
				log.Fatal(err)
				return err
			}
		}
	}
	return nil
}

func (store MysqlStore) Read_quiz(name string) (models.Quiz, error) {
	var quiz models.Quiz
	q := `SELECT quiz.id, quiz.title, customer.username 
	FROM quiz 
	INNER JOIN customer ON quiz.author = customer.id 
	WHERE quiz.title = ?;`

	rows, err := store.d.Query(q, name)

	if err != nil {
		return quiz, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&quiz.Id, &quiz.Title, &quiz.Author)
	}

	q = `SELECT question.id, question.title
	FROM question 
	INNER JOIN quiz ON quiz.id = question.quiz_id 
	WHERE quiz.id = ?`
	rows, err = store.d.Query(q, quiz.Id)
	if err != nil {
		return quiz, err
	}
	defer rows.Close()
	for rows.Next() {
		qe := models.Question_entry{}
		rows.Scan(&qe.Id, &qe.Title)
		quiz.Questions = append(quiz.Questions, &qe)
	}

	for _, question := range quiz.Questions {
		q = `SELECT question_choice.id, choice, is_correct
		FROM question_choice 
		INNER JOIN question ON question.id = question_choice.question_id 
		WHERE question.id = ?`
		rows, err = store.d.Query(q, question.Id)

		if err != nil {
			return quiz, err
		}
		defer rows.Close()
		for rows.Next() {
			qc := models.Question_choice{}
			rows.Scan(&qc.Id, &qc.Title, &qc.Is_correct)
			question.Choices = append(question.Choices, &qc)

		}
	}
	return quiz, nil
}

func (store MysqlStore) Quiz_exists(title string) bool {
	rows, err := store.d.Query("SELECT COUNT(*) as count FROM quiz WHERE title = ?", title)
	if err != nil {
		log.Fatal(err.Error())
	}
	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	return count != 0
}

func (store MysqlStore) Query_Quizzes(query string) []models.Quiz {
	var quizzes []models.Quiz
	q := `SELECT quiz.id, quiz.title, customer.username  
	FROM quiz  
	INNER JOIN customer ON quiz.author = customer.id  
	WHERE quiz.title 
	LIKE CONCAT('%', ? , '%')
	LIMIT 25;`

	rows, err := store.d.Query(q, query)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer rows.Close()
	fmt.Printf("Looking w/ query (%v) \n", query)
	for rows.Next() {
		var q models.Quiz
		rows.Scan(&q.Id, &q.Title, &q.Author)
		fmt.Printf("Found %v \n", q.Title)
		quizzes = append(quizzes, q)
	}
	return quizzes
}
