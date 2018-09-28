package main

import (
	"CypherDesk-main/init"
	"CypherDesk-main/router"
)

const port string = ":8080"

func main() {
	router := router.New()
	initPkg.ProjectInit()
	router.Run(port)
}
