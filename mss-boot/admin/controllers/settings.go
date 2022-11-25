/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/19 03:10:31
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/19 03:10:31
 */

package controllers

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/pkg/middlewares"
	"github.com/mss-boot-io/mss-boot/pkg/response"
	"github.com/mss-boot-io/mss-boot/pkg/response/curd"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/form"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/models"
)

func init() {
	e := new(Settings)
	e.TableName = "settings"
	e.Auth = true
	response.AppendController(e)
}

type Settings struct {
	curd.DefaultController
}

// Get 获取
// @Summary 获取settings
// @Description 获取settings
// @Tags settings
// @Accept  application/json
// @Product application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response{data=[]form.KV}
// @Router /admin/api/v1/settings [get]
// @Security Bearer
func (e Settings) Get(c *gin.Context) {
	err := e.Make(c).Error
	if err != nil {
		e.Err(http.StatusUnprocessableEntity, err)
		return
	}
	u := middlewares.GetLoginUser(c)
	if u == nil {
		e.Err(http.StatusUnauthorized, errors.New("user is empty"))
		return
	}
	user, err := models.GetUserByNameAndDomain(c, c.Request.Host, u.Name)
	if err != nil {
		e.Err(http.StatusUnauthorized, err)
		return
	}
	settings := &models.Settings{}
	err = settings.C().FindOne(context.TODO(), bson.M{"userID": user.ID, "tenantID": user.TenantID}).Decode(settings)
	if err != nil {
		e.Log.Errorf("get settings error: %s", err.Error())
		e.Err(http.StatusInternalServerError, err)
		return
	}
	e.OK(settings.Items)
}

// Update 更新
// @Summary 更新settings
// @Description 更新settings
// @Tags settings
// @Accept  application/json
// @Product application/json
// @Param id path string true "id"
// @Param data body form.SettingsUpdateReq true "data"
// @Success 200 {object} response.Response
// @Router /admin/api/v1/settings [put]
// @Security Bearer
func (e Settings) Update(c *gin.Context) {
	req := &form.SettingsUpdateReq{}
	err := e.Make(c).Bind(req).Error
	if err != nil {
		e.Err(http.StatusUnprocessableEntity, err)
		return
	}
	u := middlewares.GetLoginUser(c)
	if u == nil {
		e.Err(http.StatusUnauthorized, errors.New("user is empty"))
		return
	}
	user, err := models.GetUserByNameAndDomain(c, c.Request.Host, u.Name)
	if err != nil {
		e.Err(http.StatusUnauthorized, err)
		return
	}
	req.TenantID = user.TenantID
	req.UserID = user.ID
	req.UpdatedAt = time.Now()
	settings := &models.Settings{}
	count, err := settings.C().CountDocuments(context.TODO(), bson.M{"userID": settings.ID, "tenantID": settings.TenantID})
	if err != nil {
		e.Log.Errorf("count settings error: %s", err.Error())
		e.Err(http.StatusInternalServerError, err)
		return
	}
	if count == 0 {
		req.ID = primitive.NewObjectID().Hex()
		req.CreatedAt = time.Now()
		_, err = settings.C().InsertOne(context.TODO(), req)
		if err != nil {
			e.Log.Errorf("insert settings error: %s", err.Error())
			e.Err(http.StatusInternalServerError, err)
			return
		}
	} else {
		_, err = settings.C().UpdateOne(context.TODO(), bson.M{"userID": settings.ID, "tenantID": settings.TenantID}, bson.M{"$set": req})
	}
	if err != nil {
		e.Log.Errorf("update settings error: %s", err.Error())
		e.Err(http.StatusInternalServerError, err)
		return
	}
	e.OK(nil)
}

func (e Settings) Other(r *gin.RouterGroup) {

}
