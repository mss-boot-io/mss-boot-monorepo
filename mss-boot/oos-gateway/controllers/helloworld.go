/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/24 14:59:13
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/24 14:59:13
 */

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/pkg/response"

	"github.com/WhiteMatrixTech/matrix-cloud-monorepo/oos-gateway/form"
)

func init() {
	e := &Helloworld{}
	//e.Auth = false
	response.AppendController(e)
}

type Helloworld struct {
	response.Api
	//curd.DefaultController
}

func (Helloworld) Path() string {
	return "/oos-gateway"
}

func (e *Helloworld) Other(r *gin.RouterGroup) {
	//r.Use((&middlewares.AuthMiddleware{}).AuthMiddleware())
	r.POST("/call", e.Call)
}

// Call somebody
// @Summary Call somebody
// @Description Call somebody
// @Tags oos-gateway
// @Accept  application/json
// @Product application/json
// @Param data body form.HelloworldCallReq true "data"
// @Success 200 {object} response.Response{data=form.HelloworldCallResp}
// @Router /oos-gateway/api/v1/call [post]
func (e Helloworld) Call(c *gin.Context) {
	req := &form.HelloworldCallReq{}
	err := e.Make(c).Bind(req).Error
	if err != nil {
		e.Err(http.StatusUnprocessableEntity, err)
		return
	}
	result := &form.HelloworldCallResp{
		Message: "Hello " + req.Name,
	}
	e.OK(result)
}
