package handlers

import (
	"curswork-trpo/internal/middleware"
	"curswork-trpo/internal/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all routes
func SetupRouter(
	expenseHandler *ExpenseHandler,
	authHandler *AuthHandler,
	budgetHandler *BudgetHandler,
) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", HealthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
		}

		// Expense routes (protected)
		expenses := api.Group("/expenses")
		expenses.Use(middleware.AuthMiddleware())
		{
			expenses.POST("", expenseHandler.CreateExpenseRequest)
			expenses.GET("", expenseHandler.GetExpenseRequests)
			expenses.GET("/:id", expenseHandler.GetExpenseRequest)

			// Management only routes
			expenses.PUT(
				"/:id/status",
				middleware.RoleMiddleware(models.RoleManagement),
				expenseHandler.UpdateExpenseRequestStatus,
			)
			expenses.GET(
				"/statistics",
				middleware.RoleMiddleware(models.RoleManagement),
				expenseHandler.GetStatistics,
			)
		}

		reports := api.Group("/reports/expenses")
		reports.Use(middleware.AuthMiddleware())
		{
			reports.GET("", expenseHandler.GetTopExpenses)
		}

		// Budget routes (protected)
		budget := api.Group("/budget")
		budget.Use(middleware.AuthMiddleware())
		{
			budget.GET("/current", budgetHandler.GetCurrentBudget)
		}
	}

	return router
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(
		200, gin.H{
			"status":  "ok",
			"service": "expense-system-api",
		},
	)
}
