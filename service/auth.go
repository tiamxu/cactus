package service

import (
	"errors"

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
