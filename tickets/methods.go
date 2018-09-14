package tickets

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"encoding/json"
	"time"
)

func getTickets(chnMsg *chanMessage) {
	mysql := db.CreateMysqlUser()
	id := chnMsg.Message.Account.ID
	user := mysql.GetUser("id", id)
	if !user.Exist() || !user.Filled() {
		sendResponse(false, "get", "У Вас нет прав на это действие!", chnMsg.conn)
		return
	}
	tickets := mysql.GetUserTickets(id)
	byteJsonTickets := chk(json.Marshal(tickets)).([]byte)
	sendResponse(true, "get", string(byteJsonTickets), chnMsg.conn)
}

func sendTicket(chnMsg *chanMessage) {
	mysql := db.CreateMysqlUser()
	id := chnMsg.Message.Account.ID
	user := mysql.GetUser("id", id)
	if !user.Exist() || !user.Filled() {
		sendResponse(false, "create", "У Вас нет прав на это действие!", chnMsg.conn)
		return
	}
	args := make(map[string]string)
	err := json.Unmarshal([]byte(chnMsg.Message.Data), &args)
	if err != nil {
		sendResponse(false, "create", "Ошибка на сервере", chnMsg.conn)
		return
	}
	caption, description := args["caption"], args["description"]
	if alias.EmptyStr(caption) || alias.EmptyStr("description") { // TODO: chk for nil????
		sendResponse(false, "create", "Неправильный запрос", chnMsg.conn)
		return
	}
	ticket := &db.Ticket{
		Caption:     caption,
		Description: description,
		Sender:      id,
		Status:      "opened",
	}
	mysql.CreateTicket(ticket)
	ticketAdmin := mysql.GetDepartmentTicketAdmin(user.Department)
	log := &db.TicketLog{
		Ticket:   mysql.GetLastTicketBySender(id),
		UserFrom: id,
		UserTo:   ticketAdmin.ID,
		Action:   "send",
		Time:     time.Now(),
	}

	reciever := mysql.GetUser("id", ticketAdmin.ID)
	if ClientsByLogin[reciever.Login] != nil {
		sendResponse(true, "create", string(chk(json.Marshal(ticket)).([]byte)), ClientsByLogin[reciever.Login].Connection)
	}

	mysql.TransferTicket(log)
	sendResponse(true, "create", "null", chnMsg.conn)
}

func forwardTicket(chnMsg *chanMessage) {
	mysql := db.CreateMysqlUser()
	id := chnMsg.Message.Account.ID
	user := mysql.GetUser("id", id)
	if !user.Exist() || !user.Filled() || user.Role != "ticketModerator" {
		sendResponse(false, "forward", "У Вас нет прав на это действие!", chnMsg.conn)
		return
	}
	args := make(map[string]string)
	err := json.Unmarshal([]byte(chnMsg.Message.Data), args)
	if err != nil {
		sendResponse(false, "forward", "Ошибка на сервере", chnMsg.conn)
		return
	}
	to, ticketID := args["to"], args["ticketID"]
	if alias.EmptyStr(to) || alias.EmptyStr(ticketID) {
		sendResponse(false, "forward", "Неправильный запрос!", chnMsg.conn)
		return
	}
	tID := chk(alias.STI(ticketID)).(int)
	ticket := mysql.GetTicket(id)

	if !ticket.Exist() {
		sendResponse(false, "forward", "Неправильный запрос!", chnMsg.conn)
		return
	}

	log := &db.TicketLog{
		Ticket:   tID,
		UserFrom: id,
		UserTo:   chk(alias.STI(to)).(int),
		Action:   "forward",
		Time:     time.Now(),
	}
	mysql.TransferTicket(log)

	ticketStr := string(chk(json.Marshal(ticket)).([]byte))
	if ClientsByLogin[chnMsg.Message.RecipientLogin] != nil {
		ClientsByLogin[chnMsg.Message.RecipientLogin].Connection.WriteJSON(StandartResponse{"event": "newTicket", "ok": true, "data": ticketStr})
	}

	sendResponse(true, "forward", "null", chnMsg.conn)
}
