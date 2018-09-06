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
		if user.ActivationKey != "" && user.ActivationType == 1 {
			writePongoTemplate("templates/fillAccount/activationMessage.html", pongo2.Context{}, c)
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
		c.JSON(http.StatusSeeOther, gin.H{"ok": false, "err": "Вы не авторизованы", "redirect": "/"})
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	if user.Filled() {
		c.JSON(http.StatusSeeOther, gin.H{"ok": false, "err": "Вы уже заполнили аккаунт!", "redirect": "/account"})
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
	user.WriteIn(updUser)

	if !user.ChkLogin() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Логин не подходит по регулярке"})
	} else if !user.ChkPass() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Пароль не подходит по регулярке"})
	} else if user.Login != updUser.Login && mysql.GetUser("login", updUser.Login).Exist() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Такой логин уже существует!"})
	} else {
		user.HashPass()
		activationKey := alias.StringWithCharset(20, loginCharset) + user.Mail + time.Now().String()
		user.SetActivationKey(activationKey)
		mysql.UpdateUser(user)

		mailMsg := &feedback.MailMessage{
			Subject:    "Активация аккаунта",
			Body:       "Здравствуйте! Спасибо за использование системы CypherDesk. Для активации акаунта перейдите по ссылке: " + Protocol + "://" + Host + ":" + Port + "/activate/" + user.ActivationKey,
			Recipients: []string{user.Mail},
		}

		feedback.SendMail(mailMsg)

		c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil, "redirect": "/account"})
	}
}

func fillUserAccountHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	updUser := new(db.User)
	updUser.Name, updUser.Surname, updUser.Partonymic = c.PostForm("name"), c.PostForm("surname"), c.PostForm("partonymic")
	updUser.Recourse = c.PostForm("recourse")
	updUser.Login, updUser.Pass = c.PostForm("login"), c.PostForm("pass")
	if alias.EmptyStrArr([]string{updUser.Name, updUser.Surname, updUser.Partonymic, updUser.Recourse,
		updUser.Login, updUser.Pass}) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Заполните все поля!"})
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	prevLogin := user.Login
	user.WriteIn(updUser)

	if !user.ChkLogin() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Логин не подходит по регулярке"})
	} else if !user.ChkPass() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Пароль не подходит по регулярке"})
	} else if user.Login != prevLogin && mysql.GetUser("login", updUser.Login).Exist() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Такой логин уже существует!"})
	} else {
		user.HashPass()
		mysql.UpdateUser(user)
		c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
	}
}

func activateAccountHandler(c *gin.Context) {
	key := c.Param("key")
	isAuthorized, id := getID(c)
	if !isAuthorized {
		writePongoTemplate("templates/fillAccount/failActivation.html", pongo2.Context{
			"err": "Вы не авторизованы!",
		}, c)
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	if user.ActivationKey == key && user.ActivationType == 1 {
		user.ActivationKey = ""
		user.ActivationType = 0
		mysql.UpdateUser(user)
		c.Redirect(http.StatusSeeOther, "/account")
	} else {
		writePongoTemplate("templates/fillAccount/failActivation.html", pongo2.Context{
			"err": "Активационный ключ не действителен!",
		}, c)
	}
}

func remindPassHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, _ := getID(c)
	if isAuthorized {
		c.Redirect(http.StatusSeeOther, "/account")
		return
	}

	credentials := c.PostForm("credentials")
	if alias.EmptyStr(credentials) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Неправильный запрос"})
		return
	}

	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("login", credentials)
	if !user.Exist() {
		user = mysql.GetUser("mail", credentials)
		if !user.Exist() {
			c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Такого пользователя не существует!"})
			return
		}
	}

	user.SetRemindKey(alias.StringWithCharset(20, loginCharset) + time.Now().String())

	mailMsg := &feedback.MailMessage{
		Subject: "Восстановление пароля CypherDesk",
		Body: "Ваш логин: " + user.Login + ". Для восстановления пароля перейдите по ссылке: " + Protocol + "://" + Host + ":" +
			Port + "/remindPass/chk/" + user.Login + "/" + user.ActivationKey,
		Recipients: []string{user.Mail},
	}

	mysql.UpdateUser(user)
	feedback.SendMail(mailMsg)
	c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
}

func checkChangeCredentialsKeyHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, _ := getID(c)
	if isAuthorized {
		c.Redirect(http.StatusSeeOther, "/account")
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("login", c.Param("login"))
	if user.ActivationType != 2 || user.ActivationKey != c.Param("key") {
		writePongoTemplate("templates/fillAccount/failActivation.html", pongo2.Context{}, c)
		return
	}
	session := sessions.Default(c)
	session.Set("updatePass", user.ActivationKey)
	session.Save()
	writePongoTemplate("templates/fillAccount/changePass.html", pongo2.Context{
		"login":      user.Login,
		"mail":       user.Mail,
		"name":       user.Name,
		"surname":    user.Surname,
		"partonymic": user.Partonymic,
	}, c)
}

func changeCredentialsHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, _ := getID(c)
	if isAuthorized {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Вы уже авторизованы!", "redirect": "/account"})
		return
	}

	userLogin, userNewPass := c.PostForm("login"), c.PostForm("pass")
	if alias.EmptyStr(userLogin) || alias.EmptyStr(userNewPass) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Некорректный запрос!"})
		return
	}

	session := sessions.Default(c)
	updatePassKey := session.Get("updatePass")

	if alias.EmptyStr(updatePassKey.(string)) {
		writePongoTemplate("templates/fillAccount/failAccount.html", pongo2.Context{}, c)
		return
	}

	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("login", userLogin)

	if user.ActivationKey != updatePassKey {
		writePongoTemplate("templates/fillAccount/failAccount.html", pongo2.Context{}, c)
		return
	}

	user.Pass = userNewPass
	user.ActivationKey = ""
	user.ActivationType = 0
	user.HashPass()
	mysql.UpdateUser(user)

	session.Delete("id")
	session.Save()

	c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil, "redirect": "/"})
}
