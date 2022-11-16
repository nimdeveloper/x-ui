package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"x-ui/database/model"
	requestBody "x-ui/web/request_body"
	"x-ui/web/service/mirror"
)

type MirrorStoreController struct {
	storeMirrorService mirror.StoreService
}

func NewStoreMirrorController(g *gin.RouterGroup) *MirrorStoreController {
	a := &MirrorStoreController{}
	a.initRouter(g)
	return a
}

func (a *MirrorStoreController) initRouter(g *gin.RouterGroup) {
	g.POST("/api/mirrors", a.storeMirror)
}

func (a *MirrorStoreController) storeMirror(c *gin.Context) {
	var body requestBody.StoreMirrorRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	m := model.Mirror{
		Ip:   body.Ip,
		Port: body.Port,
	}
	err := a.storeMirrorService.SaveMirror(&m)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   m,
	})
}
