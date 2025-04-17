package service

import (
	"errors"
	"fmt"

	"github.com/tiamxu/cactus/models"
	"github.com/tiamxu/cactus/utils"
	"github.com/tiamxu/kit/log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}
func (s *AuthService) Authenticate(username, password string) (*models.User, string, error) {
	user, err := models.GetUserByUsername(username)

	if err != nil {
		log.Errorf("数据库查询错误: %v\n", err)
		return nil, "", errors.New("用户查询失败")
	}
	if user == nil {
		log.Infoln("用户不存在")
		return nil, "", errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("密码验证失败")
	}

	token := utils.GenerateToken(user.ID)

	return user, token, nil
}

func (s *AuthService) ChangePassword(uid int, oldPwd, newPwd string) error {
	// 1. 获取当前密码哈希
	currentPasswordHash, err := models.GetPasswordHash(uid)
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
	if err := models.UpdatePassword(uid, string(newPasswordHash)); err != nil {
		return err
	}

	return nil
}
