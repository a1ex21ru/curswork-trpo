package handlers

import (
	"net/http"

	"curswork-trpo/internal/middleware"
	"curswork-trpo/internal/models"
	"curswork-trpo/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication
type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register godoc
// @Summary Register user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterUserDTO true "Registration data"
// @Success 201 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var dto models.RegisterUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	user, err := h.userService.RegisterUser(c.Request.Context(), &dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login user
// @Description Login and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginDTO true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var dto models.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	user, err := h.userService.AuthenticateUser(c.Request.Context(), dto.Email, dto.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "invalid credentials",
		})
		return
	}

	token, err := middleware.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  *user,
	})
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get current authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/me [get]
// @Security BearerAuth
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "authorization get failed",
		})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
