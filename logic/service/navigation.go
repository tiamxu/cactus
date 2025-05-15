package service

import (
	"sort"

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

func (s *NavigationService) GetAllLinks() ([]model.NavigationLink, error) {

	return s.db.GetAllLinks()
}

func (s *NavigationService) GetLinkByID(id int) (model.NavigationLink, error) {
	return s.db.GetLinkByID(id)
}

func (s *NavigationService) CreateLink(req model.CreateLinkRequest) (int, error) {
	return s.db.CreateLink(req)
}

func (s *NavigationService) UpdateLink(id int, req model.UpdateLinkRequest) error {
	return s.db.UpdateLink(id, req)
}

func (s *NavigationService) DeleteLink(id int) error {
	return s.db.DeleteLink(id)
}

func (s *NavigationService) RenderIndexPage() ([]inout.GroupedLink, error) {
	links, err := s.db.GetAllLinks()
	if err != nil {
		return nil, err
	}
	groups := make(map[string][]model.NavigationLink)
	// 先按category分组
	for _, link := range links {
		category := link.Category
		if category == "" {
			category = "未分类"
		}
		groups[category] = append(groups[category], link)
	}
	// 转换为切片并排序
	var result []inout.GroupedLink
	for category, links := range groups {
		result = append(result, inout.GroupedLink{
			Category: category,
			Links:    links,
		})
	}
	// 按category名称排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Category < result[j].Category
	})
	return result, nil
}
