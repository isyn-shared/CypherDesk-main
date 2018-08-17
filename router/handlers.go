package router

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	defer rec(c)
	writePongoTemplate("templates/front-page/index.html", pongo2.Context{}, c)
}

func authorizeHandler(c *gin.Context) {
	login, pass := c.PostForm("login"), c.PostForm("pass")
	if alias.EmptyStr(login) || alias.EmptyStr(pass) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Заполните все поля!"})
		return
	}
}

func testHandler(c *gin.Context) {
	defer rec(c)
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("login", "admin")
	str := chk(json.Marshal(user)).([]byte)
	fmt.Println(string(str))
}
