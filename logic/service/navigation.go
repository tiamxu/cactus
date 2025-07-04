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

func (s *NavigationService) List(ctx context.Context, req *types.NavigationListReq) (*types.DataListResp, error) {
	var status *int
	if req.Status != nil {
		s := *req.Status
		status = &s
	}
	links, total, err := repo.ListNavigationLinks(ctx, req.Title, req.Category, status, req.PageNo, req.PageSize)
	if err != nil {
		return nil, errors.New("查询导航链接信息失败")
	}

	return &types.DataListResp{
		Total:    total,
		PageData: links,
	}, nil
}

func (s *NavigationService) Get(ctx context.Context, id int) (*model.NavigationLink, error) {
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
	nav.URL = req.URL
	nav.Status = *req.Status

	return repo.UpdateNavigationWithId(ctx, id, nav)
}

func (s *NavigationService) Delete(ctx context.Context, id int) error {
	return repo.DeleteNavigationWithId(ctx, id)
}

func (s *NavigationService) ListPublicNavigationList(page, pageSize int) ([]*model.NavigationLink, int64, error) {
	ctx := context.Background()

	return repo.GetAllLinks(ctx, page, pageSize)
}
