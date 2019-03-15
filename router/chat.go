package router

import (
	"CypherDesk-main/db"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func chatHandler(c *gin.Context) {
	defer rec(c)
	isAuthorized, id := getID(c)
	if !isAuthorized {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	mysql := db.CreateMysqlUser()
	user := mysql.GetUser("id", id)
	extUser := mysql.GetExtenedUser(user)

	if user.Role == "admin" {
		departments := mysql.GetDepartments()
		department := user.GetDepartment()
		users := mysql.GetUsersByDecField("role", "user")
		admins := mysql.GetUsersByDecField("role", "admin")
		tms := mysql.GetUsersByDecField("role", "ticketModerator")

		users = append(users, admins...)
		users = append(users, tms...)

		writePongoTemplate("templates/chat/index.html", pongo2.Context{
			"name":            user.Name,
			"surname":         user.Surname,
			"partonymic":      user.Partonymic,
			"recourse":        user.Recourse,
			"mail":            user.Mail,
			"login":           user.Login,
			"id":              user.ID,
			"department":      department.Name,
			"departments":     departments,
			"phone":           extUser.Phone,
			"Address":         extUser.Address,
			"usersToTransfer": users,
		}, c)
	} else if user.Role == "ticketModerator" {
		department := user.GetDepartment()
		departments := mysql.GetDepartments()
		usersInDep := mysql.GetDepartmentUsers("id", department.ID)
		admins := mysql.GetUsersByDecField("role", "admin")
		moderators := mysql.GetUsersByDecField("role", "ticketModerator")
		usersToTransfer := append(usersInDep, admins...)
		usersToTransfer = append(usersToTransfer, moderators...)

		for _, u := range usersToTransfer {
			u.HidePrivateInfo()
		}

		writePongoTemplate("templates/chat/index.html", pongo2.Context{
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
			"phone":           extUser.Phone,
			"Address":         extUser.Address,
		}, c)
	} else {
		department := user.GetDepartment()
		usersInDep := mysql.GetDepartmentUsers("id", department.ID)
		admins := mysql.GetUsersByDecField("role", "admin")
		usersToTransfer := append(usersInDep, admins...)
		writePongoTemplate("templates/chat/index.html", pongo2.Context{
			"isModerator":     false,
			"id":              user.ID,
			"name":            user.Name,
			"surname":         user.Surname,
			"partonymic":      user.Partonymic,
			"recourse":        user.Recourse,
			"mail":            user.Mail,
			"login":           user.Login,
			"department":      department.Name,
			"phone":           extUser.Phone,
			"Address":         extUser.Address,
			"usersToTransfer": usersToTransfer,
		}, c)
	}
}
