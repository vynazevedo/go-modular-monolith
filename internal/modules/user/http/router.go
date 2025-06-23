package http

import "github.com/gofiber/fiber/v2"

func (h *UserHandlers) RegisterRoutes(router fiber.Router) {
	router.Post("/", h.CreateUser)
	router.Get("/:id", h.GetUser)
	router.Put("/:id", h.UpdateUser)
	router.Delete("/:id", h.DeleteUser)
	router.Put("/:id/activate", h.ActivateUser)
	router.Put("/:id/deactivate", h.DeactivateUser)
	router.Get("/", h.ListUsers)
}
