package router

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"CypherDesk-main/feedback"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createUserHandler(c *gin.Context) {
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
	if user.Role != "admin" {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "У Вас нет прав на это действие"})
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
	mysql.InsertUser(NewUser)

	mailMsg := &feedback.MailMessage{
		Subject:    "Регистрация CypherDesk",
		Body:       "Здравствуйте. Для авторизации используйте логин: " + NewUser.Login + " и пароль: " + NewUser.Pass + ". Приятного пользования!",
		Recipients: []string{user.Mail},
	}

	feedback.SendMail(mailMsg)
}
