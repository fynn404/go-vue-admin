package service

import (
	"context"
	"errors"

	models "github.com/fynn404/go-vue-admin/internal/model"
	"github.com/fynn404/go-vue-admin/internal/repository"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
)

// UserService 定义用户服务接口
type UserService interface {
	// Register 注册新用户
	Register(ctx context.Context, username, password string, role models.Role) (*models.User, error)
	// Login 用户登录
	Login(ctx context.Context, username, password string) (*models.User, error)
	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, id uint) error
	// GetUserByID 根据ID获取用户信息
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	// ListUsers 获取用户列表
	ListUsers(ctx context.Context, page, pageSize int) ([]*models.User, int64, error)
}

// userService 实现 UserService 接口
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// Register 实现用户注册
func (s *userService) Register(ctx context.Context, username, password string, role models.Role) (*models.User, error) {
	// 检查用户是否已存在
	existingUser, err := s.userRepo.FindByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 创建新用户
	user := &models.User{
		Username: username,
		Password: password,
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 清空密码后返回用户信息
	user.Password = ""
	return user, nil
}

// Login 实现用户登录
func (s *userService) Login(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.ValidatePassword(password) {
		return nil, ErrInvalidPassword
	}

	// 清空密码后返回用户信息
	user.Password = ""
	return user, nil
}

// UpdateUser 实现更新用户信息
func (s *userService) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	// 如果更新包含密码，需要特殊处理
	if password, ok := updates["password"]; ok {
		user.Password = password.(string)
		delete(updates, "password")
	}

	// 更新其他字段
	for key, value := range updates {
		switch key {
		case "username":
			user.Username = value.(string)
		case "role":
			user.Role = value.(models.Role)
		}
	}

	return s.userRepo.Update(ctx, user)
}

// DeleteUser 实现删除用户
func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	if _, err := s.userRepo.FindByID(ctx, id); err != nil {
		return ErrUserNotFound
	}
	return s.userRepo.Delete(ctx, id)
}

// GetUserByID 实现根据ID获取用户信息
func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	user.Password = "" // 清空密码
	return user, nil
}

// ListUsers 实现获取用户列表
func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]*models.User, int64, error) {
	offset := (page - 1) * pageSize
	users, total, err := s.userRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 清空所有用户的密码
	for _, user := range users {
		user.Password = ""
	}

	return users, total, nil
}
