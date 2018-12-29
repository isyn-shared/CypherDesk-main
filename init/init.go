package initPkg

import (
	"CypherDesk-main/tickets"
	"CypherDesk-main/db"
	"image"
	"image/jpeg"
	"image/png"
)

func initImg() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

// Project init makes basic init
func ProjectInit() {
	initImg()
	db.MysqlInit()
	tickets.Start()
}
