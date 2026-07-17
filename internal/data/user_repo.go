package data

import (
	"errors"
	"login/internal/data/model"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// UserRepository 定义数据访问对象，内部持有数据库连接状态
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 相当于这个“类”的构造函数，用于在服务启动时注入依赖
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindUserByUsername 根据用户名查询用户，用于登录校验
// 注意这里的方法接收者变成了 (r *UserRepository)，它现在是一个对象方法了
func (r *UserRepository) FindUserByUsername(username string) (*model.User, error) {
	var user model.User

	// 相当于执行: SELECT * FROM users WHERE username = ? LIMIT 1;
	// 【核心改变】：使用对象内部的 r.db，彻底告别全局变量 DB
	result := r.db.Where("username = ?", username).First(&user)

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
// 同样，变成了对象方法
func (r *UserRepository) CreateUser(user *model.User) error {
	// 相当于执行: INSERT INTO users ...
	// 【核心改变】：使用 r.db
	err := r.db.Create(user).Error
	if err == nil {
		return nil
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrDuplicateKey
	}
	return err
}
