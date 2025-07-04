package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tiamxu/cactus/logic/model"
	"github.com/tiamxu/cactus/logic/repo"
	"github.com/tiamxu/cactus/types"
)

type EnvironmentService struct{}

func NewEnvironmentService() *EnvironmentService {
	return &EnvironmentService{}
}

// CreateEnvironment 创建环境
func (s *EnvironmentService) Add(ctx context.Context, req *types.EnvironmentCreateReq) (*model.Environment, error) {
	// 验证参数
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("环境名称不能为空")
	}
	if strings.TrimSpace(req.Code) == "" {
		return nil, errors.New("环境代码不能为空")
	}

	env := &model.Environment{
		Name:        req.Name,
		Code:        strings.ToUpper(req.Code),
		Description: req.Description,
		Status:      1, // 默认启用
	}

	if err := repo.CreateEnvironment(ctx, env); err != nil {
		return nil, fmt.Errorf("创建环境失败: %v", err)
	}

	return env, nil
}

// GetEnvironment 获取环境详情
func (s *EnvironmentService) Get(ctx context.Context, id int) (*model.Environment, error) {
	env, err := repo.GetEnvironmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取环境失败: %v", err)
	}
	return env, nil
}

// ListEnvironments 获取环境列表
func (s *EnvironmentService) List(ctx context.Context, req *types.EnvironmentListReq) (*types.DataListResp, error) {
	var status *int
	if req.Status != nil {
		s := *req.Status
		status = &s
	}

	envs, total, err := repo.ListEnvironments(ctx, req.Name, req.Code, status, req.PageNo, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("查询环境列表失败: %v", err)
	}

	return &types.DataListResp{
		PageData: envs,
		Total:    total,
	}, nil
}

// UpdateEnvironment 更新环境
func (s *EnvironmentService) Update(ctx context.Context, id int, req *types.EnvironmentUpdateReq) error {
	env, err := repo.GetEnvironmentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("环境不存在")
	}

	if req.Name != nil {
		env.Name = *req.Name
	}
	if req.Code != nil {
		env.Code = *req.Code
	}
	if req.Description != nil {
		env.Description = *req.Description
	}
	if req.Status != nil {
		env.Status = *req.Status
	}

	return repo.UpdateEnvironment(ctx, env)
}

// DeleteEnvironment 删除环境
func (s *EnvironmentService) Delete(ctx context.Context, id int) error {

	return repo.DeleteEnvironment(ctx, id)
}
