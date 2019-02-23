package tickets

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type EventArguments map[string]string

func ticketMethodsRecovery(chnMsg *chanMessage, eventName string) {
	if r := recover(); r != nil {
		fmt.Printf("%T\n %q", r, r)
		log.Fatal("Error in event " + eventName + ": ")
		// sendResponse(false, eventName, r.(string), chnMsg.conn)
	}
}

func getEventArgs(chnMsg *chanMessage) EventArguments {
	var args EventArguments
	err := json.Unmarshal([]byte(chnMsg.Message.Data), &args)
	if err != nil {
		panic("Некорректный формат входных данных")
	}
	return args
}

func exchangePublicKeys(chnMsg *chanMessage) {
	defer ticketMethodsRecovery(chnMsg, "publicKey")
	args := getEventArgs(chnMsg)
	clientPubKey := getPublicKeyFromPem(args["key"])

	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", chnMsg.Message.Account.ID)

	response := make(map[string][]byte)
	base64ServerKey := []byte(alias.Base64Enc(string(ClientsByLogin[user.Login].ServerKey)))
	base64ClientKey := []byte(alias.Base64Enc(string(ClientsByLogin[user.Login].ClientKey)))
	response["server"] = encryptWithPublicKey(base64ServerKey, clientPubKey)
	response["client"] = encryptWithPublicKey(base64ClientKey, clientPubKey)

	fmt.Println("IM HERE!")
	clientsBySocket[chnMsg.conn].SecureConnection = true

	sendResponse(true, "publicKey", string(chk(json.Marshal(response)).([]byte)), chnMsg.conn)
}

func getTickets(chnMsg *chanMessage) {
	defer ticketMethodsRecovery(chnMsg, "get")
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", chnMsg.Message.Account.ID)
	if !user.Exist() || !user.Filled() {
		return
	}
	tickets := mysql.GetUserTickets(user.ID, true)
	byteJSONTickets := chk(json.Marshal(tickets)).([]byte)
	sendResponse(true, "get", string(byteJSONTickets), chnMsg.conn)
}

func sendTicket(userFrom *db.User, userTo *db.User, args EventArguments, chnMsg *chanMessage) {
	mysql := db.CreateMysqlUser()

	ticket := &db.Ticket{
		Caption:     args["caption"],
		Description: args["description"],
		Sender:      userFrom.ID,
		Status:      "opened",
	}

	tmpTicket := *ticket
	ntID := mysql.CreateTicket(ticket)
	tmpTicket.ID = ntID

	log := &db.TicketLog{
		Ticket:   mysql.GetLastTicketBySender(userFrom.ID),
		UserFrom: userFrom.ID,
		UserTo:   userTo.ID,
		Action:   "send",
		Time:     time.Now(),
	}

	extTicket := db.ExtTicket{
		Ticket:      &tmpTicket,
		ForwardFrom: userFrom.ID,
		ForwardTo:   userTo.ID,
		Time:        time.Now(),
	}

	mysql.TransferTicket(log)

	if ClientsByLogin[userTo.Login] != nil {
		sendResponse(true, "incoming", string(chk(json.Marshal(extTicket)).([]byte)), ClientsByLogin[userTo.Login].Connection)
	}
	sendResponse(true, "create", string(chk(json.Marshal(extTicket)).([]byte)), chnMsg.conn)
}

func sendModeratorTicket(chnMsg *chanMessage) {
	defer ticketMethodsRecovery(chnMsg, "createM")

	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", chnMsg.Message.Account.ID)
	if !user.Exist() || !user.Filled() {
		panic("У Вас нет прав на это действие!")
	}

	args := getEventArgs(chnMsg)
	if args["caption"] == "" || args["description"] == "" || args["id"] == "" {
		panic("Неправильный запрос")
	}
	toID, err := alias.STI(args["id"])
	if err != nil {
		panic("Невозможное значение ID")
	}

	toUser := mysql.GetUser("id", toID)
	if !user.Exist() {
		panic("Такого пользователя не существует")
	}

	sendTicket(user, toUser, args, chnMsg)
}

func sendUserTicket(chnMsg *chanMessage) {
	defer ticketMethodsRecovery(chnMsg, "create")
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", chnMsg.Message.Account.ID)
	if !user.Exist() || !user.Filled() {
		panic("У Вас нет прав на это действие!")
	}
	args := getEventArgs(chnMsg)
	if args["caption"] == "" || args["description"] == "" {
		panic("Неправильный запрос")
	}

	ticketAdmin := mysql.GetDepartmentTicketAdmin(user.Department)
	sendTicket(user, ticketAdmin, args, chnMsg)
}

func forwardTicket(chnMsg *chanMessage) {
	defer ticketMethodsRecovery(chnMsg, "forward")
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", chnMsg.Message.Account.ID)
	if !user.Exist() || !user.Filled() || user.Role != "ticketModerator" {
		panic("У Вас нет прав на это действие!")
	}
	args := getEventArgs(chnMsg)
	// TODO: rights to send ticket to this user

	to, ticketID := args["to"], args["ticketID"]
	if args["to"] == "" || args["ticketID"] == "" {
		panic("Неправильный закон!")
	}
	tID := chk(alias.STI(ticketID)).(int)
	ticket := mysql.GetTicket(tID)

	if !ticket.Exist() {
		panic("Неправильный запрос!")
	}

	log := &db.TicketLog{
		Ticket:   tID,
		UserFrom: user.ID,
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
	defer ticketMethodsRecovery(chnMsg, "close")
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", chnMsg.Message.Account.ID)
	if !user.Exist() || !user.Filled() {
		panic("У Вас нет прав на это действие!")
	}
	args := getEventArgs(chnMsg)

	if args["id"] == "" {
		sendResponse(false, "close", "Неправильный запрос", chnMsg.conn)
	}

	ticketID := chk(alias.STI(args["id"])).(int)
	ticket := mysql.GetTicket(ticketID)

	sender := mysql.GetUser("id", ticket.Sender)
	mysql.UpdateTicketStatus(ticketID, "closed")

	ticketStr := string(chk(json.Marshal(ticket)).([]byte))

	if ClientsByLogin[sender.Login] != nil {
		ClientsByLogin[sender.Login].Connection.WriteJSON(StandartResponse{"event": "closedTicket", "ok": true, "data": ticketStr})
	}
	sendResponse(true, "close", ticketStr, chnMsg.conn)
}
