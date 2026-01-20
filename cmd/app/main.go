package main

import (
	"context"
	"log"
	"os"

	"curswork-trpo/internal/handlers"
	"curswork-trpo/internal/repository"
	"curswork-trpo/internal/service"
	"curswork-trpo/pkg/adapters/postgres"

	_ "curswork-trpo/docs"
)

// @title Expense System API
// @version 1.0
// @description REST API для системы учета внутриофисных расходов
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	ctx := context.Background()

	// Initialize database client
	dbClient, err := postgres.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbClient.Close(ctx)

	// Initialize database schema
	if err = initSchema(ctx, dbClient); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize repositories
	expenseRepo := repository.NewExpenseRepository(dbClient)
	userRepo := repository.NewUserRepository(dbClient)
	budgetRepo := repository.NewBudgetRepository(dbClient)

	// Initialize services
	expenseService := service.NewExpenseService(expenseRepo, budgetRepo, userRepo)
	userService := service.NewUserService(userRepo)
	budgetService := service.NewBudgetService(budgetRepo)

	// Initialize handlers
	expenseHandler := handlers.NewExpenseHandler(expenseService, userService, budgetService)
	authHandler := handlers.NewAuthHandler(userService)
	budgetHandler := handlers.NewBudgetHandler(budgetService)

	// Setup router
	router := handlers.SetupRouter(expenseHandler, authHandler, budgetHandler)

	// Start server
	port := os.Getenv("PORT")
	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	if err = router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initSchema(ctx context.Context, client *postgres.Client) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		role VARCHAR(20) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS expense_requests (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		category VARCHAR(100) NOT NULL,
		amount DECIMAL(10, 2) NOT NULL,
		vendor VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'pending',
		employee_id INTEGER NOT NULL REFERENCES users(id),
		reviewer_id INTEGER REFERENCES users(id),
		comments TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		reviewed_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS budgets (
		id SERIAL PRIMARY KEY,
		year INTEGER NOT NULL,
		month INTEGER NOT NULL,
		total DECIMAL(12, 2) NOT NULL,
		spent DECIMAL(12, 2) NOT NULL DEFAULT 0,
		remaining DECIMAL(12, 2) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		UNIQUE(year, month)
	);

	CREATE INDEX IF NOT EXISTS idx_expense_requests_employee_id ON expense_requests(employee_id);
	CREATE INDEX IF NOT EXISTS idx_expense_requests_status ON expense_requests(status);
	CREATE INDEX IF NOT EXISTS idx_budgets_year_month ON budgets(year, month);
	`

	_, err := client.Exec(ctx, schema)
	if err != nil {
		return err
	}

	log.Println("Database schema initialized successfully")
	return nil
}
