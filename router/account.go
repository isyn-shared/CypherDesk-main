package router

import (
	"CypherDesk-main/db"
	"net/http"

	"github.com/flosch/pongo2"
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
		writePongoTemplate("templates/fillAccount/index.html", pongo2.Context{}, c)
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
