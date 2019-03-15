package chat

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type EventArguments map[string]string

func chatMethodsRecovery(chnMsg *chanMessage, eventName string) {
	if r := recover(); r != nil {
		fmt.Printf("%T\n %q", r, r)
		log.Fatal("Error in chat event " + eventName + ": ")
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

func sendUserChatMessage(chnMsg *chanMessage) {
	defer chatMethodsRecovery(chnMsg, "send")

	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", chnMsg.Message.Account.ID)

	if !user.Exist() || !user.Filled() {
		panic("У Вас нет прав на это действие!")
	}

	args := getEventArgs(chnMsg)
	if args["text"] == "" || args["to"] == "" {
		panic("Неправильный запрос")
	}

	toID, err := alias.STI(args["to"])
	if err != nil {
		panic("Невозможное значение ID")
	}

	toUser := mysql.GetUser("id", toID)
	if !toUser.Exist() {
		panic("Такого пользователя не существует")
	}

	chatMessage := &db.ChatMsg{
		From:   user.ID,
		To:     toID,
		Date:   time.Now(),
		Text:   args["text"],
		Status: 0,
	}

	chatMessage.ID = mysql.InsertChatMessage(chatMessage)

	if ClientsByLogin[toUser.Login] != nil {
		sendResponse(true, "newMessage", string(chk(json.Marshal(chatMessage)).([]byte)), ClientsByLogin[toUser.Login].Connection)
	}
	sendResponse(true, "newMessage", string(chk(json.Marshal(chatMessage)).([]byte)), chnMsg.conn)
}

func getChatUserMessages(chnMsg *chanMessage) {
	defer chatMethodsRecovery(chnMsg, "get")

	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", chnMsg.Message.Account.ID)

	if !user.Exist() || !user.Filled() {
		panic("У Вас нет прав на это действие!")
	}

	chatMessages := mysql.GetUsersChatMessages(user)
	sendResponse(true, "get", string(chk(json.Marshal(chatMessages)).([]byte)), chnMsg.conn)
}
