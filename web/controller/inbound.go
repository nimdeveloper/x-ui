package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	_ "github.com/samber/lo"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/web/global"
	"x-ui/web/response"
	"x-ui/web/service"
	"x-ui/web/session"
)

type InboundController struct {
	inboundService service.InboundService
	xrayService    service.XrayService
}

func NewInboundController(g *gin.RouterGroup) *InboundController {
	a := &InboundController{}
	a.initRouter(g)
	a.startTask()
	return a
}

func (a *InboundController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/inbounds")

	g.GET(":id", a.getInbound)
	g.GET("", a.getInbounds)
	g.POST("", a.addInbound)
	g.DELETE(":id", a.deleteInbound)
	g.PUT(":id", a.updateInbound)

	g.POST("/clientIps/:email", a.getClientIps)
	g.POST("/clearClientIps/:email", a.clearClientIps)

}

func (a *InboundController) startTask() {
	webServer := global.GetWebServer()
	c := webServer.GetCron()
	c.AddFunc("@every 10s", func() {
		if a.xrayService.IsNeedRestartAndSetFalse() {
			err := a.xrayService.RestartXray(false)
			if err != nil {
				logger.Error("restart xray failed:", err)
			}
		}
	})
}

func (a *InboundController) getInbounds(c *gin.Context) {
	user := session.GetLoginUser(c)
	inbounds, err := a.inboundService.GetInbounds(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
	}

	c.JSON(http.StatusOK, response.ListResponse[response.InboundResponse]{
		Data: lo.Map[*model.Inbound, *response.InboundResponse](inbounds, func(item *model.Inbound, _ int) *response.InboundResponse {
			return response.InboundResponseFromInbound(*item)
		}),
	})
}

func (a *InboundController) getInbound(c *gin.Context) {
	user := session.GetLoginUser(c)
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	inbound, err := a.inboundService.GetUserInbound(user.Id, id)

	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorMessage: "Inbound not found",
		})
	}

	c.JSON(http.StatusOK, response.InboundResponseFromInbound(*inbound))
}

func (a *InboundController) addInbound(c *gin.Context) {
	inbound := &model.Inbound{}
	err := c.ShouldBind(inbound)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	user := session.GetLoginUser(c)
	inbound.UserId = user.Id
	inbound.Enable = true
	inbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)
	inbound, err = a.inboundService.AddInbound(inbound)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.InboundResponseFromInbound(*inbound))
	a.xrayService.SetToNeedRestart()
}

func (a *InboundController) deleteInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	err = a.inboundService.DelInbound(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		SuccessMessage: I18n(c, "api.inbound.removeInboundSuccessMessage"),
	})
	a.xrayService.SetToNeedRestart()
}

func (a *InboundController) updateInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	inbound := &model.Inbound{
		Id: id,
	}
	err = c.ShouldBind(inbound)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	inbound, err = a.inboundService.UpdateInbound(inbound)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, response.ErrorResponse{
				ErrorMessage: err.Error(),
			})
			return
		}

		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.InboundResponseFromInbound(*inbound))
	a.xrayService.SetToNeedRestart()
}

func (a *InboundController) getClientIps(c *gin.Context) {
	email := c.Param("email")

	ips, err := a.inboundService.GetInboundClientIps(email)
	if err != nil {
		jsonObj(c, "No IP Record", nil)
		return
	}
	jsonObj(c, ips, nil)
}

func (a *InboundController) clearClientIps(c *gin.Context) {
	email := c.Param("email")

	err := a.inboundService.ClearClientIps(email)
	if err != nil {
		jsonMsg(c, "修改", err)
		return
	}
	jsonMsg(c, "Log Cleared", nil)
}
