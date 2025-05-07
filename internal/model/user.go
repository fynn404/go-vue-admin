package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

// User 用户模型
type User struct {
	gorm.Model
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Role      Role      `gorm:"not null" json:"role"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Courses   []Course  `gorm:"foreignKey:TeacherID" json:"courses,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeSave - GORM hook to hash password before saving
func (u *User) BeforeSave(tx *gorm.DB) error {

	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// ValidatePassword - Check if provided password matches the hash
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// TableName - Set table name for GORM
func (User) TableName() string {
	return "users_tab"
}
