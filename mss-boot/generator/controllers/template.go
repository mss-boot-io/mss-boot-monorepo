/*
 * @Author: lwnmengjing
 * @Date: 2022/3/10 22:43
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/10 22:43
 */

package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mss-boot-io/mss-boot/pkg/response"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/cfg"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/form"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/pkg"
)

func init() {
	e := &Template{}
	response.AppendController(e)
}

type Template struct {
	response.Api
}

func (Template) Path() string {
	return "template"
}

func (e Template) Other(r *gin.RouterGroup) {
	r.GET("/template/get-params", e.GetParams)
	r.POST("/template/generate", e.Generate)
}

// GetParams 获取template参数配置
// @Summary 获取template参数配置
// @Description 获取template参数配置
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Param source query string true "template source"
// @Param branch query string false "branch default:HEAD"
// @Param path query string false "path default:."
// @Success 200 {object} response.Response{data=form.TemplateGetParamsResp}
// @Router /generator/api/v1/template/get-params [get]
// @Security Bearer
func (e Template) GetParams(c *gin.Context) {
	req := &form.TemplateGetParamsReq{}
	err := e.Make(c).Bind(req).Error
	if err != nil {
		e.Err(http.StatusUnprocessableEntity, err)
		return
	}
	if req.Branch == "" {
		req.Branch = "main"
	}
	//获取模版, 存放位置: templates/provider/owner/repo
	dir := fmt.Sprintf("temp/%s/%s", strings.ReplaceAll(
		strings.ReplaceAll(req.Source, "https://", ""),
		"http://",
		"",
	), req.Branch)
	//获取最新代码
	_, err = pkg.GitClone(req.Source, req.Branch, dir, false, cfg.Cfg.Github.PersonalAccessToken)
	//更新
	if err != nil {
		e.Log.Error(err)
		e.Err(http.StatusInternalServerError, err)
		return
	}
	var keys map[string]string
	keys, err = pkg.GetParseFromTemplate(dir, req.Path)
	if err != nil {
		e.Log.Error(err)
		e.Err(http.StatusFailedDependency, err)
		return
	}

	resp := &form.TemplateGetParamsResp{
		Params: keys,
	}
	e.OK(resp)
}

// Generate 从模版生成代码
// @Summary 从模版生成代码
// @Description 从模版生成代码
// @Tags admin
// @Accept  application/json
// @Product application/json
// @Param data body form.TemplateGenerateReq true "data"
// @Success 200 {object} response.Response{data=form.TemplateGenerateResp}
// @Router /generator/api/v1/template/generate [post]
// @Security Bearer
func (e Template) Generate(c *gin.Context) {
	req := &form.TemplateGenerateReq{}
	err := e.Make(c).Bind(req).Error
	if err != nil {
		e.Err(http.StatusUnprocessableEntity, err)
		return
	}
	if req.Template.Branch == "" {
		req.Template.Branch = "main"
	}
	//获取模版, 存放位置: temp/provider/owner/repo
	dir := fmt.Sprintf("temp/%s", strings.ReplaceAll(
		strings.ReplaceAll(req.Template.Source, "https://", ""),
		"http://",
		"",
	))
	fmt.Println(req.Template.Source)
	fmt.Println(req.Template.Branch)
	//获取新代码
	_, err = pkg.GitClone(
		req.Template.Source,
		req.Template.Branch, dir, false,
		cfg.Cfg.Github.PersonalAccessToken)
	if err != nil {
		e.Log.Error(err)
		e.Err(http.StatusInternalServerError, err)
		return
	}

	//获取目的提交项目，存放路径: temp/provider/owner/repo/branch
	branch := fmt.Sprintf("generate/%s", uuid.New().String())
	codeDir := fmt.Sprintf("temp/%s/%s", strings.ReplaceAll(
		strings.ReplaceAll(req.Generate.Repo, "https://", ""),
		"http://",
		"",
	), branch)

	_, err = pkg.GitClone(req.Generate.Repo, "", codeDir, false, cfg.Cfg.Github.PersonalAccessToken)
	if err != nil {
		e.Log.Error(err)
		e.Err(http.StatusInternalServerError, err)
		return
	}
	destination := codeDir
	if req.Generate.Service != "" {
		destination = filepath.Join(destination, req.Generate.Service)
	}
	err = pkg.Generate(&pkg.TemplateConfig{
		Service:       req.Template.Path,
		TemplateLocal: dir,
		Destination:   destination,
		Params:        req.Generate.Params,
	})
	if err != nil {
		e.Log.Error(err)
		e.Err(http.StatusInternalServerError, err)
		return
	}
	err = pkg.CommitAndPushGithubRepo(codeDir, branch, cfg.Cfg.Github.PersonalAccessToken)
	if err != nil {
		e.Log.Error(err)
		e.Err(http.StatusInternalServerError, err)
		return
	}
	resp := &form.TemplateGenerateResp{
		Repo:   req.Generate.Repo,
		Branch: branch,
	}
	e.OK(resp)
}
