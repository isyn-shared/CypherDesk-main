package feedback

import (
	"log"

	"github.com/BurntSushi/toml"
)

type config struct {
	Server   string
	Port     int
	Email    string
	Password string
}

func (c *config) Read() {
	if _, err := toml.DecodeFile("feedback/mailconfig.toml", &c); err != nil {
		log.Fatal(err)
	}
}
