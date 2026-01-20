package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	FirstName string    `gorm:"not null" json:"firstName"`
	LastName  string    `gorm:"not null" json:"lastName"`
	Role      UserRole  `gorm:"type:varchar(20);not null" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserRole string

const (
	RoleEmployee   UserRole = "employee"
	RoleManagement UserRole = "management"
)
