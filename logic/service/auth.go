package service

import (
	"errors"
	"fmt"

	"github.com/tiamxu/cactus/logic/model"
	"github.com/tiamxu/cactus/logic/repo"

	"github.com/tiamxu/cactus/utils"
	"github.com/tiamxu/kit/log"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound        = errors.New("用户不存在")
	ErrPasswordMismatch    = errors.New("密码验证失败")
	ErrDatabaseQueryFailed = errors.New("数据库查询失败")
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}
func (s *AuthService) Authenticate(username, password string) (*model.User, string, error) {
	user, err := repo.GetUserByUsername(username)

	if err != nil {
		log.Errorf("数据库查询错误: %v", err)
		return nil, "", fmt.Errorf("用户查询失败: %w", err)
	}
	if user == nil {
		log.Infoln("用户不存在")
		return nil, "", ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Warnf("密码验证失败: %v", err)
		return nil, "", ErrPasswordMismatch
	}

	token := utils.GenerateToken(user.ID)

	return user, token, nil
}

func (s *AuthService) ChangePassword(uid int, oldPwd, newPwd string) error {
	// 1. 获取当前密码哈希
	currentPasswordHash, err := repo.GetPasswordHash(uid)
	if err != nil {
		return err
	}

	// 2. 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(currentPasswordHash), []byte(oldPwd)); err != nil {
		return fmt.Errorf("当前密码不正确")
	}

	// 3. 生成新密码哈希
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码哈希失败: %w", err)
	}

	// 4. 更新密码
	if err := repo.UpdatePassword(uid, string(newPasswordHash)); err != nil {
		return err
	}

	return nil
}
