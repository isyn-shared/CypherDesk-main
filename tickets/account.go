package tickets

import (
	"github.com/gorilla/websocket"
)

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
