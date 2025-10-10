package main

//	@title			SARC-NG API
//	@version		1.0
//	@description	SARC Next Generation - Resource Reservation and Management API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1
//	@schemes	http https
//
//	@securityDefinitions.oauth2.accessCode	CognitoOAuth
//	@tokenUrl								https://sarc-ng-auth-prod.auth.us-east-1.amazoncognito.com/oauth2/token
//	@authorizationUrl						https://sarc-ng-auth-prod.auth.us-east-1.amazoncognito.com/oauth2/authorize
//	@scope.openid							OpenID Connect scope
//	@scope.email							Access to user email
//	@scope.profile							Access to user profile
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT token from Cognito (use the access_token from OAuth2 login)

import (
	"fmt"
	"log"
	"os"
	"sarc-ng/internal/domain/building"
	"sarc-ng/internal/domain/class"
	"sarc-ng/internal/domain/lesson"
	"sarc-ng/internal/domain/reservation"
	"sarc-ng/internal/domain/resource"
	"sarc-ng/pkg/metrics"

	_ "sarc-ng/api/swagger" // Import generated API documentation

	"github.com/gin-gonic/gin"
)

// Version information (will be set during build)
var (
	version   = "dev"
	buildDate = "unknown"
	commit    = "none"
)

func main() {
	// Initialize metrics with version information
	metrics.Initialize(version, buildDate, commit)

	// Initialize application with Wire dependency injection
	log.Println("Initializing application with dependency injection...")
	app, err := InitializeApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Get SQL database connection for cleanup
	sqlDB, err := app.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer sqlDB.Close()

	// Migrate all domain tables with error handling
	log.Println("Running database migrations...")
	err = app.DB.AutoMigrate(
		&building.Building{},
		&class.Class{},
		&lesson.Lesson{},
		&reservation.Reservation{},
		&resource.Resource{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// Get mode from environment or use default
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	// Create Gin engine
	r := gin.New()

	// Setup routes using the Wire-injected router
	app.Router.SetupRoutes(r)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
