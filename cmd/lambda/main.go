package main

import (
	"context"
	"log"
	"os"
	"sarc-ng/internal/domain/building"
	"sarc-ng/internal/domain/class"
	"sarc-ng/internal/domain/lesson"
	"sarc-ng/internal/domain/reservation"
	"sarc-ng/internal/domain/resource"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Global variable to hold the Gin Lambda adapter
// This is initialized once and reused across Lambda invocations for better performance
var ginLambda *ginadapter.GinLambda

// init function runs once when the Lambda container starts
// This is where we set up our Gin rest and database connections
func init() {
	log.Println("Initializing Lambda function...")

	// Initialize application with Wire dependency injection
	log.Println("Initializing application with dependency injection...")
	app, err := InitializeApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Test database connection with timeout
	if err := pingDatabase(app.DB); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}

	// Auto-migrate database tables
	// In production, consider running migrations separately to avoid cold start delays
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

	// Get mode from environment or use release mode for Lambda
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode // Default to release mode for Lambda production
	}
	gin.SetMode(mode)

	// Create Gin engine
	r := gin.New()

	// Setup routes using the Wire-injected router
	app.Router.SetupRoutes(r)

	// Create the Gin Lambda adapter
	// This adapter converts API Gateway events to HTTP requests for Gin
	ginLambda = ginadapter.New(r)

	log.Println("Lambda function initialized successfully")
}

// pingDatabase tests the database connection
func pingDatabase(database *gorm.DB) error {
	sqlDB, err := database.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Handler is the Lambda function entry point
// It receives API Gateway proxy requests and returns proxy responses
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Only log requests in debug mode
	if os.Getenv("GIN_MODE") == "debug" {
		log.Printf("Received request: %s %s", request.HTTPMethod, request.Path)
	}

	// Convert API Gateway request to HTTP request and get response
	response, err := ginLambda.ProxyWithContext(ctx, request)
	return response, err
}

// main function starts the Lambda runtime
func main() {
	// Check if we're running in Lambda environment
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Starting Lambda runtime...")
		// Start the Lambda handler
		lambda.Start(Handler)
	} else {
		log.Println("Not running in Lambda environment. Use 'go run cmd/main.go' for local development.")
	}
}
