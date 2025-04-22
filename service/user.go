package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/models"
)

// 定义业务错误
var (
	// ErrUserNotFound   = errors.New("user not found")
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
func (u *UserService) GetUserDetail(userId int) (*inout.UserDetailRes, error) {
	var res inout.UserDetailRes

	// 查询用户信息
	user, err := models.GetUserByID(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	res.User = *user

	// 查询用户详情
	profile, err := models.GetProfileByUserID(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res.Profile = &models.Profile{} // 如果没有找到 profile，可以设置为空结构体
		} else {
			return nil, err
		}
	}
	res.Profile = profile

	// 查询用户角色 ID 列表
	// roleIDs, err := models.GetRolesIdByUserID(userId)
	// if err != nil {
	// 	return nil, err
	// }
	roleIDs, err := models.GetUserRoleIDs(userId)
	if err != nil {
		return nil, err
	}
	// 查询角色信息
	if len(roleIDs) > 0 {
		roles, err := models.GetRolesByID(roleIDs)
		if err != nil {
			return nil, err
		}
		res.Roles = roles
	}

	// 设置当前角色
	if len(res.Roles) > 0 {
		res.CurrentRole = res.Roles[0]
	}

	return &res, nil
}

func (u *UserService) GetUserList(gender, enable, username string, pageNo, pageSize int) (*inout.UserListRes, error) {
	var data = inout.UserListRes{
		PageData: make([]inout.UserListItem, 0),
	}
	profiles, total, err := models.GetProfilesByCondition(gender, enable, username, pageNo, pageSize)
	if err != nil {
		return nil, errors.New("查询用户资料信息失败")
	}
	data.Total = total

	for _, profile := range profiles {

		uinfo, err := models.GetUserByID(profile.UserId)
		if err != nil {
			return nil, errors.New("查询用户信息失败")
		}
		roles, err := models.GetRolesByUserId(profile.UserId)
		if err != nil {
			return nil, errors.New("查询用户角色失败")
		}
		// 组装返回数据
		data.PageData = append(data.PageData, inout.UserListItem{
			ID:         uinfo.ID,
			Username:   uinfo.Username,
			Enable:     uinfo.Enable,
			CreateTime: uinfo.CreateTime,
			UpdateTime: uinfo.UpdateTime,
			Gender:     profile.Gender,
			Avatar:     profile.Avatar,
			Address:    profile.Address,
			Email:      profile.Email,
			Roles:      roles,
		})

	}

	return &data, nil
}

func (u *UserService) UpdateProfile(params inout.PatchProfileUserReq) error {
	a := models.Profile{
		ID:       params.Id,
		Gender:   params.Gender,
		Address:  params.Address,
		Email:    params.Email,
		NickName: params.NickName,
	}
	err := models.UpdateProfileByWhere(a)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserService) Update(params inout.PatchUserReq) error {
	err := models.UpdateUserByWhere(params.Id, params.Username, params.Password, params.Enable, params.RoleIds)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserService) Add(params inout.AddUserReq) error {
	err := models.AddUserByWhere(params.Username, params.Password, params.Enable, params.RoleIds)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Delete(uid int) error {
	err := models.DeleteUserByWhere(uid)
	if err != nil {
		return err
	}
	return nil
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
// func (s *UserService) Create(req *CreateUserRequest) error {
// 	// 参数校验
// 	if req.Username == "" || req.Password == "" {
// 		return ErrInvalidRequest
// 	}
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
// 	if err != nil {
// 		return err
// 	}
// 	// 检查用户名唯一性
// 	exists, err := models.ExistsByUsername(req.Username)
// 	if err != nil {
// 		return fmt.Errorf("数据库查询失败: %v", err)
// 	}
// 	if exists {
// 		return fmt.Errorf("用户名 %s 已存在", req.Username)
// 	}

// 	user := &models.User{
// 		Username: req.Username,
// 		Password: string(hashedPassword),
// 	}

// 	return models.Create(user)
// }

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

// func (s *UserService) Update(id uint, req *UpdateUserRequest) (*models.User, error) {
// 	user, err := models.GetByID(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("用户不存在")
// 	}

// 	// 更新字段
// 	if req.Username != "" && req.Username != user.Username {
// 		exists, err := models.ExistsByUsername(req.Username)
// 		if err != nil {
// 			return nil, fmt.Errorf("数据库查询失败: %v", err)
// 		}
// 		if exists {
// 			return nil, fmt.Errorf("用户名 %s 已存在", req.Username)
// 		}
// 		user.Username = req.Username
// 	}

// 	if err := models.Update(user); err != nil {
// 		return nil, fmt.Errorf("更新失败: %v", err)
// 	}
// 	return user, nil
// }

// Delete 删除用户（业务校验）

// 根据用户ID查找对应的profile
func findProfileByUserId(profiles []models.Profile, userId int) *models.Profile {
	for _, p := range profiles {
		if p.UserId == userId {
			return &p
		}
	}
	return &models.Profile{} // 返回空profile避免nil
}

// func (s *UserService) Delete(id uint) error {
// 	if id == 0 {
// 		return ErrInvalidRequest
// 	}

// 	// 先检查用户是否存在
// 	if _, err := models.GetByID(id); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return ErrUserNotFound
// 		}
// 		return fmt.Errorf("check user existence failed: %w", err)
// 	}

// 	if err := models.Delete(id); err != nil {
// 		return fmt.Errorf("delete operation failed: %w", err)
// 	}
// 	return nil
// }

// List 获取用户列表（带分页）
// func (s *UserService) List(page, pageSize int) ([]*User, error) {
// 	// 参数校验
// 	if page < 1 || pageSize < 1 || pageSize > 100 {
// 		return nil, ErrInvalidRequest
// 	}

// 	dbUsers, err := models.List(page, pageSize)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list users: %w", err)
// 	}

// 	users := make([]*User, 0, len(dbUsers))
// 	for _, u := range dbUsers {
// 		users = append(users, &User{
// 			ID:       u.ID,
// 			Username: u.Username,
// 		})
// 	}
// 	return users, nil
// }
