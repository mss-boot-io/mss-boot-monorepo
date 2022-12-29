/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/19 15:28:20
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/19 15:28:20
 */

package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/pkg/middlewares"
	"github.com/mss-boot-io/mss-boot/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"

	tenant "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/models"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/form"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/models"
)

func init() {
	e := &Github{}
	response.AppendController(e)
}

// Github github
type Github struct {
	response.Api
}

func (Github) Path() string {
	return "github"
}

func (e Github) Other(r *gin.RouterGroup) {
	r.Use((&middlewares.AuthMiddleware{}).AuthMiddleware())
	r.POST("/github/create-or-update", e.CreatOrUpdate)
	r.GET("/github/get", e.Get)
}

// CreatOrUpdate 创建或更新github配置
// @Summary 创建或更新github配置
// @Description 创建或更新github配置
// @Tags generator
// @Accept  application/json
// @Product application/json
// @Param data body form.GithubCreateReq true "data"
// @Success 200 {object} response.Response
// @Router /generator/api/v1/github/create-or-update [post]
// @Security Bearer
func (e Github) CreatOrUpdate(c *gin.Context) {
	req := &form.GithubCreateReq{}
	err := e.Make(c).Bind(req).Error
	if err != nil {
		e.Err(http.StatusUnprocessableEntity, err)
		return
	}
	user := middlewares.GetLoginUser(c)
	if user == nil {
		e.Err(http.StatusUnauthorized, errors.New("user is empty"))
		return
	}

	t := tenant.Tenant{}
	err = t.GetTenantByDomain(c, c.Request.Host)
	if err != nil {
		e.Err(http.StatusUnauthorized, err)
		return
	}

	now := time.Now()
	g := &models.Github{
		ID:        primitive.NewObjectID().Hex(),
		TenantID:  t.ID.Hex(),
		Email:     user.Email,
		Password:  req.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err = g.C().InsertOne(c, g)
	if err != nil {
		e.Err(http.StatusInternalServerError, err)
		return
	}
	e.OK(nil)
}

// Get 获取github配置
// @Summary 获取github配置
// @Description 获取github配置
// @Tags generator
// @Accept  application/json
// @Product application/json
// @Success 200 {object} response.Response{data=form.GithubGetResp}
// @Router /generator/api/v1/github/get [get]
// @Security Bearer
func (e Github) Get(c *gin.Context) {
	e.Make(c)
	user := middlewares.GetLoginUser(c)
	if user == nil {
		e.Err(http.StatusUnauthorized, errors.New("user is empty"))
		return
	}

	t := tenant.Tenant{}
	err := t.GetTenantByDomain(c, c.Request.Host)
	if err != nil {
		e.Err(http.StatusUnauthorized, err)
		return
	}
	g, err1 := models.GetMyGithubConfig(c, t.ID.Hex(), user.Email)
	result := &form.GithubGetResp{
		Email:     user.Email,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if err1 == nil {
		result.Configured = true
	}
	e.OK(result)
}
