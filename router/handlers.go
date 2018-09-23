package router

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	defer rec(c)
	if isAuthorized, _ := getID(c); isAuthorized {
		c.Redirect(http.StatusSeeOther, "/account")
		return
	}
	writePongoTemplate("templates/frontPage/index.html", pongo2.Context{}, c)
}

func authorizeHandler(c *gin.Context) {
	defer rec(c)
	if isAuthorized, _ := getID(c); isAuthorized {
		c.Redirect(http.StatusSeeOther, "/account")
		return
	}
	login, pass := c.PostForm("login"), c.PostForm("pass")
	if alias.EmptyStr(login) || alias.EmptyStr(pass) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Заполните все поля!"})
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("login", login)

	if !user.Exist() {
		user = mysql.GetUser("mail", login)
		if !user.Exist() {
			c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Такого пользователя не существует"})
			return
		}
	}

	pass = alias.HashPass(pass)
	if pass == user.Pass {
		setID(c, user)
		c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Неправильный пароль!"})
	}
}
func testHandler(c *gin.Context) {
	// mysql := db.CreateMysqlUser()
	// ticketLogArr := mysql.GetUserTickets(2)
	// strRes, _ := json.Marshal(ticketLogArr)
	// strTicket := chk(json.Marshal(&ticket)).(string)
	// c.String(http.StatusOK, string(strRes))
}
