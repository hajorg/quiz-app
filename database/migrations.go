package database

import (
	"database/sql"
)

func CreateAndOpenDb(name string) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user(
			id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(50) NOT NULL,
			password VARCHAR(255) NOT NULL
		)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS subjects(
			id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(50) NOT NULL
		)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS questions(
			id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
			subject_id INT NOT NULL,
			FOREIGN KEY (subject_id) REFERENCES subjects(id),
			content TEXT NOT NULL
		)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS answers(
			id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
			question_id INT NOT NULL,
			FOREIGN KEY (question_id) REFERENCES questions(id),
			content TEXT NOT NULL,
			correct BOOLEAN NOT NULL DEFAULT FALSE
		)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scores(
			id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
			subject_id INT NOT NULL,
			FOREIGN KEY (subject_id) REFERENCES subjects(id),
			user_id INT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES user(id)
		)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DROP PROCEDURE IF EXISTS add_created_at_col")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE PROCEDURE add_created_at_col()
			BEGIN 
				IF NOT EXISTS 
					(SELECT * FROM information_schema.COLUMNS WHERE TABLE_NAME = 'user' AND TABLE_SCHEMA = 'quiz' AND COLUMN_NAME = 'created_at')
				THEN 
					ALTER TABLE user ADD created_at DATETIME NOT NULL;
				END IF;
			END;
		`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("call add_created_at_col()")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DROP PROCEDURE IF EXISTS add_updated_at_col")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE PROCEDURE add_updated_at_col()
			BEGIN
				IF NOT EXISTS (
					SELECT * FROM information_schema.COLUMNS WHERE TABLE_NAME = 'user' AND TABLE_SCHEMA = 'quiz' AND COLUMN_NAME = 'updated_at'
				)
				THEN ALTER TABLE user ADD updated_at DATETIME NOT NULL;
				END IF;
			END;
		`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CALL add_updated_at_col()")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DROP PROCEDURE IF EXISTS add_email_col")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE PROCEDURE add_email_col()
			BEGIN
				IF NOT EXISTS (
					SELECT * FROM information_schema.COLUMNS WHERE TABLE_NAME = 'user' AND TABLE_SCHEMA = 'quiz' AND COLUMN_NAME = 'email'
				)
				THEN ALTER TABLE user ADD email VARCHAR(50) NOT NULL UNIQUE;
				END IF;
			END;
		`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DROP PROCEDURE IF EXISTS make_username_unique")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("ALTER TABLE `user` MODIFY COLUMN `username` VARCHAR(255) NOT NULL UNIQUE;")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CALL add_email_col()")
	if err != nil {
		panic(err)
	}
}
