package repository

import (
	"fmt"
	"time"

	"curswork-trpo/internal/models"
	"curswork-trpo/pkg/adapters/postgres"
)

import (
	"context"
)

type ExpenseRepository struct {
	client *postgres.Client
}

func NewExpenseRepository(client *postgres.Client) *ExpenseRepository {
	return &ExpenseRepository{client: client}
}

// CreateExpenseRequest creates a new expense request
func (r *ExpenseRepository) CreateExpenseRequest(ctx context.Context, req *models.ExpenseRequest) error {
	query := `
		INSERT INTO expense_requests (title, category, amount, vendor, description, status, employee_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	now := time.Now().UTC()
	return r.client.QueryRow(
		ctx, query,
		req.Title, req.Category, req.Amount, req.Vendor, req.Description,
		models.StatusPending, req.EmployeeID, now, now,
	).Scan(&req.ID)
}

// GetExpenseRequestByID gets an expense request by ID
func (r *ExpenseRepository) GetExpenseRequestByID(ctx context.Context, id uint) (*models.ExpenseRequest, error) {
	query := `
		SELECT er.id, er.title, er.category, er.amount, er.vendor, er.description, 
		       er.status, er.employee_id, er.reviewer_id, er.comments, 
		       er.created_at, er.updated_at, er.reviewed_at,
		       e.id, e.email, e.first_name, e.last_name, e.role,
		       r.id, r.email, r.first_name, r.last_name, r.role
		FROM expense_requests er
		LEFT JOIN users e ON er.employee_id = e.id
		LEFT JOIN users r ON er.reviewer_id = r.id
		WHERE er.id = $1
	`

	var req models.ExpenseRequest
	var employee models.User
	var reviewerID *uint
	var comments *string
	var reviewedAt *time.Time

	// Nullable fields for reviewer
	var reviewerIDNullable *uint
	var reviewerEmail *string
	var reviewerFirstName *string
	var reviewerLastName *string
	var reviewerRole *string

	err := r.client.QueryRow(ctx, query, id).Scan(
		&req.ID, &req.Title, &req.Category, &req.Amount, &req.Vendor,
		&req.Description, &req.Status, &req.EmployeeID, &reviewerID, &comments,
		&req.CreatedAt, &req.UpdatedAt, &reviewedAt,
		&employee.ID, &employee.Email, &employee.FirstName, &employee.LastName, &employee.Role,
		&reviewerIDNullable, &reviewerEmail, &reviewerFirstName, &reviewerLastName, &reviewerRole,
	)

	if err != nil {
		return nil, err
	}

	req.Employee = employee
	req.ReviewerID = reviewerID
	req.ReviewedAt = reviewedAt

	if comments != nil {
		req.Comments = *comments
	}

	// Build reviewer object if exists
	if reviewerIDNullable != nil && reviewerEmail != nil {
		req.Reviewer = &models.User{
			ID:        *reviewerIDNullable,
			Email:     *reviewerEmail,
			FirstName: *reviewerFirstName,
			LastName:  *reviewerLastName,
			Role:      models.UserRole(*reviewerRole),
		}
	}

	return &req, nil
}

// GetExpenseRequestsByEmployee gets all expense requests for an employee
func (r *ExpenseRepository) GetExpenseRequestsByEmployee(ctx context.Context, employeeID uint, status string) ([]models.ExpenseRequest, error) {
	query := `
		SELECT er.id, er.title, er.category, er.amount, er.vendor, er.description, 
		       er.status, er.employee_id, er.reviewer_id, er.comments, 
		       er.created_at, er.updated_at, er.reviewed_at,
		       e.id, e.email, e.first_name, e.last_name, e.role,
		       r.id, r.email, r.first_name, r.last_name, r.role
		FROM expense_requests er
		LEFT JOIN users e ON er.employee_id = e.id
		LEFT JOIN users r ON er.reviewer_id = r.id
		WHERE er.employee_id = $1
	`

	args := []interface{}{employeeID}

	if status != "" && status != "all" {
		query += " AND er.status = $2"
		args = append(args, status)
	}

	query += " ORDER BY er.created_at DESC"

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.ExpenseRequest
	for rows.Next() {
		var req models.ExpenseRequest
		var employee models.User
		var reviewerID *uint
		var comments *string
		var reviewedAt *time.Time

		// Nullable fields for reviewer
		var reviewerIDNullable *uint
		var reviewerEmail *string
		var reviewerFirstName *string
		var reviewerLastName *string
		var reviewerRole *string

		err = rows.Scan(
			&req.ID, &req.Title, &req.Category, &req.Amount, &req.Vendor,
			&req.Description, &req.Status, &req.EmployeeID, &reviewerID, &comments,
			&req.CreatedAt, &req.UpdatedAt, &reviewedAt,
			&employee.ID, &employee.Email, &employee.FirstName, &employee.LastName, &employee.Role,
			&reviewerIDNullable, &reviewerEmail, &reviewerFirstName, &reviewerLastName, &reviewerRole,
		)
		if err != nil {
			return nil, err
		}

		req.Employee = employee
		req.ReviewerID = reviewerID
		req.ReviewedAt = reviewedAt

		if comments != nil {
			req.Comments = *comments
		}

		if reviewerIDNullable != nil && reviewerEmail != nil {
			req.Reviewer = &models.User{
				ID:        *reviewerIDNullable,
				Email:     *reviewerEmail,
				FirstName: *reviewerFirstName,
				LastName:  *reviewerLastName,
				Role:      models.UserRole(*reviewerRole),
			}
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// GetAllExpenseRequests gets all expense requests with optional filter
func (r *ExpenseRepository) GetAllExpenseRequests(ctx context.Context, status string) ([]models.ExpenseRequest, error) {
	query := `
		SELECT er.id, er.title, er.category, er.amount, er.vendor, er.description, 
		       er.status, er.employee_id, er.reviewer_id, er.comments, 
		       er.created_at, er.updated_at, er.reviewed_at,
		       e.id, e.email, e.first_name, e.last_name, e.role,
		       r.id, r.email, r.first_name, r.last_name, r.role
		FROM expense_requests er
		LEFT JOIN users e ON er.employee_id = e.id
		LEFT JOIN users r ON er.reviewer_id = r.id
	`

	var args []interface{}

	if status != "" && status != "all" {
		query += " WHERE er.status = $1"
		args = append(args, status)
	}

	query += " ORDER BY er.created_at DESC"

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.ExpenseRequest
	for rows.Next() {
		var req models.ExpenseRequest
		var employee models.User
		var reviewerID *uint
		var comments *string
		var reviewedAt *time.Time

		// Nullable fields for reviewer
		var reviewerIDNullable *uint
		var reviewerEmail *string
		var reviewerFirstName *string
		var reviewerLastName *string
		var reviewerRole *string

		err = rows.Scan(
			&req.ID, &req.Title, &req.Category, &req.Amount, &req.Vendor,
			&req.Description, &req.Status, &req.EmployeeID, &reviewerID, &comments,
			&req.CreatedAt, &req.UpdatedAt, &reviewedAt,
			&employee.ID, &employee.Email, &employee.FirstName, &employee.LastName, &employee.Role,
			&reviewerIDNullable, &reviewerEmail, &reviewerFirstName, &reviewerLastName, &reviewerRole,
		)
		if err != nil {
			return nil, err
		}

		req.Employee = employee
		req.ReviewerID = reviewerID
		req.ReviewedAt = reviewedAt

		if comments != nil {
			req.Comments = *comments
		}

		if reviewerIDNullable != nil && reviewerEmail != nil {
			req.Reviewer = &models.User{
				ID:        *reviewerIDNullable,
				Email:     *reviewerEmail,
				FirstName: *reviewerFirstName,
				LastName:  *reviewerLastName,
				Role:      models.UserRole(*reviewerRole),
			}
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// UpdateExpenseRequestStatus updates the status of an expense request
func (r *ExpenseRepository) UpdateExpenseRequestStatus(ctx context.Context, id uint, reviewerID uint, status models.RequestStatus, comments string) error {
	query := `
		UPDATE expense_requests 
		SET status = $1, reviewer_id = $2, comments = $3, reviewed_at = $4, updated_at = $5
		WHERE id = $6
	`
	now := time.Now().UTC()
	_, err := r.client.Exec(ctx, query, status, reviewerID, comments, now, now, id)
	return err
}

// GetStatistics gets expense statistics
func (r *ExpenseRepository) GetStatistics(ctx context.Context) (*models.StatsResponse, error) {
	var stats models.StatsResponse

	// Total pending
	query := `SELECT COALESCE(SUM(amount), 0) FROM expense_requests WHERE status = $1`
	err := r.client.QueryRow(ctx, query, models.StatusPending).Scan(&stats.TotalPending)
	if err != nil {
		return nil, err
	}

	// Pending count
	query = `SELECT COUNT(*) FROM expense_requests WHERE status = $1`
	var count int64
	err = r.client.QueryRow(ctx, query, models.StatusPending).Scan(&count)
	if err != nil {
		return nil, err
	}
	stats.PendingCount = int(count)

	// Total approved this month
	now := time.Now().UTC()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	query = `SELECT COALESCE(SUM(amount), 0) FROM expense_requests WHERE status = $1 AND created_at >= $2`
	err = r.client.QueryRow(ctx, query, models.StatusApproved, startOfMonth).Scan(&stats.TotalApproved)
	if err != nil {
		return nil, err
	}

	query = `SELECT COUNT(*) FROM expense_requests WHERE status = $1 AND created_at >= $2`
	err = r.client.QueryRow(ctx, query, models.StatusApproved, startOfMonth).Scan(&count)
	if err != nil {
		return nil, err
	}
	stats.ApprovedThisMonth = int(count)

	// Budget info
	query = `SELECT spent, remaining FROM budgets WHERE year = $1 AND month = $2`
	err = r.client.QueryRow(ctx, query, now.Year(), int(now.Month())).Scan(&stats.BudgetUsed, &stats.BudgetRemaining)
	if err != nil {
		stats.BudgetUsed = 0
		stats.BudgetRemaining = 0
	}

	return &stats, nil
}

func (r *ExpenseRepository) GetTopExpenses(ctx context.Context) ([]models.TopExpenseRequest, error) {
	var expenses []models.TopExpenseRequest

	query := `
		SELECT title, amount, category, status
		FROM expense_requests
		ORDER BY amount DESC 
		LIMIT 3;`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetTopExpenses: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var exp models.TopExpenseRequest
		if err = rows.Scan(&exp.Title, &exp.Amount, &exp.Category, &exp.Status); err != nil {
			return nil, fmt.Errorf("GetTopExpenses scan: %w", err)
		}
		expenses = append(expenses, exp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTopExpenses rows: %w", err)
	}

	return expenses, nil
}

func (r *ExpenseRepository) GetExpensesByCategory(ctx context.Context) ([]models.CategoryExpenseRequest, error) {
	var expenses []models.CategoryExpenseRequest

	query := `
		SELECT category, SUM(amount) AS amount, COUNT(*) AS count
		FROM expense_requests
		GROUP BY category;`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetExpensesByCategory: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var exp models.CategoryExpenseRequest
		if err = rows.Scan(&exp.Category, &exp.TotalAmount, &exp.Count); err != nil {
			return nil, fmt.Errorf("GetExpensesByCategory scan: %w", err)
		}
		expenses = append(expenses, exp)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetExpensesByCategory rows: %w", err)
	}
	return expenses, nil
}

// UserRepository handles user operations
type UserRepository struct {
	client *postgres.Client
}

func NewUserRepository(client *postgres.Client) *UserRepository {
	return &UserRepository{client: client}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password, first_name, last_name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	now := time.Now().UTC()
	return r.client.QueryRow(
		ctx, query,
		user.Email, user.Password, user.FirstName, user.LastName, user.Role, now, now,
	).Scan(&user.ID)
}

// GetUserByEmail gets a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.client.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID gets a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.client.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// BudgetRepository handles budget operations
type BudgetRepository struct {
	client *postgres.Client
}

func NewBudgetRepository(client *postgres.Client) *BudgetRepository {
	return &BudgetRepository{client: client}
}

// GetOrCreateCurrentBudget gets or creates budget for current month
func (r *BudgetRepository) GetOrCreateCurrentBudget(ctx context.Context) (*models.Budget, error) {
	now := time.Now().UTC()
	year, month := now.Year(), int(now.Month())

	query := `SELECT id, year, month, total, spent, remaining, created_at, updated_at FROM budgets WHERE year = $1 AND month = $2`

	var budget models.Budget
	err := r.client.QueryRow(ctx, query, year, month).Scan(
		&budget.ID, &budget.Year, &budget.Month, &budget.Total,
		&budget.Spent, &budget.Remaining, &budget.CreatedAt, &budget.UpdatedAt,
	)

	if err != nil {
		// Create new budget if not found
		insertQuery := `
			INSERT INTO budgets (year, month, total, spent, remaining, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, year, month, total, spent, remaining, created_at, updated_at
		`
		now = time.Now().UTC()
		err = r.client.QueryRow(ctx, insertQuery, year, month, 100000.0, 0.0, 100000.0, now, now).Scan(
			&budget.ID, &budget.Year, &budget.Month, &budget.Total,
			&budget.Spent, &budget.Remaining, &budget.CreatedAt, &budget.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return &budget, nil
}

// UpdateBudgetSpent updates the spent amount in budget
func (r *BudgetRepository) UpdateBudgetSpent(ctx context.Context, year, month int, amount float64) error {
	query := `
		UPDATE budgets 
		SET spent = spent + $1, remaining = remaining - $2, updated_at = $3
		WHERE year = $4 AND month = $5
	`

	_, err := r.client.Exec(ctx, query, amount, amount, time.Now().UTC(), year, month)
	return err
}

// GetBudgetByMonth gets budget for specific month
func (r *BudgetRepository) GetBudgetByMonth(ctx context.Context, year, month int) (*models.Budget, error) {
	query := `SELECT id, year, month, total, spent, remaining, created_at, updated_at 
				FROM budgets 
				WHERE year = $1 AND month = $2;`

	var budget models.Budget
	err := r.client.QueryRow(ctx, query, year, month).Scan(
		&budget.ID, &budget.Year, &budget.Month, &budget.Total,
		&budget.Spent, &budget.Remaining, &budget.CreatedAt, &budget.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("budget not found for %d-%d", year, month)
	}
	return &budget, nil
}
