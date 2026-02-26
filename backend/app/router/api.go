package router

import (
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/auth"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	App             fiber.Router
	Cfg             *config.Config
	AuthRouter      *auth.AuthRouter
	WorkspaceRouter *workspace.WorkspaceRouter
}

func NewRouter(
	fiber *fiber.App,
	cfg *config.Config,
	authRouter *auth.AuthRouter,
	workspaceRouter *workspace.WorkspaceRouter,
) *Router {
	return &Router{
		App:             fiber,
		Cfg:             cfg,
		AuthRouter:      authRouter,
		WorkspaceRouter: workspaceRouter,
	}
}

func (r *Router) Register() {
	// Test Routes
	r.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong! 👋")
	})

	// routes of modules
	r.AuthRouter.RegisterAuthRoutes()
	r.WorkspaceRouter.RegisterWorkspaceRoutes()
}
