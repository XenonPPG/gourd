package gateway

import (
	"gourd/internal/config"
	"gourd/internal/constants"
	"gourd/internal/domain/counter"
	"gourd/internal/domain/result"
	"gourd/internal/domain/rule"
	"gourd/internal/domain/rule/adapter"
	"gourd/internal/gateway/dto"
	"gourd/internal/gateway/dto/dtoregistry"
	"time"

	"gourd/internal/domain"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type Service struct {
	app             *fiber.App
	documentService *domain.Service
	counterService  *counter.Service
	config          config.Config
	dtoRegistry     *dtoregistry.Registry
}

func New(documentService *domain.Service, counterService *counter.Service, config config.Config) *Service {
	fiberApp := fiber.New(fiber.Config{
		BodyLimit: constants.MaxTotalSize,
	})

	s := &Service{
		app:             fiberApp,
		documentService: documentService,
		counterService:  counterService,
		config:          config,
		dtoRegistry:     dtoregistry.NewRegistry(),
	}

	s.setupRoutes()
	s.registerDTO()

	return s
}

func (s *Service) registerDTO() {
	ruleToDTO := func(r rule.Rule) dto.Rule {
		availableFormats := make([]string, 0, len(r.Adapters))
		for _, a := range r.Adapters.AvailableFileTypes() {
			availableFormats = append(availableFormats, string(a))
		}

		return dto.Rule{
			ID:               r.ID,
			Name:             r.Name,
			Description:      r.Description,
			EnabledByDefault: r.EnabledByDefault,
			AvailableFormats: availableFormats,
		}
	}

	// rule
	s.dtoRegistry.Register(ruleToDTO)

	// category
	s.dtoRegistry.Register(func(category rule.Category) dto.Category {
		rules := make([]dto.Rule, 0, len(category.Rules))
		for _, r := range category.Rules {
			rules = append(rules, ruleToDTO(r))
		}

		return dto.Category{
			Name:  category.Name,
			Rules: rules,
		}
	})

	// indices result
	s.dtoRegistry.Register(func(ir adapter.IndicesResult) dto.IndicesResult {
		return dto.IndicesResult{
			Indices:    ir.Indices,
			TotalCount: ir.TotalCount,
		}
	})

	// result
	s.dtoRegistry.Register(func(res result.Result) dto.Result {
		d, _ := s.dtoRegistry.Match(res.Result)

		return dto.Result{
			Result: d,
			Key:    res.Key,
			RuleID: res.RuleID,
		}
	})
}

func (s *Service) setupRoutes() {
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{xProcessMeta},
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
