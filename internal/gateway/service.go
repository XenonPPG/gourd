package gateway

import (
	"gourd/internal/counter"
	"gourd/internal/document"

	"github.com/gofiber/fiber/v3"
)

type Service struct {
	documentService *document.Service
	counterService  *counter.Service
}

func New(documentService *document.Service, counterService *counter.Service) *Service {
	return &Service{
		documentService: documentService,
		counterService:  counterService,
	}
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
