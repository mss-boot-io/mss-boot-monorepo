/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/12/29 16:01:32
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/12/29 16:01:32
 */

package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/pkg/response"
)

func init() {
	e := new(Casbin)
	response.AppendController(e)
}

type Casbin struct {
	response.Api
}

func (e Casbin) Path() string {
	return "/casbin"
}

func (e Casbin) Put(c *gin.Context) {
	e.Make(c)
	e.Log.Info("Casbin Put")
}

// Get 获取
// @Summary 获取casbin
// @Description 获取casbin
// @Tags casbin
// @Accept  application/json
// @Product application/json
// @Success 200 {object} response.Response{data=form.MenuGetResp}
// @Router /admin/api/v1/menu/{id} [get]
// @Security Bearer
func (e Casbin) Get(c *gin.Context) {
	e.Make(c)
	e.Log.Info("Casbin Get")
}

func (e Casbin) Other(engine *gin.RouterGroup) {
	//engine.GET("/casbin", e.Get)
	//engine.PUT("/casbin", e.Put)
}
