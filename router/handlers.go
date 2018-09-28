package router

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	defer rec(c)
	writePongoTemplate("templates/front-page/index.html", pongo2.Context{}, c)
}
