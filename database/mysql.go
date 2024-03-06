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
title varchar(64) NOT NULL
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
	fmt.Println("TGBN")

	exists := store.User_exists(u.Name)
	if exists {
		return errors.New("user already exists")
	}
	fmt.Println("ZAAAAAA")

	//Add user to db
	q := `INSERT INTO customer (username, auth_password, email, quiz_done_count) VALUES (?, ?, ?, 0);`
	_, err := store.d.Exec(q, u.Name, u.Password, u.Email)
	if err != nil {
		fmt.Println("GJJJJJG")

		log.Fatal(err)
		return err
	}
	fmt.Printf("User %v added", u.Name)
	return nil
}
func (store MysqlStore) ReadUser(name string) (models.User, error) {
	fmt.Println("RRRRRRRRRRRR")

	var u models.User
	q := `SELECT username, auth_password, email, quiz_done_count FROM customer WHERE username = ? `
	rows, err := store.d.Query(q, name)
	fmt.Println("UUUUUUUUUUUUUUu")

	if err != nil {
		log.Println("Couldnt get user")
		return models.User{}, err
	}
	fmt.Println("DDDDDDDDDDDDD")

	defer rows.Close()
	for rows.Next() {
		fmt.Println("KKKKKKKKKKKAAAAAAAAA")

		rows.Scan(&u.Name, &u.Password, &u.Email, &u.Quiz_done_count)
	}
	fmt.Printf("User (%v, %v, %v) retrieved", u.Name, u.Password, u.Email)
	return u, nil
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

func (store MysqlStore) TempTemp(u models.User) error { return fmt.Errorf("TMEP REVMOE") }
