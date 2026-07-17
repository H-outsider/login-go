package service

import (
	"errors"
	"login/api"
	"login/internal/data"
	"login/internal/data/model"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserExists         = errors.New("用户名已存在")
)

type UserRepository interface {
	FindUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) error
}

// UserService 定义业务逻辑对象，内部持有数据仓储层 (Repo) 的引用
type UserService struct {
	repo UserRepository
}

// NewUserService 构造函数，用于依赖注入
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterService 处理用户注册逻辑
// 变成 UserService 对象的方法
func (s *UserService) RegisterService(req api.RegisterRequest) error {
	// 1. 检查用户是否已存在
	// 【核心改变】：使用注入进来的 s.repo 对象调用方法
	existUser, err := s.repo.FindUserByUsername(req.Username)
	if err != nil {
		return err // 数据库查询出错
	}
	if existUser != nil {
		return ErrUserExists
	}

	// 2. 密码加密 (Bcrypt)
	// GenerateFromPassword 第二个参数是加密强度（默认 10 即可）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 3. 组装实体类并保存到数据库
	newUser := model.User{
		Username: req.Username,
		Password: string(hashedPassword), // 存入加密后的哈希值
		Email:    req.Email,
	}

	// 【核心改变】：使用 s.repo
	if err := s.repo.CreateUser(&newUser); err != nil {
		if errors.Is(err, data.ErrDuplicateKey) {
			return ErrUserExists
		}
		return err
	}
	return nil
}

// LoginService 处理用户登录逻辑
// 变成 UserService 对象的方法
func (s *UserService) LoginService(req api.LoginRequest) (*api.UserResponse, error) {
	// 1. 根据用户名查找用户
	// 【核心改变】：使用 s.repo
	user, err := s.repo.FindUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials // 故意模糊提示，防暴力破解
	}

	// 2. 校验密码是否正确
	// CompareHashAndPassword 专门用于比对明文和哈希值
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 3. 登录成功，组装返回给前端的数据 (脱敏)
	return &api.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
