package controller

import (
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/auth/service"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Controller struct {
	Auth AuthController
}

func NewController(authService service.AuthService, cfg *config.Config, sessStore *session.Store) *Controller {
	return &Controller{
		Auth: NewAuthController(authService, cfg, sessStore),
	}
}
