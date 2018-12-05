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
			writePongoTemplate("templates/fillAccount/admin.html", pongo2.Context{
				"mail": user.Mail,
			}, c)
		} else {
			writePongoTemplate("templates/fillAccount/user.html", pongo2.Context{}, c)
		}
	} else {
		if user.ActivationKey != "" && user.ActivationType == 1 {
			writePongoTemplate("templates/fillAccount/activationMessage.html", pongo2.Context{}, c)
		} else {
			if user.Role == "admin" {
				departments := mysql.GetDepartments()
				department := user.GetDepartment()
				writePongoTemplate("templates/adminPanel/index.html", pongo2.Context{
					"name":        user.Name,
					"surname":     user.Surname,
					"partonymic":  user.Partonymic,
					"recourse":    user.Recourse,
					"mail":        user.Mail,
					"login":       user.Login,
					"department":  department.Name,
					"departments": departments,
				}, c)
			} else if user.Role == "ticketModerator" {
				department := user.GetDepartment()
				departments := mysql.GetDepartments()
				usersInDep := mysql.GetDepartmentUsers("id", department.ID)
				admins := mysql.GetUsers("role", "admin")
				moderators := mysql.GetUsers("role", "ticketModerator")
				usersToTransfer := append(usersInDep, admins...)
				usersToTransfer = append(usersToTransfer, moderators...)

				for _, u := range usersToTransfer {
					u.HidePrivateInfo()
				}
				writePongoTemplate("templates/homePage/index.html", pongo2.Context{
					"isModerator":     true,
					"id":              user.ID,
					"name":            user.Name,
					"surname":         user.Surname,
					"partonymic":      user.Partonymic,
					"recourse":        user.Recourse,
					"mail":            user.Mail,
					"login":           user.Login,
					"department":      department.Name,
					"departments":     departments,
					"usersToTransfer": usersToTransfer,
				}, c)
			} else {
				department := user.GetDepartment()
				writePongoTemplate("templates/homePage/index.html", pongo2.Context{
					"isModerator": false,
					"id":          user.ID,
					"name":        user.Name,
					"surname":     user.Surname,
					"partonymic":  user.Partonymic,
					"recourse":    user.Recourse,
					"mail":        user.Mail,
					"login":       user.Login,
					"department":  department.Name,
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

	copyMail := updUser.Mail

	var newMail bool
	if alias.StandartRefact(user.Mail, true, db.StInfoKey) != updUser.Mail {
		newMail = true
	}

	if alias.EmptyStrArr([]string{updUser.Name, updUser.Surname, updUser.Partonymic, updUser.Recourse,
		updUser.Mail, updUser.Login, updUser.Pass}) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Заполните все поля!"})
		return
	}

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

		var mailMsg string
		if newMail {
			activationKey := alias.StringWithCharset(20, loginCharset) + user.Mail + time.Now().String()
			user.SetActivationKey(activationKey)

			mailMsg = "Здравствуйте! Спасибо за использование системы CypherDesk. Для активации акаунта перейдите по ссылке: " + Protocol + "://" + Host + ":" + Port + "/activate/" + user.ActivationKey
		} else {
			mailMsg = "Здравствуйте! Все прошло успешно! Спасибо за использование системы CypherDesk."
		}

		mysql.UpdateUser(user)
		r := feedback.NewMailRequest([]string{copyMail}, "Восстановление пароля CypherDesk")
		r.Send("templates/mail/body.html", map[string]string{"text": mailMsg})
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
	mysql.UpdateUser(user)

	mailMsg := "Ваш логин: " + user.Login + ". Для восстановления пароля перейдите по ссылке: " + Protocol + "://" + Host + ":" +
		Port + "/remindPass/chk/" + user.Login + "/" + user.ActivationKey

	r := feedback.NewMailRequest([]string{user.Mail}, "Восстановление пароля CypherDesk")
	r.Send("templates/mail/body.html", map[string]string{"text": mailMsg})

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
