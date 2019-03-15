package chat

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
