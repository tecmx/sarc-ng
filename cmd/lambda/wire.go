//go:build wireinject

package main

import (
	"context"
	"fmt"
	"os"

	"sarc-ng/internal/adapter/db"
	buildingAdapter "sarc-ng/internal/adapter/gorm/building"
	classAdapter "sarc-ng/internal/adapter/gorm/class"
	lessonAdapter "sarc-ng/internal/adapter/gorm/lesson"
	reservationAdapter "sarc-ng/internal/adapter/gorm/reservation"
	resourceAdapter "sarc-ng/internal/adapter/gorm/resource"
	"sarc-ng/internal/adapter/secrets"
	"sarc-ng/internal/config"
	"sarc-ng/internal/domain/auth"
	"sarc-ng/internal/domain/building"
	"sarc-ng/internal/domain/class"
	"sarc-ng/internal/domain/lesson"
	"sarc-ng/internal/domain/reservation"
	"sarc-ng/internal/domain/resource"
	authService "sarc-ng/internal/service/auth"
	buildingService "sarc-ng/internal/service/building"
	classService "sarc-ng/internal/service/class"
	lessonService "sarc-ng/internal/service/lesson"
	reservationService "sarc-ng/internal/service/reservation"
	resourceService "sarc-ng/internal/service/resource"
	"sarc-ng/internal/transport/rest"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// Application holds all the application dependencies
type Application struct {
	DB                 *gorm.DB
	Config             *config.Config
	Router             *rest.Router
	BuildingService    building.Usecase
	ClassService       class.Usecase
	LessonService      lesson.Usecase
	ResourceService    resource.Usecase
	ReservationService reservation.Usecase
}

// ProviderSet for the application
var ProviderSet = wire.NewSet(
	// Configuration
	config.LoadConfig,

	// Database
	provideDatabaseConnection,

	// Authentication
	provideTokenValidator,
	wire.Bind(new(auth.TokenValidator), new(*authService.JWTValidator)),

	// GORM Adapters - these provide the repository implementations
	buildingAdapter.NewGormAdapter,
	classAdapter.NewGormAdapter,
	lessonAdapter.NewGormAdapter,
	resourceAdapter.NewGormAdapter,
	reservationAdapter.NewGormAdapter,

	// Repository interface bindings
	wire.Bind(new(building.Repository), new(*buildingAdapter.GormAdapter)),
	wire.Bind(new(class.Repository), new(*classAdapter.GormAdapter)),
	wire.Bind(new(lesson.Repository), new(*lessonAdapter.GormAdapter)),
	wire.Bind(new(resource.Repository), new(*resourceAdapter.GormAdapter)),
	wire.Bind(new(reservation.Repository), new(*reservationAdapter.GormAdapter)),

	// Services
	buildingService.NewService,
	classService.NewService,
	lessonService.NewService,
	resourceService.NewService,
	reservationService.NewService,

	// Service interface bindings
	wire.Bind(new(building.Usecase), new(*buildingService.Service)),
	wire.Bind(new(class.Usecase), new(*classService.Service)),
	wire.Bind(new(lesson.Usecase), new(*lessonService.Service)),
	wire.Bind(new(resource.Usecase), new(*resourceService.Service)),
	wire.Bind(new(reservation.Usecase), new(*reservationService.Service)),

	// REST Router
	rest.NewRouter,

	// Application structure
	wire.Struct(new(Application), "*"),
)

// provideDatabaseConnection provides a database connection using Secrets Manager or config
func provideDatabaseConnection(config *config.Config) (*gorm.DB, error) {
	ctx := context.Background()

	// Check if using Secrets Manager (Lambda environment)
	secretArn := os.Getenv("DB_SECRET_ARN")

	var dbConfig db.Config
	if secretArn != "" {
		// Get credentials from Secrets Manager
		creds, err := secrets.GetDatabaseCredentials(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get database credentials from Secrets Manager: %w", err)
		}

		// Build config from secret
		dbConfig = db.Config{
			Host:            creds.Host,
			Port:            parsePort(creds.Port),
			User:            creds.Username,
			Password:        creds.Password,
			Database:        creds.Database,
			MaxOpenConns:    config.Database.MaxOpenConns,
			MaxIdleConns:    config.Database.MaxIdleConns,
			ConnMaxLifetime: config.Database.ConnMaxLifetime,
			ConnMaxIdleTime: config.Database.ConnMaxIdleTime,
		}
	} else {
		// Use config file (local development)
		dbConfig = db.Config{
			Host:            config.Database.Host,
			Port:            config.Database.Port,
			User:            config.Database.User,
			Password:        config.Database.Password,
			Database:        config.Database.Name,
			MaxOpenConns:    config.Database.MaxOpenConns,
			MaxIdleConns:    config.Database.MaxIdleConns,
			ConnMaxLifetime: config.Database.ConnMaxLifetime,
			ConnMaxIdleTime: config.Database.ConnMaxIdleTime,
		}
	}

	return db.Connect(dbConfig)
}

// parsePort converts string port to int
func parsePort(portStr string) int {
	var port int
	fmt.Sscanf(portStr, "%d", &port)
	if port == 0 {
		port = 3306 // Default MySQL port
	}
	return port
}

// provideTokenValidator creates a new JWT token validator
func provideTokenValidator(cfg *config.Config) *authService.JWTValidator {
	return authService.NewJWTValidator(
		cfg.Cognito.Region,
		cfg.Cognito.UserPoolID,
		cfg.Cognito.ClientID,
		cfg.Cognito.JWKSCacheExp,
	)
}

// InitializeApplication initializes the application with all dependencies
func InitializeApplication() (*Application, error) {
	wire.Build(ProviderSet)
	return &Application{}, nil
}
