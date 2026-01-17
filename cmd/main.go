// @title RCS Onboarding API
// @version 1.0
// @description Backend service for RCS customer onboarding workflow
// @termsOfService http://example.com/terms/

// @contact.name RCS API Support
// @contact.email support@rcs.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8084
// @BasePath /api/v1
package main

import (
	"rcs-onboarding/internal/config"
	"rcs-onboarding/internal/handlers"
	"rcs-onboarding/internal/middleware"
	"rcs-onboarding/internal/models"
	"rcs-onboarding/internal/repositories"
	"rcs-onboarding/internal/services"
	"rcs-onboarding/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	utils.InitLogger()
	cfg := config.LoadConfig()

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	if err := db.AutoMigrate(&models.User{}, &models.FormVersion{}, &models.Submission{}, &models.AuditLog{}); err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}

	utils.SeedTemplates(db)
	utils.SeedUsers(db)

	userRepo := repositories.NewUserRepo(db)
	formRepo := repositories.NewFormRepo(db)
	submissionRepo := repositories.NewSubmissionRepo(db)
	auditRepo := repositories.NewAuditRepo(db)

	authService := services.NewAuthService(userRepo)
	formService := services.NewFormService(formRepo)
	submissionService := services.NewSubmissionService(submissionRepo, formRepo, auditRepo)
	auditService := services.NewAuditService(auditRepo)

	authHandler := handlers.NewAuthHandler(authService)
	formHandler := handlers.NewFormHandler(formService)
	submissionHandler := handlers.NewSubmissionHandler(submissionService, auditService)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)

		forms := api.Group("/forms")
		forms.Use(middleware.AuthMiddleware())
		{
			forms.POST("/:type", middleware.RoleMiddleware(models.Admin), formHandler.Create)
			forms.GET("/:type/versions", formHandler.ListVersions)
			forms.GET("/:type/versions/latest", formHandler.GetLatest)
		}

		submissions := api.Group("/submissions")
		submissions.Use(middleware.AuthMiddleware())
		{
			// Register specific route first (longer path)
			submissions.POST("/:id/review", middleware.RoleMiddleware(models.TPM, models.Sales), submissionHandler.Review)

			// Then the general wildcard route
			submissions.POST("/:id", middleware.RoleMiddleware(models.Customer), submissionHandler.Submit)

			submissions.GET("", submissionHandler.GetFiltered)
			submissions.GET("/:id", submissionHandler.GetByID)
			submissions.PUT("/:id", middleware.RoleMiddleware(models.Customer), submissionHandler.UpdateDraft)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Info().Msg("Server starting on :8083")
	if err := r.Run(":8083"); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
