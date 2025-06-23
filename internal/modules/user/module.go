package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vynazevedo/go-modular-monolith/internal/modules/user/app"
	"github.com/vynazevedo/go-modular-monolith/internal/modules/user/domain"
	"github.com/vynazevedo/go-modular-monolith/internal/modules/user/http"
	"github.com/vynazevedo/go-modular-monolith/internal/modules/user/infra"
	"github.com/vynazevedo/go-modular-monolith/internal/shared/module"
	"gorm.io/gorm"
)

type Module struct {
	service  *app.UserService
	handlers *http.UserHandlers
}

func NewModule(db *gorm.DB) *Module {

	repo := infra.NewGormUserRepository(db)
	service := app.NewUserService(repo)
	handlers := http.NewUserHandlers(service)

	return &Module{
		service:  service,
		handlers: handlers,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	m.handlers.RegisterRoutes(router.Group("/users"))
}

func (m *Module) QueryService() domain.UserQueryService {
	return m.service
}

func (m *Module) GetModels() []interface{} {
	return []interface{}{
		&infra.UserModel{},
	}
}

var _ module.Module = (*Module)(nil)
