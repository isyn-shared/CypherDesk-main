package db

import (
	"encoding/json"
	"time"
)

type Document struct {
	ID       int
	UserID   int
	Name     string
	Keywords []string
	Date     time.Time
}

type Frame struct {
	DocID        int
	UserID       int
	RelationType int
}

func (m *MysqlUser) InsertDocument(d *Document) {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "INSERT INTO docs (user_id, name, upload_date, keywords) VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	jsonKeywoards, _ := json.Marshal(&d.Keywords)
	exec(stmt, []interface{}{d.UserID, d.Name, d.Date, jsonKeywoards})

	var docID int
	stmt = prepare(db, "SELECT LAST_INSERT_ID()")
	stmt.QueryRow().Scan(&docID)

	stmt = prepare(db, "INSERT INTO frames (user_id, doc_id, relation_type) VALUES (?, ?, ?)")
	exec(stmt, []interface{}{d.UserID, docID, 1})
}
