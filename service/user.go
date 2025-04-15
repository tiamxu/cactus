package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tiamxu/cactus/models"
	"golang.org/x/crypto/bcrypt"
)

// 定义业务错误
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUsernameExists = errors.New("username already exists")
	ErrInvalidRequest = errors.New("invalid request parameters")
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

// User 业务层结构体
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // 仅在创建/更新时使用明文
}
type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}
func (u *UserService) ListUsers() (user []models.User, err error) {
	user, err = models.ListUsers()
	if err != nil {
		return nil, err
	}
	return user, nil

}

// GetByID 获取用户详情
func (s *UserService) GetByID(id uint) (*User, error) {
	if id == 0 {
		return nil, ErrInvalidRequest
	}

	dbUser, err := s.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
	}, nil
}

// Create 创建用户（带业务校验）
func (s *UserService) Create(req *CreateUserRequest) error {
	// 参数校验
	if req.Username == "" || req.Password == "" {
		return ErrInvalidRequest
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return err
	}
	// 检查用户名唯一性
	exists, err := models.ExistsByUsername(req.Username)
	if err != nil {
		return fmt.Errorf("数据库查询失败: %v", err)
	}
	if exists {
		return fmt.Errorf("用户名 %s 已存在", req.Username)
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	return models.Create(user)
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

func (s *UserService) Update(id uint, req *UpdateUserRequest) (*models.User, error) {
	user, err := models.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 更新字段
	if req.Username != "" && req.Username != user.Username {
		exists, err := models.ExistsByUsername(req.Username)
		if err != nil {
			return nil, fmt.Errorf("数据库查询失败: %v", err)
		}
		if exists {
			return nil, fmt.Errorf("用户名 %s 已存在", req.Username)
		}
		user.Username = req.Username
	}

	if err := models.Update(user); err != nil {
		return nil, fmt.Errorf("更新失败: %v", err)
	}
	return user, nil
}

// Delete 删除用户（业务校验）

func (s *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidRequest
	}

	// 先检查用户是否存在
	if _, err := models.GetByID(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("check user existence failed: %w", err)
	}

	if err := models.Delete(id); err != nil {
		return fmt.Errorf("delete operation failed: %w", err)
	}
	return nil
}

// List 获取用户列表（带分页）
func (s *UserService) List(page, pageSize int) ([]*User, error) {
	// 参数校验
	if page < 1 || pageSize < 1 || pageSize > 100 {
		return nil, ErrInvalidRequest
	}

	dbUsers, err := models.List(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*User, 0, len(dbUsers))
	for _, u := range dbUsers {
		users = append(users, &User{
			ID:       u.ID,
			Username: u.Username,
		})
	}
	return users, nil
}
