package data

import (
	"errors"
	"login/internal/data/model"

	"gorm.io/gorm"
)

// FindUserByUsername 根据用户名查询用户，用于登录校验
func FindUserByUsername(username string) (*model.User, error) {
	var user model.User

	// 相当于执行: SELECT * FROM users WHERE username = ? LIMIT 1;
	result := DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		// 区分是“没找到数据”还是“数据库挂了”
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, result.Error // 其他底层错误
	}

	return &user, nil
}

// CreateUser 创建新用户，用于注册
func CreateUser(user *model.User) error {
	// 相当于执行: INSERT INTO users ...
	return DB.Create(user).Error
}
