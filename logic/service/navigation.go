package service

import (
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

func (s *NavigationService) List(pageNo, pageSize int) (*inout.NavListRes, error) {
	var data = inout.NavListRes{
		PageData: make([]model.NavigationLink, 0),
	}
	navs, total, err := s.db.GetAllLinks(pageNo, pageSize)
	if err != nil {
		return nil, errors.New("查询角色信息失败")
	}
	data.Total = total
	data.PageData = navs
	return &data, nil
}

func (s *NavigationService) GetLinkByID(id int) (model.NavigationLink, error) {
	return s.db.GetLinkByID(id)
}

func (s *NavigationService) Add(req inout.CreateLinkRequest) error {
	return s.db.Create(req)
}

func (s *NavigationService) Update(id int, req inout.UpdateLinkRequest) error {
	return s.db.UpdateNavigationWithId(id, req)
}

func (s *NavigationService) Delete(id int) error {
	return s.db.DeleteNavigationWithId(id)
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
