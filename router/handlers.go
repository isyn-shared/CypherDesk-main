package router

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"CypherDesk-main/feedback"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	defer rec(c)
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
	setID(c, user)
	c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
}

func testHandler(c *gin.Context) {
	msg := new(feedback.MailMessage)
	msg.Recipients = []string{"nikita.surnachev03@gmail.com"}
	msg.Subject = "KekLOLOLOLOL"
	msg.Body = "Heeeey maaaaan"
	feedback.SendMail(msg)
}
