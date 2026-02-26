package main

import (
	"go.uber.org/fx"

	"git.dev.siap.id/kukuhkkh/app-diagram/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/auth"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/router"
	"git.dev.siap.id/kukuhkkh/app-diagram/internal/bootstrap"
	"git.dev.siap.id/kukuhkkh/app-diagram/internal/bootstrap/database"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/session"
	fxzerolog "github.com/efectn/fx-zerolog"
	_ "go.uber.org/automaxprocs"
)

// @title                       Aplikasi Diagram API
// @version                     1.0
// @description                 This is a sample API documentation.
// @termsOfService              http://swagger.io/terms/
// @contact.name                Kukuh Rahmadani
// @contact.email               krahmadani1@gmail.com
// @license.name                MIT
// @license.url                 https://opensource.org/license/mit
// @host                        localhost:8080
// @schemes                     http https
// @BasePath                    /
func main() {
	fx.New(
		// config
		fx.Provide(config.NewConfig),
		// logging
		fx.Provide(bootstrap.NewLogger),
		// fiber
		fx.Provide(bootstrap.NewFiber),
		// database
		fx.Provide(database.NewDatabase),
		// session
		fx.Provide(session.NewStore),
		// middleware
		fx.Provide(middleware.NewMiddleware),
		fx.Provide(middleware.NewAuthMiddleware),
		// router
		fx.Provide(router.NewRouter),

		// provide modules
		auth.NewAuthModule,

		// start aplication
		fx.Invoke(bootstrap.Start),

		// define logger
		fx.WithLogger(fxzerolog.Init()),
	).Run()
}
