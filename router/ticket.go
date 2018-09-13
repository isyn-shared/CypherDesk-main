package router

import (
	"CypherDesk-main/alias"
	"CypherDesk-main/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func sendTicketHandler(c *gin.Context) {
	defer rec(c)
	mysql := db.CreateMysqlUser()
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Вы не авторизованы"})
		return
	}
	user := mysql.GetUser("id", id)
	if !user.Exist() || !user.Filled() { //TODO: модератор тикетов - как назвать?
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "У Вас нет прав"})
		return
	}
	caption, description := c.PostForm("caption"), c.PostForm("description")
	if alias.EmptyStr(caption) || alias.EmptyStr("description") {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Неправильный запрос"})
		return
	}
	ticket := &db.Ticket{
		Caption:     caption,
		Description: description,
		Sender:      id,
		Status:      "opened",
	}
	mysql.CreateTicket(ticket)
	ticketAdmin := mysql.GetDepartmentTicketAdmin(user.Department)
	log := &db.TicketLog{
		Ticket:   mysql.GetLastLogId(),
		UserFrom: id,
		UserTo:   ticketAdmin.ID,
		Action:   "send",
		Time:     time.Now(),
	}
	mysql.TransferTicket(log)
}

func forwardTicketHandler(c *gin.Context) {
	defer rec(c)
	mysql := db.CreateMysqlUser()
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Вы не авторизованы"})
		return
	}
	user := mysql.GetUser("id", id)
	if !user.Exist() || !user.Filled() || user.Role != "" { //TODO: модератор тикетов - как назвать?
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "У Вас нет прав!"})
		return
	}
	to, ticketID := c.PostForm("to"), c.PostForm("ticketID")
	if alias.EmptyStr(to) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "err": "Неправильный запрос"})
		return
	}
	log := &db.TicketLog{
		Ticket:   chk(alias.STI(ticketID)).(int),
		UserFrom: id,
		UserTo:   chk(alias.STI(to)).(int),
		Action:   "forward",
		Time:     time.Now(),
	}
	mysql.TransferTicket(log)
}
