package app

import "github.com/gin-gonic/gin"

type Router struct {
	engine *gin.Engine
}

func NewRouter(engine *gin.Engine) *Router {
	return &Router{engine: engine}
}

func (r *Router) Setup() {
	// Router setup is handled in server.go setupRoutes method
}
