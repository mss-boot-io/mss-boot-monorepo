package form

type TemplateGetBranchesReq struct {
	Source string `query:"source" form:"source" binding:"required"`
}

type TemplateGetBranchesResp struct {
	Branches []string `json:"branches"`
}

type TemplateGetPathReq struct {
	Source string `query:"source" form:"source" binding:"required"`
	Branch string `query:"branch" form:"branch"`
}

type TemplateGetPathResp struct {
	Path []string `json:"path"`
}

type TemplateGetParamsReq struct {
	Source string `query:"source" form:"source" binding:"required"`
	Branch string `query:"branch" form:"branch"`
	Path   string `query:"path" form:"path"`
}

type TemplateGetParamsResp struct {
	Params []TemplateParam `json:"params"`
}

type TemplateParam struct {
	Name string `json:"name"`
	Tip  string `json:"tip"`
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
