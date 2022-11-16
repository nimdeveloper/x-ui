package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"x-ui/web/service/mirror"
)

type MirrorIndexController struct {
	mirrorIndexService mirror.IndexService
}

func NewMirrorIndexController(g *gin.RouterGroup) *MirrorIndexController {
	a := &MirrorIndexController{}
	a.initRouter(g)
	return a
}

func (a *MirrorIndexController) initRouter(g *gin.RouterGroup) {
	g.GET("/api/mirrors", a.getAllMirrors)
}

func (a *MirrorIndexController) getAllMirrors(c *gin.Context) {
	mirrors, err := a.mirrorIndexService.GetMirrors()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   mirrors,
	})
}
