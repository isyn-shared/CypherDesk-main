package initPkg

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/chat"
	"CypherDesk-main/db"
	"CypherDesk-main/tickets"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
)

const (
	TextColorFail    = "\033[91m"
	TextColorTitle   = "\033[36m"
	TextColorBold    = "\033[1m"
	TextColorEnd     = "\033[0m"
	TextColorOKGreen = "\033[92m"
	TextUnderline    = "\033[4m"
	TextColorWarning = "\033[93m"
)

var AESKeys = []string{"activationKey.toml", "departmentKey.toml", "extendedUserKey.toml",
	"idkey.toml", "passkey.toml", "ticketKey.toml", "userdatakey.toml"}

var DEBUG bool = true

func SetMysqlCredentials(login, password, db string) {
	err := alias.WriteToFile([]byte(login+";"+password+";"+db), "keys/mysql.key")
	if err != nil {
		panic("You have not permissions to edit mysql.key file")
	}
}

func initImg() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func GenerateAESKeys() {
	for _, filename := range AESKeys {
		fmt.Print(".")
		key, kErr := alias.RandomHex(12)
		iv, iErr := alias.RandomHex(8)
		if kErr != nil || iErr != nil {
			panic("Error when generate random hex")
		}
		alias.WriteToFile([]byte("key=\""+key+"\"\niv=\""+iv+"\"\n"), "keys/"+filename)
	}
}

// Project init makes basic init
func ProjectInit() {
	initImg()
	db.MysqlInit()
	tickets.Start()
	chat.Start()

	db.EncryptedAdminValue = alias.StandartRefact("admin", false, db.StInfoKey)
	db.EncryptedTicketModeratorValue = alias.StandartRefact("ticketModerator", false, db.StInfoKey)
	db.EncryptedUserValue = alias.StandartRefact("user", false, db.StInfoKey)
}
