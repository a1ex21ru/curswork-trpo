package models

// UpdateExpenseStatusDTO for updating expense request status
type UpdateExpenseStatusDTO struct {
	Status   RequestStatus `json:"status" binding:"required,oneof=approved rejected"`
	Comments string        `json:"comments" binding:"required,min=10"`
}

// RegisterUserDTO for user registration
type RegisterUserDTO struct {
	Email     string   `json:"email" binding:"required,email"`
	Password  string   `json:"password" binding:"required,min=6"`
	FirstName string   `json:"firstName" binding:"required"`
	LastName  string   `json:"lastName" binding:"required"`
	Role      UserRole `json:"role" binding:"required,oneof=employee management"`
}

// LoginDTO for user login
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse for login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// StatsResponse for statistics
type StatsResponse struct {
	TotalPending      float64 `json:"totalPending"`
	TotalApproved     float64 `json:"totalApproved"`
	PendingCount      int     `json:"pendingCount"`
	ApprovedThisMonth int     `json:"approvedThisMonth"`
	BudgetUsed        float64 `json:"budgetUsed"`
	BudgetRemaining   float64 `json:"budgetRemaining"`
}
