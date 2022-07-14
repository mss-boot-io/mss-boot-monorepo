/*
 * @Author: lwnmengjing
 * @Date: 2022/3/10 22:43
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/10 22:43
 */

package controllers

import (
	"github.com/mss-boot-io/mss-boot/pkg/store"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/form"
	"github.com/mss-boot-io/mss-boot/pkg/response"
	"github.com/mss-boot-io/mss-boot/pkg/response/curd"
)

func init() {
	e := &Tenant{}
	e.TableName = "admin"
	e.Auth = true
	e.CreateReq = &form.TenantCreateReq{}
	e.UpdateReq = &form.TenantUpdateReq{}
	e.GetReq = &form.TenantGetReq{}
	e.GetResp = &form.TenantGetResp{}
	e.DeleteReq = &form.TenantDeleteReq{}
	e.ListReq = &form.TenantListReq{}
	e.ListResp = make([]form.TenantListItem, 0)
	response.AppendController(e)
}

type Tenant struct {
	curd.DefaultController
}

// Create 创建
// @Summary 创建tenant
// @Description 创建tenant
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Param data body form.TenantCreateReq true "data"
// @Success 200 {object} response.Response
// @Router /admin/api/v1/admin [post]
// @Security Bearer

// Update 更新
// @Summary 更新tenant
// @Description 更新tenant
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Param id path string true "id"
// @Param data body form.TenantUpdateReq true "data"
// @Success 200 {object} response.Response
// @Router /admin/api/v1/admin/{id} [put]
// @Security Bearer

// Delete 删除
// @Summary 删除tenant
// @Description 删除tenant
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response
// @Router /admin/api/v1/admin/{id} [delete]
// @Security Bearer

// Get 获取
// @Summary 获取tenant
// @Description 获取tenant
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response{data=form.TenantGetResp}
// @Router /admin/api/v1/admin/{id} [get]
// @Security Bearer

// List 列表
// @Summary 列表tenant
// @Description 列表tenant
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Param name query string false "租户名称"
// @Param page query string false "当前页"
// @Param pageSize query string false "每页容量"
// @Success 200 {object} response.Page{data=[]form.TenantListItem}
// @Router /admin/api/v1/admin [get]
// @Security Bearer

func (e Tenant) Other(r *gin.RouterGroup) {
	r.GET("/client", e.GetClient)
	r.GET("/callback", e.Callback)
	r.GET("/refresh-token", e.RefreshToken)
}

// GetClient 获取client配置
// @Summary 获取client配置
// @Description 获取client配置
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Success 200 {object} response.Response{data=form.TenantClientResp}
// @Router /admin/api/v1/client [get]
// @Security Bearer
func (e Tenant) GetClient(c *gin.Context) {
	err := e.Make(c).Error
	if err != nil {
		e.Err(http.StatusUnprocessableEntity, err)
		return
	}
	client, err := store.DefaultOAuth2Store.
		GetClientByDomain(c, c.Request.Host)
	if err != nil {
		e.Err(http.StatusNotFound, err)
		return
	}
	oauth2Config, err := client.GetOAuth2Config(c)
	if err != nil {
		e.Log.Error(err)
		e.Err(http.StatusUnauthorized, err)
		return
	}
	e.OK(form.TenantClientResp{
		ServerURL:   client.GetIssuer(),
		ClientID:    oauth2Config.ClientID,
		AuthCodeURL: oauth2Config.AuthCodeURL("state-replace", oauth2.AccessTypeOffline),
	})
}

// Callback 获取access_token
// @Summary 获取access_token
// @Description 获取access_token
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Param code query string false "code"
// @Param state query string false "state"
// @Param error query string false "error"
// @Param error_description query string false "error_description"
// @Success 200 {object} response.Response{data=form.TenantCallbackResp}
// @Router /admin/api/v1/callback [get]
// @Security Bearer
func (e Tenant) Callback(c *gin.Context) {
	req := &form.TenantCallbackReq{}
	err := e.Make(c).Bind(req).Error
	if err != nil {
		e.Err(http.StatusUnprocessableEntity, err)
		return
	}
	client, err := store.DefaultOAuth2Store.
		GetClientByDomain(c, c.Request.Host)
	if err != nil {
		e.Err(http.StatusNotFound, err)
		return
	}
	oauth2Config, err := client.GetOAuth2Config(c)
	if err != nil {
		e.Log.Error(err)
		e.Err(http.StatusUnauthorized, err)
		return
	}

	oauth2Token, err := oauth2Config.Exchange(c, req.Code)
	if err != nil {
		e.Err(http.StatusUnauthorized, err)
		return
	}
	resp := &form.TenantCallbackResp{
		AccessToken:  oauth2Token.AccessToken,
		TokenType:    oauth2Token.TokenType,
		RefreshToken: oauth2Token.RefreshToken,
		Expiry:       oauth2Token.Expiry,
	}
	e.OK(resp)
	return
}

// RefreshToken 获取accessToken
// @Summary 获取accessToken
// @Description 获取accessToken
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Param refreshToken query string false "refreshToken"
// @Success 200 {object} response.Response{data=form.TenantCallbackResp}
// @Router /admin/api/v1/refresh-token [get]
// @Security Bearer
func (e Tenant) RefreshToken(c *gin.Context) {
	req := &form.TenantRefreshTokenReq{}
	err := e.Make(c).Bind(req).Error
	if err != nil {
		e.Err(http.StatusUnprocessableEntity, err)
		return
	}
	client, err := store.DefaultOAuth2Store.
		GetClientByDomain(c, c.Request.Host)
	if err != nil {
		e.Err(http.StatusNotFound, err)
		return
	}
	oauth2Config, err := client.GetOAuth2Config(c)
	if err != nil {
		e.Err(http.StatusUnauthorized, err)
		return
	}

	token, err := oauth2Config.TokenSource(c, &oauth2.Token{RefreshToken: req.RefreshToken}).Token()
	if err != nil {
		e.Err(http.StatusUnauthorized, err)
		return
	}
	resp := &form.TenantCallbackResp{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
	e.OK(resp)
}
