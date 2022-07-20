/*
 * @Author: lwnmengjing
 * @Date: 2022/3/10 14:23
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/10 14:23
 */

package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/controllers"
	_ "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/docs"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/middleware"
	"github.com/mss-boot-io/mss-boot/pkg/response"
)

func Init(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	middleware.Init(v1)
	for i := range response.Controllers {
		response.Controllers[i].Other(v1)
	}
}
