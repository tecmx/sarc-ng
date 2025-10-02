//go:build wireinject

package main

import (
	"sarc-ng/internal/adapter/db"
	buildingAdapter "sarc-ng/internal/adapter/gorm/building"
	classAdapter "sarc-ng/internal/adapter/gorm/class"
	lessonAdapter "sarc-ng/internal/adapter/gorm/lesson"
	reservationAdapter "sarc-ng/internal/adapter/gorm/reservation"
	resourceAdapter "sarc-ng/internal/adapter/gorm/resource"
	"sarc-ng/internal/config"
	"sarc-ng/internal/domain/building"
	"sarc-ng/internal/domain/class"
	"sarc-ng/internal/domain/lesson"
	"sarc-ng/internal/domain/reservation"
	"sarc-ng/internal/domain/resource"
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

// provideDatabaseConnection provides a database connection using the new structure
func provideDatabaseConnection(config *config.Config) (*gorm.DB, error) {
	dbConfig := db.Config{
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
	return db.Connect(dbConfig)
}

// InitializeApplication initializes the application with all dependencies
func InitializeApplication() (*Application, error) {
	wire.Build(ProviderSet)
	return &Application{}, nil
}
