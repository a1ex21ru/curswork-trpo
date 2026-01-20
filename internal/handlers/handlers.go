package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"curswork-trpo/internal/models"
	"curswork-trpo/internal/service"

	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	expenseService *service.ExpenseService
	userService    *service.UserService
	budgetService  *service.BudgetService
}

func NewExpenseHandler(
	expenseService *service.ExpenseService,
	userService *service.UserService,
	budgetService *service.BudgetService,
) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
		userService:    userService,
		budgetService:  budgetService,
	}
}

// CreateExpenseRequest godoc
// @Summary Create expense request
// @Description Create a new expense request
// @Tags expenses
// @Accept json
// @Produce json
// @Param request body models.CreateExpenseRequestDTO true "Expense request data"
// @Success 201 {object} models.ExpenseRequest
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/expenses [post]
// @Security BearerAuth
func (h *ExpenseHandler) CreateExpenseRequest(c *gin.Context) {
	var dto models.CreateExpenseRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(
			http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			},
		)
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(
			http.StatusUnauthorized, ErrorResponse{
				Error: "unauthorized",
			},
		)
		return
	}

	request, err := h.expenseService.CreateExpenseRequest(c.Request.Context(), &dto, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx := context.Background()
	// TODO Запрос к базе для получения данных юзера
	// т.к. объект employee возвращается пустым
	employ, err := h.userService.GetUserByID(ctx, userID)

	request.Employee = *employ

	c.JSON(http.StatusCreated, request)
}

// GetExpenseRequests godoc
// @Summary Get expense requests
// @Description Get expense requests for current user or all (for management)
// @Tags expenses
// @Produce json
// @Param status query string false "Status filter (all, pending, approved, rejected)"
// @Success 200 {array} models.ExpenseRequest
// @Failure 401 {object} ErrorResponse
// @Router /api/expenses [get]
// @Security BearerAuth
func (h *ExpenseHandler) GetExpenseRequests(c *gin.Context) {
	userID := c.GetUint("userID")
	userRole := c.GetString("userRole")
	status := c.DefaultQuery("status", "all")

	var requests []models.ExpenseRequest
	var err error

	if userRole == string(models.RoleManagement) {
		requests, err = h.expenseService.GetAllExpenseRequests(c.Request.Context(), status)
	} else {
		requests, err = h.expenseService.GetExpenseRequestsByEmployee(c.Request.Context(), userID, status)
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, requests)
}

// GetExpenseRequest godoc
// @Summary Get expense request by ID
// @Description Get a specific expense request
// @Tags expenses
// @Produce json
// @Param id path int true "Expense request ID"
// @Success 200 {object} models.ExpenseRequest
// @Failure 404 {object} ErrorResponse
// @Router /api/expenses/{id} [get]
// @Security BearerAuth
func (h *ExpenseHandler) GetExpenseRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, ErrorResponse{
				Error: "invalid id",
			},
		)
		return
	}

	request, err := h.expenseService.GetExpenseRequest(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(
			http.StatusNotFound, ErrorResponse{
				Error: "request not found",
			},
		)
		return
	}

	c.JSON(http.StatusOK, request)
}

// UpdateExpenseRequestStatus godoc
// @Summary Update expense request status
// @Description Approve or reject an expense request (management only)
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path int true "Expense request ID"
// @Param request body models.UpdateExpenseStatusDTO true "Status update data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/expenses/{id}/status [put]
// @Security BearerAuth
func (h *ExpenseHandler) UpdateExpenseRequestStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, ErrorResponse{
				Error: "invalid id",
			},
		)
		return
	}
	fmt.Println(id)

	var dto models.UpdateExpenseStatusDTO
	if err = c.ShouldBindJSON(&dto); err != nil {
		c.JSON(
			http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			},
		)
		return
	}
	fmt.Println(dto)

	reviewerID := c.GetUint("userID")
	fmt.Println(reviewerID)

	if dto.Status == models.StatusApproved {
		err = h.expenseService.ApproveExpenseRequest(c.Request.Context(), uint(id), reviewerID, dto.Comments)
	} else {
		err = h.expenseService.RejectExpenseRequest(c.Request.Context(), uint(id), reviewerID, dto.Comments)
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "status updated successfully"})
}

// GetStatistics godoc
// @Summary Get expense statistics
// @Description Get expense statistics (management only)
// @Tags expenses
// @Produce json
// @Success 200 {object} models.StatsResponse
// @Router /api/expenses/statistics [get]
// @Security BearerAuth
func (h *ExpenseHandler) GetStatistics(c *gin.Context) {
	stats, err := h.expenseService.GetStatistics(c.Request.Context())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *ExpenseHandler) GetTopExpenses(c *gin.Context) {
	expenses, err := h.expenseService.GetTopExpenses(c.Request.Context())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			},
		)
		return
	}

	categories, err := h.expenseService.GetExpensesByCategory(c.Request.Context())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			},
		)
		return
	}

	report := models.Report{
		Category:   categories,
		TopExpense: expenses,
	}

	c.JSON(
		http.StatusOK, report,
	)
}

// BudgetHandler handles budget operations
type BudgetHandler struct {
	budgetService *service.BudgetService
}

func NewBudgetHandler(budgetService *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
	}
}

// GetCurrentBudget godoc
// @Summary Get current budget
// @Description Get budget for current month
// @Tags budget
// @Produce json
// @Success 200 {object} models.Budget
// @Router /api/budget/current [get]
// @Security BearerAuth
func (h *BudgetHandler) GetCurrentBudget(c *gin.Context) {
	budget, err := h.budgetService.GetCurrentBudget(c.Request.Context())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, budget)
}

// ErrorResponse Response types
type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
