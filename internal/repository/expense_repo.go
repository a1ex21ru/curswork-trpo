package repository

//
//import (
//	"context"
//	"time"
//
//	"curswork-trpo/internal/models"
//	"curswork-trpo/pkg/adapters/postgres"
//)
//
//type ExpenseRepository struct {
//	client *postgres.Client
//}
//
//func NewExpenseRepository(client *postgres.Client) *ExpenseRepository {
//	return &ExpenseRepository{client: client}
//}
//
//// CreateExpenseRequest creates a new expense request
//func (r *ExpenseRepository) CreateExpenseRequest(ctx context.Context, req *models.ExpenseRequest) error {
//	query := `
//		INSERT INTO expense_requests (title, category, amount, vendor, description, status, employee_id, created_at, updated_at)
//		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
//		RETURNING id
//	`
//	now := time.Now()
//	return r.client.QueryRow(
//		ctx, query,
//		req.Title, req.Category, req.Amount, req.Vendor, req.Description,
//		models.StatusPending, req.EmployeeID, now, now,
//	).Scan(&req.ID)
//}
//
//// GetExpenseRequestByID gets an expense request by ID
//func (r *ExpenseRepository) GetExpenseRequestByID(ctx context.Context, id uint) (*models.ExpenseRequest, error) {
//	query := `
//		SELECT er.id, er.title, er.category, er.amount, er.vendor, er.description,
//		       er.status, er.employee_id, er.reviewer_id, er.comments,
//		       er.created_at, er.updated_at, er.reviewed_at,
//		       e.id, e.email, e.first_name, e.last_name, e.role,
//		       COALESCE(r.id, 0), COALESCE(r.email, ''), COALESCE(r.first_name, ''),
//		       COALESCE(r.last_name, ''), COALESCE(r.role, '')
//		FROM expense_requests er
//		LEFT JOIN users e ON er.employee_id = e.id
//		LEFT JOIN users r ON er.reviewer_id = r.id
//		WHERE er.id = $1
//	`
//
//	var req models.ExpenseRequest
//	var employee models.User
//	var reviewer models.User
//	var reviewerID *uint
//	var reviewedAt *time.Time
//
//	err := r.client.QueryRow(ctx, query, id).Scan(
//		&req.ID, &req.Title, &req.Category, &req.Amount, &req.Vendor,
//		&req.Description, &req.Status, &req.EmployeeID, &reviewerID, &req.Comments,
//		&req.CreatedAt, &req.UpdatedAt, &reviewedAt,
//		&employee.ID, &employee.Email, &employee.FirstName, &employee.LastName, &employee.Role,
//		&reviewer.ID, &reviewer.Email, &reviewer.FirstName, &reviewer.LastName, &reviewer.Role,
//	)
//
//	if err != nil {
//		return nil, err
//	}
//
//	req.Employee = employee
//	req.ReviewerID = reviewerID
//	req.ReviewedAt = reviewedAt
//
//	if reviewerID != nil && reviewer.ID != 0 {
//		req.Reviewer = &reviewer
//	}
//
//	return &req, nil
//}
//
//// GetExpenseRequestsByEmployee gets all expense requests for an employee
//func (r *ExpenseRepository) GetExpenseRequestsByEmployee(ctx context.Context, employeeID uint, status string) ([]models.ExpenseRequest, error) {
//	query := `
//		SELECT er.id, er.title, er.category, er.amount, er.vendor, er.description,
//		       er.status, er.employee_id, er.reviewer_id, er.comments,
//		       er.created_at, er.updated_at, er.reviewed_at,
//		       e.id, e.email, e.first_name, e.last_name, e.role,
//		       COALESCE(r.id, 0), COALESCE(r.email, ''), COALESCE(r.first_name, ''),
//		       COALESCE(r.last_name, ''), COALESCE(r.role, '')
//		FROM expense_requests er
//		LEFT JOIN users e ON er.employee_id = e.id
//		LEFT JOIN users r ON er.reviewer_id = r.id
//		WHERE er.employee_id = $1
//	`
//
//	args := []interface{}{employeeID}
//
//	if status != "" && status != "all" {
//		query += " AND er.status = $2"
//		args = append(args, status)
//	}
//
//	query += " ORDER BY er.created_at DESC"
//
//	rows, err := r.client.Query(ctx, query, args...)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var requests []models.ExpenseRequest
//	for rows.Next() {
//		var req models.ExpenseRequest
//		var employee models.User
//		var reviewer models.User
//		var reviewerID *uint
//		var reviewedAt *time.Time
//
//		err = rows.Scan(
//			&req.ID, &req.Title, &req.Category, &req.Amount, &req.Vendor,
//			&req.Description, &req.Status, &req.EmployeeID, &reviewerID, &req.Comments,
//			&req.CreatedAt, &req.UpdatedAt, &reviewedAt,
//			&employee.ID, &employee.Email, &employee.FirstName, &employee.LastName, &employee.Role,
//			&reviewer.ID, &reviewer.Email, &reviewer.FirstName, &reviewer.LastName, &reviewer.Role,
//		)
//		if err != nil {
//			return nil, err
//		}
//
//		req.Employee = employee
//		req.ReviewerID = reviewerID
//		req.ReviewedAt = reviewedAt
//
//		if reviewerID != nil && reviewer.ID != 0 {
//			req.Reviewer = &reviewer
//		}
//
//		requests = append(requests, req)
//	}
//
//	return requests, nil
//}
//
//// GetAllExpenseRequests gets all expense requests with optional filter
//func (r *ExpenseRepository) GetAllExpenseRequests(ctx context.Context, status string) ([]models.ExpenseRequest, error) {
//	query := `
//		SELECT er.id, er.title, er.category, er.amount, er.vendor, er.description,
//		       er.status, er.employee_id, er.reviewer_id, er.comments,
//		       er.created_at, er.updated_at, er.reviewed_at,
//		       e.id, e.email, e.first_name, e.last_name, e.role,
//		       COALESCE(r.id, 0), COALESCE(r.email, ''), COALESCE(r.first_name, ''),
//		       COALESCE(r.last_name, ''), COALESCE(r.role, '')
//		FROM expense_requests er
//		LEFT JOIN users e ON er.employee_id = e.id
//		LEFT JOIN users r ON er.reviewer_id = r.id
//	`
//
//	var args []interface{}
//
//	if status != "" && status != "all" {
//		query += " WHERE er.status = $1"
//		args = append(args, status)
//	}
//
//	query += " ORDER BY er.created_at DESC"
//
//	rows, err := r.client.Query(ctx, query, args...)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var requests []models.ExpenseRequest
//	for rows.Next() {
//		var req models.ExpenseRequest
//		var employee models.User
//		var reviewer models.User
//		var reviewerID *uint
//		var reviewedAt *time.Time
//
//		err := rows.Scan(
//			&req.ID, &req.Title, &req.Category, &req.Amount, &req.Vendor,
//			&req.Description, &req.Status, &req.EmployeeID, &reviewerID, &req.Comments,
//			&req.CreatedAt, &req.UpdatedAt, &reviewedAt,
//			&employee.ID, &employee.Email, &employee.FirstName, &employee.LastName, &employee.Role,
//			&reviewer.ID, &reviewer.Email, &reviewer.FirstName, &reviewer.LastName, &reviewer.Role,
//		)
//		if err != nil {
//			return nil, err
//		}
//
//		req.Employee = employee
//		req.ReviewerID = reviewerID
//		req.ReviewedAt = reviewedAt
//
//		if reviewerID != nil && reviewer.ID != 0 {
//			req.Reviewer = &reviewer
//		}
//
//		requests = append(requests, req)
//	}
//
//	return requests, nil
//}
//
//// UpdateExpenseRequestStatus updates the status of an expense request
//func (r *ExpenseRepository) UpdateExpenseRequestStatus(ctx context.Context, id uint, reviewerID uint, status models.RequestStatus, comments string) error {
//	query := `
//		UPDATE expense_requests
//		SET status = $1, reviewer_id = $2, comments = $3, reviewed_at = $4, updated_at = $5
//		WHERE id = $6
//	`
//	now := time.Now()
//	_, err := r.client.Exec(ctx, query, status, reviewerID, comments, now, now, id)
//	return err
//}
//
//// GetStatistics gets expense statistics
//func (r *ExpenseRepository) GetStatistics(ctx context.Context) (*models.StatsResponse, error) {
//	var stats models.StatsResponse
//
//	// Total pending
//	query := `SELECT COALESCE(SUM(amount), 0) FROM expense_requests WHERE status = $1`
//	err := r.client.QueryRow(ctx, query, models.StatusPending).Scan(&stats.TotalPending)
//	if err != nil {
//		return nil, err
//	}
//
//	// Pending count
//	query = `SELECT COUNT(*) FROM expense_requests WHERE status = $1`
//	var count int64
//	err = r.client.QueryRow(ctx, query, models.StatusPending).Scan(&count)
//	if err != nil {
//		return nil, err
//	}
//	stats.PendingCount = int(count)
//
//	// Total approved this month
//	now := time.Now()
//	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
//
//	query = `SELECT COALESCE(SUM(amount), 0) FROM expense_requests WHERE status = $1 AND created_at >= $2`
//	err = r.client.QueryRow(ctx, query, models.StatusApproved, startOfMonth).Scan(&stats.TotalApproved)
//	if err != nil {
//		return nil, err
//	}
//
//	query = `SELECT COUNT(*) FROM expense_requests WHERE status = $1 AND created_at >= $2`
//	err = r.client.QueryRow(ctx, query, models.StatusApproved, startOfMonth).Scan(&count)
//	if err != nil {
//		return nil, err
//	}
//	stats.ApprovedThisMonth = int(count)
//
//	// Budget info
//	query = `SELECT spent, remaining FROM budgets WHERE year = $1 AND month = $2`
//	err = r.client.QueryRow(ctx, query, now.Year(), int(now.Month())).Scan(&stats.BudgetUsed, &stats.BudgetRemaining)
//	if err != nil {
//		stats.BudgetUsed = 0
//		stats.BudgetRemaining = 0
//	}
//
//	return &stats, nil
//}
