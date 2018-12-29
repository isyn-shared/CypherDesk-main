package router

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"CypherDesk-main/feedback"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func createUserHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Вы не авторизованы"})
		return
	}
	userMail, userRole, userDepartment := c.PostForm("mail"), c.PostForm("role"), c.PostForm("department")
	if alias.EmptyStrArr([]string{userMail, userRole, userDepartment}) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Укажите все данные!"})
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	if user.Role != "admin" || !user.Filled() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "У Вас нет прав на это действие"})
		return
	}
	encMail := alias.StandartRefact(userMail, false, db.StInfoKey)
	user = mysql.GetUser("mail", encMail)
	if user.Exist() {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Пользователь с такой почтой уже существует!"})
		return
	}
	uDid, _ := alias.STI(userDepartment)
	NewUser := &db.User{
		Mail:       userMail,
		Department: uDid,
		Role:       userRole,
	}
	NewUser.GenerateLogin(12)
	NewUser.GeneratePass(15)
	privPass := NewUser.Pass
	NewUser.HashPass()

	fmt.Println("NewUSerMail:", NewUser.Mail)

	mailText := "Для авторизации используйте логин: " + NewUser.Login + " и пароль: " + privPass + ". Приятного пользования!"
	r := feedback.NewMailRequest([]string{NewUser.Mail}, "Восстановление пароля CypherDesk")

	mysql.InsertUser(NewUser)
	r.Send("templates/mail/body.html", map[string]string{"text": mailText})

	c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
}

func findUserHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	mysql := db.CreateMysqlUser()
	userAdmin := mysql.GetUser("id", id)
	if userAdmin.Role != "admin" || !userAdmin.Filled() {
		c.JSON(http.StatusOK, "У Вас нет прав на это действие")
		return
	}

	keyPhrases := strings.Split(c.PostForm("key"), " ")
	if len(keyPhrases) == 0 || len(keyPhrases) > 5 {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Неправильные POST данные"})
		return
	}

	users := make([]*db.User, 0)
	defer func() {
		for _, u := range users {
			u.HidePrivateInfo()
		}
		strRes := string(chk(json.Marshal(users)).([]byte))
		c.String(http.StatusOK, strRes)
	}()

	if len(keyPhrases) == 1 {
		switch keyPhrases[0] {
		case "*":
			users = mysql.GetUsers("*", "*")
			return
		}
	}

	users = mysql.FindUser(keyPhrases)
}

func createDepartmentHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Вы не авторизованы!"})
		return
	}
	mysql := db.CreateMysqlUser()

	user := mysql.GetUser("id", id)
	if !user.Exist() || user.Role != "admin" {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "У Вас нет прав на это действие!!"})
		return
	}

	depName := c.PostForm("name")
	mysql.InsertDepartment(depName)
	c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
}

func deleteUserHandler(c *gin.Context) {
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Вы не авторизованы", "redirect": "/"})
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	if !user.Exist() || user.Role != "admin" {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "У Вас нет прав на это действие!!"})
		return
	}
	login := c.PostForm("login")
	if alias.EmptyStr(login) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Неправильный запрос"})
		return
	}
	login = alias.StandartRefact(login, false, db.StInfoKey)
	mysql.DeleteUser("login", login)
	c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
}

func changeUserHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Вы не авторизованы", "redirect": "/"})
		return
	}
	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	if !user.Exist() || user.Role != "admin" {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "У Вас нет прав на это действие!!"})
		return
	}
	login, newLogin, newName, newSurname, newPartonymic, newRecourse, newDepartment := c.PostForm("login"),
		c.PostForm("newLogin"), c.PostForm("newName"), c.PostForm("newSurname"), c.PostForm("newPartonymic"),
		c.PostForm("newRecourse"), c.PostForm("newDepartment")

	encLogin := alias.StandartRefact(login, false, db.StInfoKey)
	user = mysql.GetUser("login", encLogin)

	newEncLogin := alias.StandartRefact(newLogin, false, db.StInfoKey)
	tmpUser := mysql.GetUser("login", newEncLogin)
	if tmpUser.Exist() && tmpUser.ID != user.ID {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Пользователь с таким логином уже существует!"})
		return
	}

	user.Login, user.Name, user.Surname, user.Partonymic, user.Recourse, user.Department = newLogin, newName,
		newSurname, newPartonymic, newRecourse, chk(alias.STI(newDepartment)).(int)
	mysql.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
}

func changeDepartment(c *gin.Context) {
	prevName, newName := c.PostForm("prevName"), c.PostForm("name")
	if alias.EmptyStr(prevName) || alias.EmptyStr(newName) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Неправильный запрос"})
		return
	}
	mysql := db.CreateMysqlUser()
	mysql.UpdateDepartment(&db.Department{ID: 0, Name: prevName}, newName)
	c.JSON(http.StatusOK, gin.H{"ok": true, "err": nil})
}
