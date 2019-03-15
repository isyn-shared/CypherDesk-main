package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

// Message
type Message struct {
	Event          string    `json:"event"`
	Account        *Account  `json:"account"`
	Data           string    `json:"data"`
	RecipientLogin string    `json:"reclogin"`
	Time           time.Time `json:"time"`
}

type systemMessage struct {
	Event string
	Value interface{}
}

// Client includes account structure and connection structure
type Client struct {
	Account          *Account
	Connection       *websocket.Conn
	ClientKey        []byte
	ServerKey        []byte
	SecureConnection bool
}

type Account struct {
	ID         int    `json:"id"`
	Mail       string `json:"mail"`
	Login      string `json:"login"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Partonymic string `json:"partonymic"`
	Role       string `json:"role"`
}
