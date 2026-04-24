package models

type Role struct {
ID   uint   `gorm:"primaryKey"`
Name string `gorm:"not null;unique"`
}

type UserRole struct {
ID     uint `gorm:"primaryKey"`
UserID uint `gorm:"not null"`
RoleID uint `gorm:"not null"`
}

func (Role) TableName() string {
return "roles"
}

func (UserRole) TableName() string {
return "user_roles"
}
