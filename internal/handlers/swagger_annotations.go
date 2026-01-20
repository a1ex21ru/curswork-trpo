package handlers

// Swagger model definitions

// @Description User model
type SwaggerUser struct {
	ID        uint   `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	FirstName string `json:"firstName" example:"Иван"`
	LastName  string `json:"lastName" example:"Иванов"`
	Role      string `json:"role" example:"employee"`
} // @name User

// @Description Expense Request model
type SwaggerExpenseRequest struct {
	ID          uint         `json:"id" example:"1"`
	Title       string       `json:"title" example:"Закупка офисной мебели"`
	Category    string       `json:"category" example:"furniture"`
	Amount      float64      `json:"amount" example:"45000"`
	Vendor      string       `json:"vendor" example:"IKEA"`
	Description string       `json:"description" example:"Необходимо приобрести 3 рабочих стола"`
	Status      string       `json:"status" example:"pending"`
	EmployeeID  uint         `json:"employeeId" example:"1"`
	Employee    SwaggerUser  `json:"employee"`
	ReviewerID  *uint        `json:"reviewerId,omitempty" example:"2"`
	Reviewer    *SwaggerUser `json:"reviewer,omitempty"`
	Comments    string       `json:"comments,omitempty" example:"Одобрено"`
	CreatedAt   string       `json:"createdAt" example:"2025-01-15T10:30:00Z"`
	UpdatedAt   string       `json:"updatedAt" example:"2025-01-15T10:30:00Z"`
	ReviewedAt  *string      `json:"reviewedAt,omitempty" example:"2025-01-16T14:20:00Z"`
} // @name ExpenseRequest

// @Description Budget model
type SwaggerBudget struct {
	ID        uint    `json:"id" example:"1"`
	Year      int     `json:"year" example:"2025"`
	Month     int     `json:"month" example:"1"`
	Total     float64 `json:"total" example:"100000"`
	Spent     float64 `json:"spent" example:"28500"`
	Remaining float64 `json:"remaining" example:"71500"`
} // @name Budget

// @Description Statistics response
type SwaggerStatsResponse struct {
	TotalPending      float64 `json:"totalPending" example:"45000"`
	TotalApproved     float64 `json:"totalApproved" example:"28500"`
	PendingCount      int     `json:"pendingCount" example:"1"`
	ApprovedThisMonth int     `json:"approvedThisMonth" example:"3"`
	BudgetUsed        float64 `json:"budgetUsed" example:"28500"`
	BudgetRemaining   float64 `json:"budgetRemaining" example:"71500"`
} // @name StatsResponse

// @Description Create expense request DTO
type SwaggerCreateExpenseRequestDTO struct {
	Title       string  `json:"title" example:"Закупка офисной мебели" binding:"required,min=3"`
	Category    string  `json:"category" example:"furniture" binding:"required"`
	Amount      float64 `json:"amount" example:"45000" binding:"required,gt=0"`
	Vendor      string  `json:"vendor" example:"IKEA" binding:"required,min=2"`
	Description string  `json:"description" example:"Необходимо приобрести 3 рабочих стола" binding:"required,min=10"`
} // @name CreateExpenseRequestDTO

// @Description Update expense status DTO
type SwaggerUpdateExpenseStatusDTO struct {
	Status   string `json:"status" example:"approved" binding:"required,oneof=approved rejected"`
	Comments string `json:"comments" example:"Одобрено. Необходимо для работы команды." binding:"required,min=10"`
} // @name UpdateExpenseStatusDTO

// @Description Register user DTO
type SwaggerRegisterUserDTO struct {
	Email     string `json:"email" example:"user@example.com" binding:"required,email"`
	Password  string `json:"password" example:"password123" binding:"required,min=6"`
	FirstName string `json:"firstName" example:"Иван" binding:"required"`
	LastName  string `json:"lastName" example:"Иванов" binding:"required"`
	Role      string `json:"role" example:"employee" binding:"required,oneof=employee management"`
} // @name RegisterUserDTO

// @Description Login DTO
type SwaggerLoginDTO struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Password string `json:"password" example:"password123" binding:"required"`
} // @name LoginDTO

// @Description Login response
type SwaggerLoginResponse struct {
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  SwaggerUser `json:"user"`
} // @name LoginResponse
