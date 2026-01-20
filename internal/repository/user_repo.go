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
//// UserRepository handles user operations
//type UserRepository struct {
//	client *postgres.Client
//}
//
//func NewUserRepository(client *postgres.Client) *UserRepository {
//	return &UserRepository{client: client}
//}
//
//// CreateUser creates a new user
//func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
//	query := `
//		INSERT INTO users (email, password, first_name, last_name, role, created_at, updated_at)
//		VALUES ($1, $2, $3, $4, $5, $6, $7)
//		RETURNING id
//	`
//	now := time.Now()
//	return r.client.QueryRow(ctx, query,
//		user.Email, user.Password, user.FirstName, user.LastName, user.Role, now, now,
//	).Scan(&user.ID)
//}
//
//// GetUserByEmail gets a user by email
//func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
//	query := `
//		SELECT id, email, password, first_name, last_name, role, created_at, updated_at
//		FROM users
//		WHERE email = $1
//	`
//
//	var user models.User
//	err := r.client.QueryRow(ctx, query, email).Scan(
//		&user.ID, &user.Email, &user.Password, &user.FirstName,
//		&user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
//	)
//
//	if err != nil {
//		return nil, err
//	}
//	return &user, nil
//}
//
//// GetUserByID gets a user by ID
//func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
//	query := `
//		SELECT id, email, password, first_name, last_name, role, created_at, updated_at
//		FROM users
//		WHERE id = $1
//	`
//
//	var user models.User
//	err := r.client.QueryRow(ctx, query, id).Scan(
//		&user.ID, &user.Email, &user.Password, &user.FirstName,
//		&user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
//	)
//
//	if err != nil {
//		return nil, err
//	}
//	return &user, nil
//}
