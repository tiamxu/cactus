package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tiamxu/cactus/logic/model"
	"github.com/tiamxu/cactus/logic/repo"
	"github.com/tiamxu/cactus/types"
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
func (u *UserService) GetUserDetail(userId int) (*types.UserDetailRes, error) {
	var res types.UserDetailRes

	// 查询用户信息
	user, err := repo.GetUserByID(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	res.User = *user

	// 查询用户详情
	profile, err := repo.GetProfileByUserID(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res.Profile = &model.Profile{} // 如果没有找到 profile，可以设置为空结构体
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
	roleIDs, err := repo.GetUserRoleIDs(userId)
	if err != nil {
		return nil, err
	}
	// 查询角色信息
	if len(roleIDs) > 0 {
		roles, err := repo.GetRolesByID(roleIDs)
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

func (u *UserService) GetUserList(gender, enable, username string, pageNo, pageSize int) (*types.UserListRes, error) {
	var data = types.UserListRes{
		PageData: make([]types.UserListItem, 0),
	}
	profiles, total, err := repo.GetProfilesByCondition(gender, enable, username, pageNo, pageSize)
	if err != nil {
		return nil, errors.New("查询用户资料信息失败")
	}
	data.Total = total

	for _, profile := range profiles {

		uinfo, err := repo.GetUserByID(profile.UserId)
		if err != nil {
			return nil, errors.New("查询用户信息失败")
		}
		roles, err := repo.GetRolesByUserId(profile.UserId)
		if err != nil {
			return nil, errors.New("查询用户角色失败")
		}
		// 组装返回数据
		data.PageData = append(data.PageData, types.UserListItem{
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

func (u *UserService) UpdateProfile(params types.PatchProfileUserReq) error {
	a := model.Profile{
		ID:       params.Id,
		Gender:   params.Gender,
		Address:  params.Address,
		Email:    params.Email,
		NickName: params.NickName,
	}
	err := repo.UpdateProfileByWhere(a)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserService) Update(params types.PatchUserReq) error {
	err := repo.UpdateUserByWhere(params.Id, params.Username, params.Password, params.Enable, params.RoleIds)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserService) Add(params types.AddUserReq) error {
	err := repo.AddUserByWhere(params.Username, params.Password, params.Enable, params.RoleIds)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Delete(uid int) error {
	err := repo.DeleteUserByWhere(uid)
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

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

// 根据用户ID查找对应的profile
func findProfileByUserId(profiles []model.Profile, userId int) *model.Profile {
	for _, p := range profiles {
		if p.UserId == userId {
			return &p
		}
	}
	return &model.Profile{} // 返回空profile避免nil
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
