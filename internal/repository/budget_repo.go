package repository

//
//import (
//	"context"
//	"fmt"
//	"time"
//
//	"curswork-trpo/internal/models"
//	"curswork-trpo/pkg/adapters/postgres"
//)
//
//// BudgetRepository handles budget operations
//type BudgetRepository struct {
//	client *postgres.Client
//}
//
//func NewBudgetRepository(client *postgres.Client) *BudgetRepository {
//	return &BudgetRepository{client: client}
//}
//
//// GetOrCreateCurrentBudget gets or creates budget for current month
//func (r *BudgetRepository) GetOrCreateCurrentBudget(ctx context.Context) (*models.Budget, error) {
//	now := time.Now()
//	year, month := now.Year(), int(now.Month())
//
//	query := `SELECT id, year, month, total, spent, remaining, created_at, updated_at FROM budgets WHERE year = $1 AND month = $2`
//
//	var budget models.Budget
//	err := r.client.QueryRow(ctx, query, year, month).Scan(
//		&budget.ID, &budget.Year, &budget.Month, &budget.Total,
//		&budget.Spent, &budget.Remaining, &budget.CreatedAt, &budget.UpdatedAt,
//	)
//
//	if err != nil {
//		// Create new budget if not found
//		insertQuery := `
//			INSERT INTO budgets (year, month, total, spent, remaining, created_at, updated_at)
//			VALUES ($1, $2, $3, $4, $5, $6, $7)
//			RETURNING id, year, month, total, spent, remaining, created_at, updated_at
//		`
//		now := time.Now()
//		err = r.client.QueryRow(ctx, insertQuery, year, month, 100000.0, 0.0, 100000.0, now, now).Scan(
//			&budget.ID, &budget.Year, &budget.Month, &budget.Total,
//			&budget.Spent, &budget.Remaining, &budget.CreatedAt, &budget.UpdatedAt,
//		)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return &budget, nil
//}
//
//// UpdateBudgetSpent updates the spent amount in budget
//func (r *BudgetRepository) UpdateBudgetSpent(ctx context.Context, year, month int, amount float64) error {
//	query := `
//		UPDATE budgets
//		SET spent = spent + $1, remaining = remaining - $2, updated_at = $3
//		WHERE year = $4 AND month = $5
//	`
//	_, err := r.client.Exec(ctx, query, amount, amount, time.Now(), year, month)
//	return err
//}
//
//// GetBudgetByMonth gets budget for specific month
//func (r *BudgetRepository) GetBudgetByMonth(ctx context.Context, year, month int) (*models.Budget, error) {
//	query := `SELECT id, year, month, total, spent, remaining, created_at, updated_at FROM budgets WHERE year = $1 AND month = $2`
//
//	var budget models.Budget
//	err := r.client.QueryRow(ctx, query, year, month).Scan(
//		&budget.ID, &budget.Year, &budget.Month, &budget.Total,
//		&budget.Spent, &budget.Remaining, &budget.CreatedAt, &budget.UpdatedAt,
//	)
//
//	if err != nil {
//		return nil, fmt.Errorf("budget not found for %d-%d", year, month)
//	}
//	return &budget, nil
//}
