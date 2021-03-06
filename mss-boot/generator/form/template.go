package form

type TemplateGetParamsReq struct {
	Source string `query:"source" form:"source" binding:"required"`
	Branch string `query:"branch" form:"branch"`
	Path   string `query:"path" form:"path"`
}

type TemplateGetParamsResp struct {
	Params map[string]string `json:"params"`
}

type TemplateGenerateReq struct {
	Template TemplateParams `json:"template"`
	Generate GenerateParams `json:"generate"`
}

type TemplateParams struct {
	Source string `json:"source" binding:"required"`
	Branch string `json:"branch"`
	Path   string `json:"path"`
}

type GenerateParams struct {
	Service string            `json:"service"`
	Repo    string            `json:"repo" binding:"required"`
	Params  map[string]string `json:"params"`
}

type TemplateGenerateResp struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}
