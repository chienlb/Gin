package handler

import (
	"net/http"
	"strconv"

	"gin-demo/internal/domain"
	"gin-demo/internal/service"
	"gin-demo/pkg/apperror"
	"gin-demo/pkg/logger"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
	log     *logger.Logger
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
		log:     logger.Get(),
	}
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user with provided data
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.CreateUserRequest true "User data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    apperror.CodeBadRequest,
			"message": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	user, appErr := h.service.CreateUser(&req)
	if appErr != nil {
		h.log.Error("Failed to create user", appErr)
		c.JSON(appErr.Status, gin.H{
			"status":  "error",
			"code":    appErr.Code,
			"message": appErr.Message,
			"details": appErr.Details,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"code":    "CREATED",
		"message": "User created successfully",
		"data":    user,
	})
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    apperror.CodeBadRequest,
			"message": "Invalid user ID format",
		})
		return
	}

	user, appErr := h.service.GetUser(id)
	if appErr != nil {
		c.JSON(appErr.Status, gin.H{
			"status":  "error",
			"code":    appErr.Code,
			"message": appErr.Message,
			"details": appErr.Details,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"code":    "OK",
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// GetAllUsers retrieves all users
// @Summary Get all users
// @Description Get list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, appErr := h.service.GetAllUsers()
	if appErr != nil {
		h.log.Error("Failed to get users", appErr)
		c.JSON(appErr.Status, gin.H{
			"status":  "error",
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"code":    "OK",
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

// UpdateUser updates a user by ID
// @Summary Update user by ID
// @Description Update user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body domain.UpdateUserRequest true "User data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    apperror.CodeBadRequest,
			"message": "Invalid user ID format",
		})
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    apperror.CodeBadRequest,
			"message": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	user, appErr := h.service.UpdateUser(id, &req)
	if appErr != nil {
		h.log.Error("Failed to update user", appErr)
		c.JSON(appErr.Status, gin.H{
			"status":  "error",
			"code":    appErr.Code,
			"message": appErr.Message,
			"details": appErr.Details,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"code":    "OK",
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser deletes a user by ID
// @Summary Delete user by ID
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    apperror.CodeBadRequest,
			"message": "Invalid user ID format",
		})
		return
	}

	appErr := h.service.DeleteUser(id)
	if appErr != nil {
		h.log.Error("Failed to delete user", appErr)
		c.JSON(appErr.Status, gin.H{
			"status":  "error",
			"code":    appErr.Code,
			"message": appErr.Message,
			"details": appErr.Details,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"code":    "OK",
		"message": "User deleted successfully",
	})
}
