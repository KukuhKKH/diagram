package auth

import (
	"git.dev.siap.id/kukuhkkh/app-diagram/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/auth/controller"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/auth/service"
	user_repo "git.dev.siap.id/kukuhkkh/app-diagram/app/module/user/repository"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// struct of AuthRouter
type AuthRouter struct {
	App            fiber.Router
	Controller     *controller.Controller
	AuthMiddleware *middleware.AuthMiddleware
}

// register bulky of auth module
var NewAuthModule = fx.Options(
	// register repository of auth module
	fx.Provide(user_repo.NewUserRepository),

	// register service of auth module
	fx.Provide(service.NewAuthService),

	// register controller of auth module
	fx.Provide(controller.NewController),

	// register router of auth module
	fx.Provide(NewAuthRouter),
)

// init AuthRouter
func NewAuthRouter(fiber *fiber.App, controller *controller.Controller, authMw *middleware.AuthMiddleware) *AuthRouter {
	return &AuthRouter{
		App:            fiber,
		Controller:     controller,
		AuthMiddleware: authMw,
	}
}

// register routes of auth module
func (_i *AuthRouter) RegisterAuthRoutes() {
	// define controllers
	authController := _i.Controller.Auth

	// define routes
	_i.App.Route("/auth", func(router fiber.Router) {
		router.Get("/login", authController.Login)
		router.Get("/callback", authController.Callback)
		router.Get("/me", _i.AuthMiddleware.RequireAuth(), authController.Me)
		router.Post("/logout", _i.AuthMiddleware.RequireAuth(), authController.Logout)
	})
}
