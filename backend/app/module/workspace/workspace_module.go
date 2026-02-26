package workspace

import (
	"git.dev.siap.id/kukuhkkh/app-diagram/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace/controller"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace/repository"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// WorkspaceRouter adalah router untuk workspace module
type WorkspaceRouter struct {
	App        fiber.Router
	Controller *controller.Controller
	AuthMW     *middleware.AuthMiddleware
}

// Module adalah FX module untuk workspace
var NewWorkspaceModule = fx.Options(
	// register repository
	fx.Provide(repository.NewWorkspaceRepository),

	// register service
	fx.Provide(service.NewWorkspaceService),

	// register controller
	controller.Module,

	// register router
	fx.Provide(NewWorkspaceRouter),
)

// NewWorkspaceRouter membuat instance baru dari WorkspaceRouter
func NewWorkspaceRouter(
	app *fiber.App,
	ctrl *controller.Controller,
	authMW *middleware.AuthMiddleware,
) *WorkspaceRouter {
	return &WorkspaceRouter{
		App:        app,
		Controller: ctrl,
		AuthMW:     authMW,
	}
}

// RegisterWorkspaceRoutes mendaftarkan routes untuk workspace
func (_i *WorkspaceRouter) RegisterWorkspaceRoutes() {
	// define controllers
	workspaceController := _i.Controller.Workspace

	_i.App.Route("/api/v1", func(router fiber.Router) {
		workspaceRoutes := router.Group("/workspaces", _i.AuthMW.RequireAuth())

		workspaceRoutes.Post("", workspaceController.CreateWorkspace)
		workspaceRoutes.Get("", workspaceController.ListWorkspaces)
		workspaceRoutes.Get("/:id", workspaceController.GetWorkspace)
		workspaceRoutes.Put("/:id", workspaceController.UpdateWorkspace)
		workspaceRoutes.Delete("/:id", workspaceController.DeleteWorkspace)
	})
}
