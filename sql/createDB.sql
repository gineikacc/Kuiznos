CREATE DATABASE IF NOT EXISTS temp_db
    DEFAULT CHARACTER SET = 'utf8mb4';

CREATE TABLE IF NOT EXISTS customer (
  id int PRIMARY KEY AUTO_INCREMENT,
  username varchar(64) NOT NULL,
  auth_password char(128) NOT NULL,
  email varchar(64) NOT NULL,
  quiz_done_count int NOT NULL
);

CREATE TABLE IF NOT EXISTS question (
  id int PRIMARY KEY AUTO_INCREMENT,
  title varchar(64) NOT NULL,
  quiz_id int NOT NULL,
  FOREIGN KEY (quiz_id) REFERENCES quiz (id) 
) ;

CREATE TABLE IF NOT EXISTS question_choice  (
  id int PRIMARY KEY AUTO_INCREMENT,
  choice varchar(128) NOT NULL,
  is_correct BOOLEAN NOT NULL DEFAULT FALSE,
  question_id int,
  FOREIGN KEY (question_id) REFERENCES question (id) 
) ;

CREATE TABLE IF NOT EXISTS quiz (
  id int PRIMARY KEY AUTO_INCREMENT,
  title varchar(64) NOT NULL,
  author int,
  FOREIGN KEY (author) REFERENCES customer (id)
) ;
