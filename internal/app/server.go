package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-demo/internal/config"
	"gin-demo/internal/database"
	"gin-demo/internal/handler"
	"gin-demo/internal/repository"
	"gin-demo/internal/service"
	"gin-demo/pkg/logger"
	"gin-demo/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     *config.Config
	engine     *gin.Engine
	log        *logger.Logger
	httpServer *http.Server
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		engine: gin.Default(),
		log:    logger.Get(),
	}
}

func (s *Server) Initialize() error {
	// Initialize database
	if err := database.Init(s.config.Database.GetDSN()); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	s.log.Info("Database initialized successfully")

	// Run migrations
	db := database.GetDB()
	if err := database.RunMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Setup middleware
	s.setupMiddleware()

	// Setup routes
	s.setupRoutes()

	return nil
}

func (s *Server) setupMiddleware() {
	// Add middleware in order
	s.engine.Use(middleware.RequestIDMiddleware())
	s.engine.Use(middleware.LoggingMiddleware(s.log))
	s.engine.Use(middleware.CORSMiddleware())
	s.engine.Use(gin.Recovery())

	s.log.Info("Middleware initialized successfully")
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "OK",
			"timestamp": time.Now(),
		})
	})

	// Root endpoint
	s.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "Gin Demo API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Initialize dependencies
	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Setup API routes
	api := s.engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			users := v1.Group("/users")
			{
				users.POST("", userHandler.CreateUser)
				users.GET("", userHandler.GetAllUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	s.log.Info("Routes initialized successfully")
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	s.log.Info(fmt.Sprintf("Starting server at %s", address))

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         address,
		Handler:      s.engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("Server error", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	sig := <-sigChan
	s.log.Info("Received signal: " + sig.String())

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Error("Server shutdown error", err)
		return err
	}

	s.log.Info("Server shutdown gracefully")
	return nil
}

func (s *Server) Close() error {
	return database.Close()
}
