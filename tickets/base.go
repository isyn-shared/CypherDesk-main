package tickets

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func chk(v interface{}, err error) interface{} {
	if err != nil {
		fmt.Println("handle connection error: " + err.Error())
	}
	return v
}

func getID(c *gin.Context) (bool, int) {
	session := sessions.Default(c)
	id := session.Get("id")
	if id == nil || id == 0 {
		return false, 0
	}
	return true, id.(int)
}
