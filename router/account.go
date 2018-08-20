package router

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"CypherDesk-main/feedback"
	"net/http"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func accountHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	if !user.Filled() {
		if user.Role == "admin" {
			writePongoTemplate("templates/fillAccount/admin.html", pongo2.Context{}, c)
		} else {
			writePongoTemplate("templates/fillAccount/user.html", pongo2.Context{}, c)
		}
	} else {
		if user.Role == "admin" {
			departments := mysql.GetDepartments()
			writePongoTemplate("templates/adminPanel/index.html", pongo2.Context{
				"name":        user.Name,
				"surname":     user.Surname,
				"partonymic":  user.Partonymic,
				"recourse":    user.Recourse,
				"mail":        user.Mail,
				"login":       user.Login,
				"departments": departments,
			}, c)
		} else {
			department := user.GetDepartment()
			writePongoTemplate("templates/homePage/index.html", pongo2.Context{
				"name":       user.Name,
				"surname":    user.Surname,
				"partonymic": user.Partonymic,
				"recourse":   user.Recourse,
				"mail":       user.Mail,
				"login":      user.Login,
				"department": department.Name,
			}, c)
		}
	}
}

func logOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("id")
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}

func fillAdminAccountHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	updUser := new(db.User)
	updUser.Name, updUser.Surname, updUser.Partonymic = c.PostForm("name"), c.PostForm("surname"), c.PostForm("partonymic")
	updUser.Recourse = c.PostForm("recourse")
	updUser.Mail, updUser.Login, updUser.Pass = c.PostForm("mail"), c.PostForm("login"), c.PostForm("pass")
	if alias.EmptyStrArr([]string{updUser.Name, updUser.Surname, updUser.Partonymic, updUser.Recourse,
		updUser.Mail, updUser.Login, updUser.Pass}) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Заполните все поля!"})
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	user.WriteIn(updUser)

	if !user.ChkLogin() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Логин не подходит по регулярке"})
	} else if !user.ChkPass() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Пароль не подходит по регулярке"})
	} else if user.Login != updUser.Login && mysql.GetUser("login", updUser.Login).Exist() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Такой логин уже существует!"})
	} else {
		user.HashPass()
		activationKey := user.Mail + time.Now().String()
		user.SetActivationKey(activationKey)
		mysql.UpdateUser(user)

		mailMsg := &feedback.MailMessage{
			Subject:    "Активация аккаунта",
			Body:       "Здравствуйте! Спасибо за использование системы CypherDesk. Для активации акаунта перейдите по ссылке: " + Protocol + "://" + Host + ":" + Port + "/activate/" + user.ActivationKey,
			Recipients: []string{user.Mail},
		}

		feedback.SendMail(mailMsg)

		c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
	}
}

func fillUserAccountHandler(c *gin.Context) {

}

func activateAccountHandler(c *gin.Context) {
	key := c.Param("key")
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.String(http.StatusOK, "Вы не авторизованы!")
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	if user.ActivationKey == key {
		c.String(http.StatusOK, "Yeeeaaaahhh!!!")
	} else {
		c.String(http.StatusOK, "Ohh, nooo!!!")
	}
}
