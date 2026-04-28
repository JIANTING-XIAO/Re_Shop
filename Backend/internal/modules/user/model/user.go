package model

import "time"

type User struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(50);not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Nickname  string    `gorm:"column:nickname;type:varchar(100)" json:"nickname"`
	Avatar    string    `gorm:"column:avatar;type:varchar(500)" json:"avatar"`
	Phone     string    `gorm:"column:phone;type:varchar(20);uniqueIndex" json:"phone"`
	Role      int8      `gorm:"column:role;type:tinyint;not null;default:0" json:"role"`
	Status    int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
