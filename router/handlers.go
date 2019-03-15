package router

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"fmt"
	"net/http"
	"time"

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

	login = alias.StandartRefact(login, false, db.StInfoKey)

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

func createTemporaryHandler(c *gin.Context) {
	m := db.CreateMysqlUser()
	u := &db.User{
		Login:      "admin",
		Pass:       "admin",
		Role:       "admin",
		Department: 4,
	}
	u.HashPass()

	m.InsertUser(u)

	u = m.GetUserByDecField("login", "admin")

	eu := &db.ExtendedUser{
		ID:                  u.ID,
		Phone:               "+79782568334",
		Address:             "Simferopol",
		ActivationType:      "0",
		AccountCreationDate: time.Now(),
		ActivationDate:      time.Now(),
	}
	m.InsertExtendedUser(eu)
	c.String(http.StatusOK, "OK")
}

func testHandler(c *gin.Context) {
	isAuth, id := getID(c)
	if !isAuth {
		c.String(http.StatusOK, "You are not authorized")
		return
	}
	chatMsg := &db.ChatMsg{
		From:   id,
		To:     id,
		Text:   "KEKLol",
		Date:   time.Now(),
		Status: 1,
	}
	mysql := db.CreateMysqlUser()
	mysql.InsertChatMessage(chatMsg)

	chatMessages := mysql.GetUsersChatMessages(&db.User{ID: id})
	fmt.Println(chatMessages)

	for _, chm := range chatMessages {
		fmt.Println(chm.Text, chm.Date)
	}
}
