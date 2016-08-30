package main

import (
	"database/sql"
	"log"
	"os/user"

	_ "github.com/mattn/go-sqlite3"
)

func db_init() *sql.DB {
	usr, _ := user.Current()

	db, err := sql.Open("sqlite3", usr.HomeDir+"/.config/disnote.db")
	if err != nil {
		log.Println(err)
	}

	sqlStmt := `
		create table note (
				id integer not null primary key,
				text text,
				user_id text not null
		);
		`

	_, _ = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Println("%q: %s", err, sqlStmt)
	// }

	return db
}

func db_insert(db *sql.DB, text, user_id string) (id int64) {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	stmt, err := tx.Prepare("insert into note(text, user_id) values(?, ?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(text, user_id)
	tx.Commit()

	id, _ = result.LastInsertId()

	return id
}

type Rows struct {
	id      int
	text    string
	user_id string
}

func db_select(db *sql.DB) (select_rows []Rows) {
	rows, err := db.Query("select id, text, user_id from note")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var text string
		var user_id string
		err = rows.Scan(&id, &text, &user_id)

		select_rows = append(select_rows, Rows{id, text, user_id})
	}

	return select_rows
}

func db_delete(db *sql.DB, id int) bool {
	result, err := db.Exec("delete from note where id = ?", id)
	if err != nil {
		log.Println(err)
		return false
	}
	ii, _ := result.RowsAffected()
	if ii == 0 {
		return false
	}
	return true
}
