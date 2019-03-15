package chat

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type StandartResponse map[string]interface{}

var ClientsByLogin = make(map[string]*Client)
var clientsBySocket = make(map[*websocket.Conn]*Client)
var myEvents = make(map[string]func(*chanMessage))

type chanMessage struct {
	Message *Message
	conn    *websocket.Conn
}

var messages = make(chan chanMessage)
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Start() {
	bindEvents()
	go handleMessages()
}

func bindEvents() {
	myEvents["send"] = func(chMsg *chanMessage) {
		sendUserChatMessage(chMsg)
	}
	myEvents["get"] = func(chMsg *chanMessage) {
		getChatUserMessages(chMsg)
	}
}

func handleMessages() {
	for {
		ChanMsg := <-messages
		ChanMsg.Message.Account = clientsBySocket[ChanMsg.conn].Account

		if myEvents[ChanMsg.Message.Event] == nil {
			sendResponse(false, "error", "Обращение к несуществующему event-у", ChanMsg.conn)
			continue
		}

		myEvents[ChanMsg.Message.Event](&ChanMsg)
	}
}

func HandleConnections(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("handle connection error: " + err.Error())
	}
	defer conn.Close()

	isAuthorized, id := getID(c)
	if !isAuthorized {
		sendResponse(false, "error", "Вы не авторизованы!", conn)
		return
	}

	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)

	if !user.Exist() || !user.Filled() {
		sendResponse(false, "error", "У Вас нет прав на это действие!", conn)
		return
	}

	account := &Account{
		ID:         id,
		Login:      user.Login,
		Mail:       user.Mail,
		Name:       user.Name,
		Surname:    user.Surname,
		Partonymic: user.Partonymic,
		Role:       user.Role,
	}

	addClient(conn, account)

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			deleteClient(clientsBySocket[conn])
			break
		}

		messages <- chanMessage{&msg, conn}
		clientsBySocket[conn].SecureConnection = true
	}
}

func deleteClient(client *Client) {
	delete(ClientsByLogin, client.Account.Login)
	delete(clientsBySocket, client.Connection)
	// updateFriendsStat(client)
}

func addClient(conn *websocket.Conn, acc *Account) {
	serverKey, clientKey := alias.GenAESKey(), alias.GenAESKey()
	client := &Client{
		Account:    acc,
		Connection: conn,
		ClientKey:  clientKey,
		ServerKey:  serverKey,
	}
	clientsBySocket[conn] = client
	ClientsByLogin[acc.Login] = client
	ClientsByLogin[acc.Login].SecureConnection = false
}

func sendResponse(ok bool, event string, message string, conn *websocket.Conn) {
	err := conn.WriteJSON(StandartResponse{"ok": ok, "data": message, "event": event})

	if err != nil {
		fmt.Println("handleMessage error: " + err.Error())
		deleteClient(clientsBySocket[conn])
	}
}