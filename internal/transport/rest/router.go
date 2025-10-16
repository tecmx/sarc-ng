package rest

import (
	"sarc-ng/internal/domain/auth"
	"sarc-ng/internal/domain/building"
	"sarc-ng/internal/domain/class"
	"sarc-ng/internal/domain/lesson"
	"sarc-ng/internal/domain/reservation"
	"sarc-ng/internal/domain/resource"
	buildingRest "sarc-ng/internal/transport/rest/building"
	classRest "sarc-ng/internal/transport/rest/class"
	lessonRest "sarc-ng/internal/transport/rest/lesson"
	reservationRest "sarc-ng/internal/transport/rest/reservation"
	resourceRest "sarc-ng/internal/transport/rest/resource"
	"sarc-ng/pkg/rest/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "sarc-ng/api/swagger" // Import generated docs
)

// Router contains all the dependencies for setting up routes
type Router struct {
	buildingService    building.Usecase
	classService       class.Usecase
	lessonService      lesson.Usecase
	reservationService reservation.Usecase
	resourceService    resource.Usecase
	tokenValidator     auth.TokenValidator
}

// NewRouter creates a new router with all dependencies
func NewRouter(
	buildingService building.Usecase,
	classService class.Usecase,
	lessonService lesson.Usecase,
	reservationService reservation.Usecase,
	resourceService resource.Usecase,
	tokenValidator auth.TokenValidator,
) *Router {
	return &Router{
		buildingService:    buildingService,
		classService:       classService,
		lessonService:      lessonService,
		reservationService: reservationService,
		resourceService:    resourceService,
		tokenValidator:     tokenValidator,
	}
}

// SetupRoutes configures all the routes for the API
func (r *Router) SetupRoutes(router *gin.Engine) {
	// Apply global middleware
	r.setupMiddleware(router)

	// Setup system routes
	r.setupSystemRoutes(router)

	// Setup API routes
	r.setupAPIRoutes(router)
}

// setupMiddleware applies global middleware to the router
func (r *Router) setupMiddleware(router *gin.Engine) {
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Metrics())
}

// setupSystemRoutes configures system routes like health and swagger
func (r *Router) setupSystemRoutes(router *gin.Engine) {
	// Swagger routes with OAuth2 configuration
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.PersistAuthorization(true),
	))

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "sarc-ng",
		})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// setupAPIRoutes configures API v1 routes with authentication
func (r *Router) setupAPIRoutes(router *gin.Engine) {
	// Public API routes (no authentication required)
	// These are read-only endpoints accessible to everyone
	publicV1 := router.Group("/api/v1")
	{
		buildingRest.RegisterRoutes(publicV1, r.buildingService)
		classRest.RegisterRoutes(publicV1, r.classService)
		lessonRest.RegisterRoutes(publicV1, r.lessonService)
		resourceRest.RegisterRoutes(publicV1, r.resourceService)
	}

	// Protected API routes (authentication required)
	protectedV1 := router.Group("/api/v1")
	protectedV1.Use(middleware.AuthMiddleware(r.tokenValidator))
	{
		reservationRest.RegisterRoutes(protectedV1, r.reservationService)
	}
}
