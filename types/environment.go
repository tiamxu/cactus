package types

type EnvironmentCreateReq struct {
	Name        string `form:"name" json:"name" binding:"required"` // 环境名称
	Code        string `form:"code" json:"code" binding:"required"` // 环境代码
	Description string `form:"description" json:"description"`      // 环境描述
}

type EnvironmentUpdateReq struct {
	Name        *string `json:"name"`        // 环境名称
	Code        *string `json:"code"`        // 环境代码
	Description *string `json:"description"` // 环境描述
	Status      *int    `json:"status"`      // 状态 (0-禁用, 1-启用)
}

// ListEnvironmentRequest 环境列表请求
type EnvironmentListReq struct {
	Name   string `form:"name"` // 环境名称(模糊查询)
	Code   string `form:"code"`
	Status *int   `form:"status"` // 状态 (0-禁用, 1-启用)
	BasePage
}
