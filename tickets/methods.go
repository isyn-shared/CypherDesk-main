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
	tickets := mysql.GetUserTickets(id, true)
	byteJSONTickets := chk(json.Marshal(tickets)).([]byte)
	sendResponse(true, "get", string(byteJSONTickets), chnMsg.conn)
}

func sendModeratorTicket(chnMsg *chanMessage) {
	mysql := db.CreateMysqlUser()
	id := chnMsg.Message.Account.ID
	user := mysql.GetUser("id", id)
	if !user.Exist() || !user.Filled() {
		sendResponse(false, "createM", "У Вас нет прав на это действие!", chnMsg.conn)
		return
	}
	args := make(map[string]string)
	err := json.Unmarshal([]byte(chnMsg.Message.Data), &args)
	if err != nil {
		sendResponse(false, "createM", "Некорректный формат входных данных", chnMsg.conn)
		return
	}
	caption, description, toIdStr := args["caption"], args["description"], args["id"]
	if alias.EmptyStr(caption) || alias.EmptyStr(description) || alias.EmptyStr(toIdStr) {
		sendResponse(false, "createM", "Неправильный запрос", chnMsg.conn)
		return
	}
	toId, err := alias.STI(toIdStr)
	if err != nil {
		sendResponse(false, "createM", "Невозможное значение ID", chnMsg.conn)
		return
	}

	toUser := mysql.GetUser("id", toId)
	if !user.Exist() {
		sendResponse(false, "createM", "Такого пользователя не существует", chnMsg.conn)
	}

	ticket := &db.Ticket{
		Caption:     caption,
		Description: description,
		Sender:      id,
		Status:      "opened",
	}

	tmpTicket := *ticket
	ntID := mysql.CreateTicket(ticket)
	tmpTicket.ID = ntID

	log := &db.TicketLog{
		Ticket:   mysql.GetLastTicketBySender(id),
		UserFrom: id,
		UserTo:   toId,
		Action:   "send",
		Time:     time.Now(),
	}

	extTicket := db.ExtTicket{
		Ticket:      &tmpTicket,
		ForwardFrom: id,
		ForwardTo:   toId,
		Time:        time.Now(),
	}

	mysql.TransferTicket(log)

	if ClientsByLogin[toUser.Login] != nil {
		sendResponse(true, "incoming", string(chk(json.Marshal(extTicket)).([]byte)), ClientsByLogin[toUser.Login].Connection)
	}
	sendResponse(true, "create", string(chk(json.Marshal(extTicket)).([]byte)), chnMsg.conn)
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
	if alias.EmptyStr(caption) || alias.EmptyStr(description) {
		sendResponse(false, "create", "Неправильный запрос", chnMsg.conn)
		return
	}
	ticket := &db.Ticket{
		Caption:     caption,
		Description: description,
		Sender:      id,
		Status:      "opened",
	}

	tmpTicket := *ticket
	ntID := mysql.CreateTicket(ticket)
	tmpTicket.ID = ntID
	ticketAdmin := mysql.GetDepartmentTicketAdmin(user.Department)

	log := &db.TicketLog{
		Ticket:   mysql.GetLastTicketBySender(id),
		UserFrom: id,
		UserTo:   ticketAdmin.ID,
		Action:   "send",
		Time:     time.Now(),
	}

	extTicket := db.ExtTicket{
		Ticket:      &tmpTicket,
		ForwardFrom: id,
		ForwardTo:   ticketAdmin.ID,
		Time:        time.Now(),
	}

	mysql.TransferTicket(log)

	if ClientsByLogin[ticketAdmin.Login] != nil {
		sendResponse(true, "incoming", string(chk(json.Marshal(extTicket)).([]byte)), ClientsByLogin[ticketAdmin.Login].Connection)
	}
	sendResponse(true, "create", string(chk(json.Marshal(extTicket)).([]byte)), chnMsg.conn)
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
	err := json.Unmarshal([]byte(chnMsg.Message.Data), &args)
	if err != nil {
		sendResponse(false, "forward", "Ошибка на сервере"+err.Error(), chnMsg.conn)
		return
	}
	// TODO: rights to send ticket to this user

	to, ticketID := args["to"], args["ticketID"]
	if alias.EmptyStr(to) || alias.EmptyStr(ticketID) {
		sendResponse(false, "forward", "Неправильный запрос!", chnMsg.conn)
		return
	}
	tID := chk(alias.STI(ticketID)).(int)
	ticket := mysql.GetTicket(tID)

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
		ClientsByLogin[chnMsg.Message.RecipientLogin].Connection.WriteJSON(StandartResponse{"event": "incoming", "ok": true, "data": ticketStr})
	}
	sendResponse(true, "forward", "null", chnMsg.conn)
}

func closeTicket(chnMsg *chanMessage) {
    mysql := db.CreateMysqlUser()
    id := chnMsg.Message.Account.ID
    user := mysql.GetUser("id", id)
    if !user.Exist() || !user.Filled() {
        sendResponse(false, "close", "У Вас нет прав на это действие!", chnMsg.conn)
        return
    }
    args := make(map[string]string)
    err := json.Unmarshal([]byte(chnMsg.Message.Data), &args)
    if err != nil {
        sendResponse(false, "close", "Ошибка на сервере", chnMsg.conn)
        return
    }

    strTicketID := args["id"]
    if alias.EmptyStr(strTicketID) {
        sendResponse(false, "close", "Неправильный запрос", chnMsg.conn)
    }
    
    ticketID := chk(alias.STI(strTicketID)).(int)
    ticket := mysql.GetTicket(ticketID)

    sender := mysql.GetUser("id", ticket.Sender)
    mysql.UpdateTicketStatus(ticketID, "closed")

    ticketStr := string(chk(json.Marshal(ticket)).([]byte))

    if ClientsByLogin[sender.Login] != nil {
        ClientsByLogin[sender.Login].Connection.WriteJSON(StandartResponse{"event": "closedTicket", "ok": true, "data": ticketStr})
    }
    sendResponse(true, "close", ticketStr, chnMsg.conn)
}
