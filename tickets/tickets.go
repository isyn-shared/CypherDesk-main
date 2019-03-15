package tickets

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type StandartResponse map[string]interface{}

var ClientsByLogin = make(map[string]*Client)
var clientsBySocket = make(map[*websocket.Conn]*Client)

type chanMessage struct {
	Message *Message
	conn    *websocket.Conn
}

type systemChanMessage struct {
	systemMessage *systemMessage
	conn          *websocket.Conn
}

var messages = make(chan chanMessage)
var systemMessages = make(chan systemChanMessage)
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
	// go handleSystemMessages()
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
		if clientsBySocket[conn].SecureConnection {
			_, p, err := conn.ReadMessage()
			if err != nil {
				deleteClient(clientsBySocket[conn])
				break
			}
			decryptedMsg, _ := alias.DecryptAESCBC(string(p), clientsBySocket[conn].ClientKey)
			decryptedMsg = decryptedMsg[:bytes.LastIndex(decryptedMsg, []byte("}"))+1]
			err = json.Unmarshal(decryptedMsg, &msg)
			if err != nil {
				sendResponse(false, "null", "Invalid data", conn)
				continue
			}
		} else {
			err := conn.ReadJSON(&msg)
			if err != nil {
				deleteClient(clientsBySocket[conn])
				break
			}
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

var myEvents = make(map[string]func(*chanMessage))

func bindEvents() {
	myEvents["create"] = func(chMsg *chanMessage) {
		sendUserTicket(chMsg)
	}
	myEvents["forward"] = func(chMsg *chanMessage) {
		forwardTicket(chMsg)
	}
	myEvents["get"] = func(chMsg *chanMessage) {
		getTickets(chMsg)
	}
	myEvents["close"] = func(chMsg *chanMessage) {
		closeTicket(chMsg)
	}
	myEvents["createM"] = func(chMsg *chanMessage) {
		sendModeratorTicket(chMsg)
	}
	myEvents["publicKey"] = func(chMsg *chanMessage) {
		exchangePublicKeys(chMsg)
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

func handleSystemMessages() {
	for {
		SysMsg := <-systemMessages
		switch SysMsg.systemMessage.Event {
		case "SecureConnectionStatus":
			clientsBySocket[SysMsg.conn].SecureConnection = SysMsg.systemMessage.Value.(bool)
		}
	}
}

func sendResponse(ok bool, event string, message string, conn *websocket.Conn) {
	var err error
	if event == "publicKey" {
		err = conn.WriteJSON(StandartResponse{"ok": ok, "data": message, "event": event})
	} else {
		response, err := json.Marshal(&StandartResponse{"ok": ok, "data": message, "event": event})
		if err != nil {
			fmt.Println("Error in marshaling StandartResponse object!!!")
			return
		}
		encResp, err := alias.EncryptAESCBC(response, clientsBySocket[conn].ServerKey)
		if err != nil {
			fmt.Println("Error when decrypting", err.Error())
		}
		fmt.Println("ENCRYPTED: ", string(encResp))
		err = conn.WriteMessage(1, encResp)
	}
	if err != nil {
		fmt.Println("handleMessage error: " + err.Error())
		deleteClient(clientsBySocket[conn])
	}
}
