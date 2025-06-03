package service

import (
	"context"
	"errors"

	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/logic/model"
	"github.com/tiamxu/cactus/logic/repo"
)

type CreateLinkRequest struct {
	Title       string `json:"title" binding:"required"`
	URL         string `json:"url" binding:"required,url"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type NavigationService struct {
	db *repo.NavigationDB
}

func NewNavigationService(db *repo.NavigationDB) *NavigationService {
	return &NavigationService{db: db}
}

func (s *NavigationService) List(ctx context.Context, pageNo, pageSize int) (*inout.NavListRes, error) {

	links, total, err := s.db.GetAllLinks(ctx, pageNo, pageSize)
	if err != nil {
		return nil, errors.New("查询导航链接信息失败")
	}

	return &inout.NavListRes{
		Total:    total,
		PageData: links,
	}, nil
}

func (s *NavigationService) GetLinkByID(ctx context.Context, id int) (model.NavigationLink, error) {
	return s.db.GetLinkByID(ctx, id)
}

func (s *NavigationService) Add(ctx context.Context, req inout.CreateLinkRequest) error {
	return s.db.Create(ctx, req)
}

func (s *NavigationService) Update(ctx context.Context, id int, req inout.UpdateLinkRequest) error {
	return s.db.UpdateNavigationWithId(ctx, id, req)
}

func (s *NavigationService) Delete(ctx context.Context, id int) error {
	return s.db.DeleteNavigationWithId(ctx, id)
}

// func (s *NavigationService) RenderIndexPage() ([]inout.GroupedLink, error) {
// 	links, err := s.db.GetAllLinks()
// 	if err != nil {
// 		return nil, err
// 	}
// 	groups := make(map[string][]model.NavigationLink)
// 	// 先按category分组
// 	for _, link := range links {
// 		category := link.Category
// 		if category == "" {
// 			category = "未分类"
// 		}
// 		groups[category] = append(groups[category], link)
// 	}
// 	// 转换为切片并排序
// 	var result []inout.GroupedLink
// 	for category, links := range groups {
// 		result = append(result, inout.GroupedLink{
// 			Category: category,
// 			Links:    links,
// 		})
// 	}
// 	// 按category名称排序
// 	sort.Slice(result, func(i, j int) bool {
// 		return result[i].Category < result[j].Category
// 	})
// 	return result, nil
// }
