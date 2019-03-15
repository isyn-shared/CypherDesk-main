package db

import (
	"CypherDesk-main/alias"
	"database/sql"
	"log"
	"time"
)

const (
	ChatKey = "keys/chatkey.toml"
)

type ChatMsg struct {
	ID     int
	From   int
	To     int
	Date   time.Time
	Text   string
	Status int
}

func (msg *ChatMsg) Refact(dec bool) {
	msg.Text = alias.StandartRefact(msg.Text, dec, ChatKey)
}

func (m *MysqlUser) GetUsersChatMessages(user *User) []*ChatMsg {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM chat_messages WHERE (from_user = ? OR to_user = ?)")
	defer stmt.Close()

	rows := chk(stmt.Query([]interface{}{user.ID, user.ID}...)).(*sql.Rows)
	messages := make([]*ChatMsg, 0)

	for rows.Next() {
		msg := new(ChatMsg)
		err := rows.Scan(&msg.ID, &msg.Text, &msg.Status, &msg.Date, &msg.From, &msg.To)
		if err != nil {
			panic("GetUserMessages error: " + err.Error())
		}
		messages = append(messages, msg)
	}

	return messages
}

func (m *MysqlUser) InsertChatMessage(msg *ChatMsg) int {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "INSERT INTO chat_messages (text, status, date, from_user, to_user) VALUES (?, ?, ?, ?, ?)")
	defer stmt.Close()

	exec(stmt, []interface{}{msg.Text, msg.Status, msg.Date, msg.From, msg.To})

	var insertedID int
	stmt = prepare(db, "SELECT LAST_INSERT_ID()")
	err := stmt.QueryRow().Scan(&insertedID)

	if err != nil {
		log.Fatal("Error when insert new ticket in db.CreateTicket")
	}

	return insertedID
}

func (m *MysqlUser) EditChatMessage(msg *ChatMsg) int64 {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "UPDATE chat_messages SET text=?, status=? WHERE id=?")
	defer stmt.Close()

	res := exec(stmt, []interface{}{msg.Text, msg.Status, msg.Date, msg.ID})
	aff := affect(res)

	return aff
}
