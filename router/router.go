package router

import (
	"CypherDesk-main/db"
	"log"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	// Protocol which used to access to the server
	Protocol = "http"
	// Host of the server
	Host = "127.0.0.1"
	// Port which used to access to the server
	Port = "3000"

	loginCharset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// New returns pointer on gin.Engine obj with settings
func New() *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/", indexHandler)

	//Account
	router.GET("/account", accountHandler)
	router.POST("/authorize", authorizeHandler)
	router.POST("/fillAdminAccount", fillAdminAccountHandler)
	router.GET("/out", logOutHandler)
	router.GET("/activate/:key", activateAccountHandler)
	router.POST("/fillUserAccount", fillUserAccountHandler)
	router.POST("/remindPass", remindPassHandler)
	router.GET("/remindPass/chk/:login/:key", checkChangeCredentialsKeyHandler)
	router.POST("/remindPass/change", changeCredentialsHandler)

	// AdminPanel
	router.POST("/createDepartment", createDepartmentHandler)
	router.POST("/createUser", createUserHandler)
	router.POST("/findUser", findUserHandler)
	router.POST("/changeUser", changeUserHandler)
	router.GET("/deleteUser", deleteUserHandler)

	router.GET("/test", testHandler)

	//	router.LoadHTMLGlob("templates/**/template.html")
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	return router
}

func chk(obj interface{}, err error) interface{} {
	if err != nil {
		log.Fatal("panic in handlers: " + err.Error())
		panic(err.Error())
	}
	return obj
}

func getPongoTemplate(filepath string, pc pongo2.Context) []byte {
	temp := pongo2.Must(pongo2.FromFile(filepath))
	out := chk(temp.ExecuteBytes(pc))
	return out.([]byte)
}

func writePongoTemplate(filepath string, pc pongo2.Context, c *gin.Context) []byte {
	out := getPongoTemplate(filepath, pc)

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(out)

	return []byte(out)
}

func rec(c *gin.Context) {
	if err := recover(); err != nil {
		log.Fatal(err)
		c.JSON(
			http.StatusOK,
			gin.H{
				"err": "Ошибка",
				"ok":  false,
			},
		)
	}
}

func getID(c *gin.Context) (bool, int) {
	session := sessions.Default(c)
	id := session.Get("id")
	if id == nil || id == 0 {
		return false, 0
	}
	return true, id.(int)
}

func setID(c *gin.Context, user *db.User) {
	session := sessions.Default(c)
	session.Set("id", user.ID)
	session.Save()
}
