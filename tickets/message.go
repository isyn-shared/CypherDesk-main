package tickets

import (
	"time"
)

type Message struct {
	Event          string    `json:"event"`
	Account        *Account  `json:"account"`
	Data           string    `json:"data"`
	RecipientLogin string    `json:"reclogin"`
	Time           time.Time `json:"time"`
}
