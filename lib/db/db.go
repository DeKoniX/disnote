package db

import (
	"database/sql"
	"log"
	"os/user"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db     *sql.DB
	config string
}

var DataBase = DB{config: "./.config/"}

func init() {
	var err error

	usr, _ := user.Current()
	DataBase.db, err = sql.Open("sqlite3", usr.HomeDir+"/.config/disnote.db")

	if err != nil {
		log.Panicln(err)
	}

	sqlStmt := `
		create table note (
				id integer not null primary key,
				text text,
				user_id text not null
		);
		`
	_, _ = DataBase.db.Exec(sqlStmt)
}

func (DataBase *DB) Insert(text, user_id string) (id int64) {
	tx, err := DataBase.db.Begin()
	if err != nil {
		log.Panicln(err)
	}

	stmt, err := tx.Prepare("insert into note(text, user_id) values(?, ?)")
	if err != nil {
		log.Panicln(err)
	}
	defer stmt.Close()
	result, _ := stmt.Exec(text, user_id)
	tx.Commit()

	id, _ = result.LastInsertId()

	return id
}

type Rows struct {
	Id     int
	Text   string
	UserId string
}

func (DataBase *DB) Select() (select_rows []Rows) {
	rows, err := DataBase.db.Query("select id, text, user_id from note")
	if err != nil {
		log.Panicln(err)
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

func (DataBase *DB) Delete(id int) bool {
	result, err := DataBase.db.Exec("delete from note where id = ?", id)
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
