package models

import "time"

// ExpenseRequest represents an expense request
type ExpenseRequest struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	Title       string        `gorm:"not null" json:"title"`
	Category    string        `gorm:"not null" json:"category"`
	Amount      float64       `gorm:"not null" json:"amount"`
	Vendor      string        `gorm:"not null" json:"vendor"`
	Description string        `gorm:"type:text;not null" json:"description"`
	Status      RequestStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	EmployeeID  uint          `gorm:"not null" json:"employeeId"`
	Employee    User          `gorm:"foreignKey:EmployeeID" json:"employee"`
	ReviewerID  *uint         `json:"reviewerId,omitempty"`
	Reviewer    *User         `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Comments    string        `gorm:"type:text" json:"comments,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	ReviewedAt  *time.Time    `json:"reviewedAt,omitempty"`
}

type RequestStatus string

const (
	StatusPending  RequestStatus = "pending"
	StatusApproved RequestStatus = "approved"
	StatusRejected RequestStatus = "rejected"
)

// Budget represents monthly budget information
type Budget struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Year      int       `gorm:"not null" json:"year"`
	Month     int       `gorm:"not null" json:"month"`
	Total     float64   `gorm:"not null" json:"total"`
	Spent     float64   `gorm:"not null;default:0" json:"spent"`
	Remaining float64   `gorm:"not null" json:"remaining"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateExpenseRequestDTO for creating new expense requests
type CreateExpenseRequestDTO struct {
	Title       string  `json:"title" binding:"required,min=3"`
	Category    string  `json:"category" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Vendor      string  `json:"vendor" binding:"required,min=2"`
	Description string  `json:"description" binding:"required,min=10"`
}

type TopExpenseRequest struct {
	Title    string        `json:"title" binding:"required,min=3"`
	Category string        `json:"category" binding:"required"`
	Amount   float64       `json:"amount" binding:"required,gt=0"`
	Status   RequestStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
}

type CategoryExpenseRequest struct {
	Category    string  `json:"category_name" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required,gt=0"`
	Count       int     `json:"count" binding:"required,gt=0"`
}

type Report struct {
	TopExpense []TopExpenseRequest      `json:"top_expenses"`
	Category   []CategoryExpenseRequest `json:"expenses_by_category"`
}
