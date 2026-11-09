package gateway

import (
	"gourd/internal/config"
	"gourd/internal/constants"
	"time"

	"gourd/internal/counter"
	"gourd/internal/document"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type Service struct {
	app             *fiber.App
	documentService *document.Service
	counterService  *counter.Service
	config          config.Config
}

func New(documentService *document.Service, counterService *counter.Service, config config.Config) *Service {
	fiberApp := fiber.New(fiber.Config{
		BodyLimit: constants.MaxTotalSize,
	})

	s := &Service{
		app:             fiberApp,
		documentService: documentService,
		counterService:  counterService,
		config:          config,
	}

	s.setupRoutes()
	return s
}

func (s *Service) setupRoutes() {
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
	}))

	if !s.config.DevMode {
		s.app.Use(limiter.New(limiter.Config{
			Max:        100,
			Expiration: time.Minute,
			LimitReached: func(c fiber.Ctx) error {
				return c.SendStatus(fiber.StatusTooManyRequests)
			},
		}))
	}

	s.app.Get("/health", s.Health)
	s.app.Get("/rules", s.ListRules)
	s.app.Get("/counter", s.GetCounterValue)

	heavyLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	})

	docs := s.app.Group("/documents")
	if !s.config.DevMode {
		docs.Use(heavyLimiter)
	}
	docs.Post("/file", s.ProcessFile)
	docs.Post("/raw", s.ProcessRaw)
}

func (s *Service) Start(addr string) error {
	return s.app.Listen(addr)
}

func (s *Service) Shutdown(timeout time.Duration) error {
	return s.app.ShutdownWithTimeout(timeout)
}

// Health checks the service status
//
//	@Summary Health check
//	@Description Returns a simple message confirming that the service is up and running
//	@Tags service
//	@Produce plain
//	@Success 200 {string} string "Healthy :)"
//	@Router /health [get]
func (s *Service) Health(c fiber.Ctx) error {
	return c.SendString("Healthy :)")
}
