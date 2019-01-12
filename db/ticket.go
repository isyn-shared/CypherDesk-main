package db

import (
	"CypherDesk-main/alias"
	"database/sql"
	"log"
	"reflect"
	"time"
)

const (
	ticketKey = "keys/ticketKey.toml"
)

// Ticket structure
type Ticket struct {
	ID          int    `json:"ID"`
	Caption     string `json:"Caption"`
	Description string `json:"Description"`
	Sender      int    `json:"Sender"`
	Status      string `json:"Status"`
}

type ExtTicket struct {
	Ticket      *Ticket   `json:"ticket"`
	ForwardFrom int       `json:"forwardFrom"`
	ForwardTo   int       `json:"forwardTo"`
	Action      string    `json:"action"`
	Time        time.Time `json:"time"`
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

/*
	About ticketLog Action
	when user creates ticket, creates ticketLog with action "send".
	when user forward ticket, creates ticketLog with action "forward".
	when user fillfield the ticket, creates ticketLog with action "result" and userTo = "Ticket.sender"
*/

func (t *Ticket) Exist() bool {
	if t.ID == 0 {
		return false
	}
	return true
}

func (et *ExtTicket) Exist() bool {
	if et.Ticket.ID == 0 {
		return false
	}
	return true
}

func refactTicketKey(val string, dec bool) string {
	return alias.StandartRefact(val, dec, ticketKey)
}

func (t *Ticket) refact(dec bool) {
	fields := reflect.TypeOf(*t)
	values := reflect.ValueOf(*t)

	num := values.NumField()

	for i := 0; i < num; i++ {
		var input string
		value := values.Field(i)
		field := fields.Field(i)

		switch value.Kind() {
		case reflect.String:
			input = value.String()

			var enc string
			enc = refactDepartmentField(input, dec)

			reflect.ValueOf(t).Elem().FieldByName(field.Name).SetString(enc)
		case reflect.Int:
			continue
		}
	}
}

// UpdateTicketStatus updates status of ticket in db
func (m *MysqlUser) UpdateTicketStatus(ticketID int, status string) int64 {
	status = refactTicketKey(status, false)

	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "UPDATE tickets SET status = ? WHERE id = ?")
	defer stmt.Close()

	res := exec(stmt, []interface{}{status, ticketID})
	aff := affect(res)
	return aff
}

// GetTicket returns ticket obj from db using id
func (m *MysqlUser) GetTicket(id int) *Ticket {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM tickets WHERE id = ? LIMIT 1")
	defer stmt.Close()

	ticket := new(Ticket)
	err := stmt.QueryRow(id).Scan(&ticket.ID, &ticket.Caption, &ticket.Description, &ticket.Sender, &ticket.Status)

	if err != nil && err.Error() == "sql: no rows in result set" {
		return ticket
	}
	if err != nil {
		panic("db error: " + err.Error())
	}

	ticket.refact(true)

	return ticket
}

// CreateTicket creates new ticket in DB
func (m *MysqlUser) CreateTicket(ticket *Ticket) int {
	ticket.refact(false)
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "INSERT INTO tickets (caption, description, sender, status) VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	exec(stmt, []interface{}{ticket.Caption, ticket.Description, ticket.Sender, ticket.Status})

	var res int
	stmt = prepare(db, "SELECT LAST_INSERT_ID()")
	err := stmt.QueryRow().Scan(&res)

	if err != nil {
		log.Fatal("Error when insert new ticket in db.CreateTicket")
	}

	return res
}

// GetTicketLog return array of logs from DB using ticketID
func (m *MysqlUser) GetTicketLog(ticketID int) []*TicketLog {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM logs WHERE ticket = ?")
	defer stmt.Close()

	rows := chk(stmt.Query(ticketID)).(*sql.Rows)
	logs := make([]*TicketLog, 0)
	for rows.Next() {
		log := new(TicketLog)
		err := rows.Scan(&log.ID, &log.Ticket, &log.UserFrom, &log.UserTo, &log.Action, &log.Time)
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

// GetUserTickets returns array of tickets, which sended to user
func (m *MysqlUser) GetUserTickets(userID int, filter bool) []*ExtTicket {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM logs WHERE (userTo = ? OR userFrom = ?)")
	defer stmt.Close()

	rows := chk(stmt.Query([]interface{}{userID, userID}...)).(*sql.Rows)
	tickets := make([]*ExtTicket, 0)
	for rows.Next() {
		log := new(TicketLog)
		err := rows.Scan(&log.ID, &log.Ticket, &log.UserFrom, &log.UserTo, &log.Action, &log.Time)
		if err != nil {
			panic("GetUserTickets error: " + err.Error())
		}
		exT := new(ExtTicket)
		exT.Ticket = m.GetTicket(log.Ticket)
		exT.Action, exT.ForwardFrom, exT.ForwardTo, exT.Time = log.Action, log.UserFrom, log.UserTo, log.Time

		if (filter && exT.Ticket.Status != "deleted") || !filter {
			tickets = append(tickets, exT)
		}
	}
	return tickets
}

// GetLastLogId - returns last
func (m *MysqlUser) GetLastTicketBySender(senderID int) int {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT MAX(id) FROM tickets WHERE sender = ?")
	defer stmt.Close()

	var res int
	err := stmt.QueryRow(senderID).Scan(&res)

	if err != nil {
		panic("db error: " + err.Error())
	}
	return res
}
