package domain

import "time"

type User struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:name;type:varchar(100);not null;index:idx_name" json:"name"`
	Email     string     `gorm:"column:email;type:varchar(100);uniqueIndex:idx_email;not null" json:"email"`
	Password  string     `gorm:"column:password;type:varchar(255);not null" json:"-"`
	CreatedAt time.Time  `gorm:"autoCreateTime;index:idx_created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index:idx_deleted_at" json:"-"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"email"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
