package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/tiamxu/cactus/logic/model"
	"github.com/tiamxu/cactus/logic/repo"
	"github.com/tiamxu/cactus/types"
)

type NavigationService struct {
}

func NewNavigationService() *NavigationService {
	return &NavigationService{}
}

func (s *NavigationService) List(ctx context.Context, pageNo, pageSize int) (*types.DataListResp, error) {

	links, total, err := repo.GetAllLinks(ctx, pageNo, pageSize)
	if err != nil {
		return nil, errors.New("查询导航链接信息失败")
	}

	return &types.DataListResp{
		Total:    total,
		PageData: links,
	}, nil
}

func (s *NavigationService) GetLinkByID(ctx context.Context, id int) (*model.NavigationLink, error) {
	return repo.GetLinkByID(ctx, id)
}

func (s *NavigationService) Add(ctx context.Context, req *types.NavigationCreateReq) error {
	nav := &model.NavigationLink{
		Title:       req.Title,
		URL:         req.URL,
		Icon:        req.Icon,
		Category:    req.Category,
		Description: req.Description,
	}
	return repo.InsertLink(ctx, nav)
}

func (s *NavigationService) Update(ctx context.Context, id int, req *types.NavigationUpdateReq) error {
	nav, err := repo.GetLinkByID(ctx, id)
	if err != nil {
		return fmt.Errorf("nav信息不存在")
	}

	nav.Title = req.Title
	nav.Icon = req.Icon
	nav.Description = req.Description
	nav.Category = req.Category

	return repo.UpdateNavigationWithId(ctx, id, nav)
}

func (s *NavigationService) Delete(ctx context.Context, id int) error {
	return repo.DeleteNavigationWithId(ctx, id)
}
