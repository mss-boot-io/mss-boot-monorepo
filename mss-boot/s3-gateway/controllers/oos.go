/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/24 14:59:13
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/24 14:59:13
 */

package controllers

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/s3-gateway/cfg"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/pkg/response"
)

func init() {
	e := &OOS{}
	//e.Auth = false
	response.AppendController(e)
}

type OOS struct {
	response.Api
	//curd.DefaultController
}

func (e *OOS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	key = key[1:]
	if key == "" {
		key = "index.html"
	}
	client := cfg.Cfg.S3.GetClient()

	output, err := client.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: &cfg.Cfg.S3.Bucket,
		Key:    &key,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, fmt.Sprintf("get object failed, err:%s", err.Error()))
		return
	}
	defer output.Body.Close()
	w.Header().Set("Content-Type", *output.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", output.ContentLength))

	_, err = io.Copy(w, output.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, fmt.Sprintf("copy object failed, err:%s", err.Error()))
		return
	}
}

func (*OOS) Path() string {
	return "/"
}

func (e *OOS) Other(r *gin.RouterGroup) {
	//r.Use((&middlewares.AuthMiddleware{}).AuthMiddleware())
	r.GET("/:key", e.GetObject)
}

// GetObject somebody
// @Summary Call somebody
// @Description Call somebody
// @Tags s3-gateway
// @Accept  application/json
// @Product application/json
// @Param data body form.HelloworldCallReq true "data"
// @Success 200 {object} response.Response{data=form.HelloworldCallResp}
// @Router /s3-gateway/api/v1/call [post]
func (e OOS) GetObject(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		key = "index.html"
	}
	client := cfg.Cfg.S3.GetClient()

	output, err := client.GetObject(c, &s3.GetObjectInput{
		Bucket: &cfg.Cfg.S3.Bucket,
		Key:    &key,
	})
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/plain", []byte(fmt.Sprintf("get object failed, err:%v", err)))
		return
	}
	c.DataFromReader(http.StatusOK, output.ContentLength, *output.ContentType, output.Body, nil)
}
