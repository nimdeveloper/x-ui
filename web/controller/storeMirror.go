package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"x-ui/database/model"
	requestBody2 "x-ui/web/requestBody"
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
	g.POST("/mirrors", a.storeMirror)
}

func (a *MirrorStoreController) storeMirror(c *gin.Context) {
	var requestBody requestBody2.StoreMirrorRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	m := model.Mirror{
		Ip:   requestBody.Ip,
		Port: requestBody.Port,
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
