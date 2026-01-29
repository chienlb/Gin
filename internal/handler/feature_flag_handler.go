package handler

import (
	"net/http"

	"gin-demo/pkg/feature"

	"github.com/gin-gonic/gin"
)

// FeatureFlagHandler handles feature flag endpoints
type FeatureFlagHandler struct {
	manager *feature.FeatureFlagManager
}

// NewFeatureFlagHandler creates a new feature flag handler
func NewFeatureFlagHandler(manager *feature.FeatureFlagManager) *FeatureFlagHandler {
	return &FeatureFlagHandler{manager: manager}
}

// ListFlags returns all feature flags
func (h *FeatureFlagHandler) ListFlags(c *gin.Context) {
	flags := h.manager.ListFlags()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   flags,
	})
}

// GetFlag returns a specific feature flag
func (h *FeatureFlagHandler) GetFlag(c *gin.Context) {
	key := c.Param("key")

	flag, err := h.manager.GetFlag(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"code":    "NOT_FOUND",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   flag,
	})
}

// CheckFlag checks if a feature is enabled for a given context
func (h *FeatureFlagHandler) CheckFlag(c *gin.Context) {
	key := c.Param("key")

	var req struct {
		Context map[string]string `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.Context = make(map[string]string)
	}

	// Add user info from headers/context if available
	if userID := c.GetString("user_id"); userID != "" {
		req.Context["user_id"] = userID
	}

	enabled := h.manager.IsEnabled(key, req.Context)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"key":     key,
			"enabled": enabled,
		},
	})
}

// UpdateFlag updates a feature flag
func (h *FeatureFlagHandler) UpdateFlag(c *gin.Context) {
	key := c.Param("key")

	var req struct {
		Enabled bool `json:"enabled"`
		Rollout int  `json:"rollout"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "INVALID_REQUEST",
			"message": "Invalid request body",
		})
		return
	}

	if err := h.manager.UpdateFlag(key, req.Enabled, req.Rollout); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"code":    "NOT_FOUND",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Feature flag updated",
	})
}

// CreateFlag creates a new feature flag
func (h *FeatureFlagHandler) CreateFlag(c *gin.Context) {
	var flag feature.FeatureFlag

	if err := c.ShouldBindJSON(&flag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "INVALID_REQUEST",
			"message": "Invalid request body",
		})
		return
	}

	h.manager.RegisterFlag(&flag)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Feature flag created",
		"data":    flag,
	})
}
