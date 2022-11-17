package controller

import (
	"github.com/gin-gonic/gin"
)

type XUIController struct {
	BaseController

	settingController     *SettingController
	mirrorIndexController *MirrorIndexController
	mirrorStoreController *MirrorStoreController
}

func NewXUIController(g *gin.RouterGroup) *XUIController {
	a := &XUIController{}
	a.initRouter(g)
	return a
}

func (a *XUIController) initRouter(g *gin.RouterGroup) {
	a.mirrorIndexController = NewMirrorIndexController(g)
	a.mirrorStoreController = NewStoreMirrorController(g)

	g = g.Group("/xui")
	g.Use(a.checkLogin)

	g.GET("/", a.index)
	g.GET("/inbounds", a.inbounds)
	g.GET("/setting", a.setting)

	a.settingController = NewSettingController(g)
}

func (a *XUIController) index(c *gin.Context) {
	html(c, "index.html", "pages.index.title", nil)
}

func (a *XUIController) inbounds(c *gin.Context) {
	html(c, "inbounds.html", "pages.inbounds.title", nil)
}

func (a *XUIController) setting(c *gin.Context) {
	html(c, "setting.html", "pages.setting.title", nil)
}
