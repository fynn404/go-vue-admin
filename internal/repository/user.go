package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"

	models "github.com/fynn404/go-vue-admin/internal/model"
)

// UserRepository 定义用户数据访问接口
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context, offset, limit int) ([]*models.User, int64, error)
	ListByRole(ctx context.Context, role models.Role, offset, limit int) ([]*models.User, int64, error)
	Search(ctx context.Context, keyword string, offset, limit int) ([]*models.User, int64, error)
	BatchDelete(ctx context.Context, ids []uint) error
}

// userRepository 实现 UserRepository 接口
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户 Repository 实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户信息
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListByRole 根据角色查询用户
func (r *userRepository) ListByRole(ctx context.Context, role models.Role, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", role).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("role = ?", role).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Search 搜索用户
func (r *userRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// 构建模糊搜索条件
	keyword = "%" + strings.ToLower(keyword) + "%"
	query := r.db.WithContext(ctx).Where(
		"LOWER(username) LIKE ? OR LOWER(email) LIKE ?",
		keyword, keyword,
	)

	if err := query.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// BatchDelete 批量删除用户
func (r *userRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, ids).Error
}
