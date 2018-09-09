package db

import (
	"database/sql"
	"time"
)

// Ticket structure
type Ticket struct {
	ID          int    `json:"ID"`
	Caption     string `json:"Caption"`
	Description string `json:"Description"`
	Sender      int    `json:"Sender"`
	Status      string `json:"Status"`
}

// TicketLog struct))))
type TicketLog struct {
	ID       int       `json: "ID"`
	Ticket   int       `json: "Ticket"`
	UserFrom int       `json: "UserFrom"`
	UserTo   int       `json: "UserTo"`
	Action   string    `json: "Action"`
	Time     time.Time `json: "Time"`
}

// GetTicket returns ticket obj from db using id
func (m *MysqlUser) GetTicket(id int) *Ticket {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM tickets WHERE id = ? LIMIT 1")
	defer stmt.Close()

	ticket := new(Ticket)
	err := stmt.QueryRow(id).Scan(&ticket.ID, &ticket.Caption, &ticket.Description, &ticket.Sender, &ticket.Status)

	if err != nil {
		panic("db error: " + err.Error())
	}
	return ticket
}

// CreateTicket creates new ticket in DB
func (m *MysqlUser) CreateTicket(ticket *Ticket) sql.Result {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "INSERT INTO tickets (caption, description, sender, status) VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	res := exec(stmt, []interface{}{ticket.Caption, ticket.Description, ticket.Sender, ticket.Status})
	return res
}

// GetTicketLog return array of logs from DB using ticket ID
func (m *MysqlUser) GetTicketLog(ticketID int) []*TicketLog {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM tickets WHERE ticket = ?")
	defer stmt.Close()

	rows := chk(stmt.Query(ticketID)).(*sql.Rows)
	logs := make([]*TicketLog, 0)
	for rows.Next() {
		log := new(TicketLog)
		err := rows.Scan()
		if err != nil {
			panic("GetTicketLog error: " + err.Error())
		}
		logs = append(logs, log)
	}

	return logs
}

// TransferTicket pass ticket to another user
func (m *MysqlUser) TransferTicket(newLog *TicketLog) sql.Result {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "INSERT INTO logs (ticket, userFrom, userTo, action, time) VALUES (?, ?, ?, ?, ?)")
	defer stmt.Close()

	res := exec(stmt, []interface{}{newLog.Ticket, newLog.UserFrom, newLog.UserTo, newLog.Action, newLog.Time})
	return res
}

// func (m *MysqlUser) GetUserTickets(user *User) []*Users {

// }
