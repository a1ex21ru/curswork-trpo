package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"curswork-trpo/internal/models"
	"curswork-trpo/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type ExpenseService struct {
	expenseRepo *repository.ExpenseRepository
	budgetRepo  *repository.BudgetRepository
	userRepo    *repository.UserRepository
}

func NewExpenseService(
	expenseRepo *repository.ExpenseRepository,
	budgetRepo *repository.BudgetRepository,
	userRepo *repository.UserRepository,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo: expenseRepo,
		budgetRepo:  budgetRepo,
		userRepo:    userRepo,
	}
}

// CreateExpenseRequest creates a new expense request
func (s *ExpenseService) CreateExpenseRequest(ctx context.Context, dto *models.CreateExpenseRequestDTO, employeeID uint) (*models.ExpenseRequest, error) {
	// Validate employee exists
	_, err := s.userRepo.GetUserByID(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	request := &models.ExpenseRequest{
		Title:       dto.Title,
		Category:    dto.Category,
		Amount:      dto.Amount,
		Vendor:      dto.Vendor,
		Description: dto.Description,
		Status:      models.StatusPending,
		EmployeeID:  employeeID,
	}

	if err = s.expenseRepo.CreateExpenseRequest(ctx, request); err != nil {
		return nil, fmt.Errorf("failed to create expense request: %w", err)
	}
	return request, nil
}

// GetExpenseRequest gets an expense request by ID
func (s *ExpenseService) GetExpenseRequest(ctx context.Context, id uint) (*models.ExpenseRequest, error) {
	return s.expenseRepo.GetExpenseRequestByID(ctx, id)
}

// GetExpenseRequestsByEmployee gets expense requests for an employee
func (s *ExpenseService) GetExpenseRequestsByEmployee(ctx context.Context, employeeID uint, status string) ([]models.ExpenseRequest, error) {
	return s.expenseRepo.GetExpenseRequestsByEmployee(ctx, employeeID, status)
}

// GetAllExpenseRequests gets all expense requests
func (s *ExpenseService) GetAllExpenseRequests(ctx context.Context, status string) ([]models.ExpenseRequest, error) {
	return s.expenseRepo.GetAllExpenseRequests(ctx, status)
}

// GetTopExpenses gets top 3 expenses
func (s *ExpenseService) GetTopExpenses(ctx context.Context) ([]models.TopExpenseRequest, error) {
	return s.expenseRepo.GetTopExpenses(ctx)
}
func (s *ExpenseService) GetExpensesByCategory(ctx context.Context) ([]models.CategoryExpenseRequest, error) {
	return s.expenseRepo.GetExpensesByCategory(ctx)
}

// ApproveExpenseRequest approves an expense request
func (s *ExpenseService) ApproveExpenseRequest(ctx context.Context, id uint, reviewerID uint, comments string) error {
	// Get the request
	request, err := s.expenseRepo.GetExpenseRequestByID(ctx, id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}

	if request.Status != models.StatusPending {
		return errors.New("request is not pending")
	}

	// Check remaining for greater zero after expense
	budget, err := s.budgetRepo.GetBudgetByMonth(ctx, 2026, 1)
	if err != nil {
		return fmt.Errorf("budget not found: %w", err)
	}
	if budget.Remaining-request.Amount < 0 {
		return fmt.Errorf("budget remaining %d < 0", request.Amount)
	}

	// Update budget
	now := time.Now()
	if err = s.budgetRepo.UpdateBudgetSpent(ctx, now.Year(), int(now.Month()), request.Amount); err != nil {
		return fmt.Errorf("failed to update budget: %w", err)
	}

	// Update request status
	if err = s.expenseRepo.UpdateExpenseRequestStatus(
		ctx, id, reviewerID, models.StatusApproved, comments,
	); err != nil {
		return fmt.Errorf("failed to approve request: %w", err)
	}

	return nil
}

// RejectExpenseRequest rejects an expense request
func (s *ExpenseService) RejectExpenseRequest(ctx context.Context, id uint, reviewerID uint, comments string) error {
	// Get the request
	request, err := s.expenseRepo.GetExpenseRequestByID(ctx, id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}

	fmt.Println(request.Status)
	if request.Status != models.StatusPending {
		return errors.New("request is not pending")
	}

	// Update request status
	if err = s.expenseRepo.UpdateExpenseRequestStatus(
		ctx, id, reviewerID, models.StatusRejected, comments,
	); err != nil {
		return fmt.Errorf("failed to reject request: %w", err)
	}

	return nil
}

// GetStatistics gets expense statistics
func (s *ExpenseService) GetStatistics(ctx context.Context) (*models.StatsResponse, error) {
	return s.expenseRepo.GetStatistics(ctx)
}

// UserService handles user operations
type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(ctx context.Context, dto *models.RegisterUserDTO) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, dto.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:     dto.Email,
		Password:  string(hashedPassword),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Role:      dto.Role,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// AuthenticateUser authenticates a user
func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// BudgetService handles budget operations
type BudgetService struct {
	budgetRepo *repository.BudgetRepository
}

func NewBudgetService(budgetRepo *repository.BudgetRepository) *BudgetService {
	return &BudgetService{budgetRepo: budgetRepo}
}

// GetCurrentBudget gets or creates current month's budget
func (s *BudgetService) GetCurrentBudget(ctx context.Context) (*models.Budget, error) {
	return s.budgetRepo.GetOrCreateCurrentBudget(ctx)
}

// GetBudgetByMonth gets budget for specific month
func (s *BudgetService) GetBudgetByMonth(ctx context.Context, year, month int) (*models.Budget, error) {
	return s.budgetRepo.GetBudgetByMonth(ctx, year, month)
}
