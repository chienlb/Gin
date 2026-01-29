package app

import (
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	engine *gin.Engine
}

func NewMiddleware(engine *gin.Engine) *Middleware {
	return &Middleware{engine: engine}
}

// Setup adds common middleware to the engine
func (m *Middleware) Setup() {
	// CORS middleware could be added here
	// Logging middleware could be added here
	// Recovery middleware is already included in gin.Default()
}
