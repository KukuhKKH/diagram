package middleware

import (
	"net/http"
	"time"

	"git.dev.siap.id/kukuhkkh/app-diagram/utils"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Middleware is a struct that contains all the middleware functions
type Middleware struct {
	App *fiber.App
	Cfg *config.Config
}

func NewMiddleware(app *fiber.App, cfg *config.Config, sessStore *session.Store) *Middleware {
	return &Middleware{
		App: app,
		Cfg: cfg,
	}
}

// Register registers all the middleware functions
func (m *Middleware) Register() {
	// CORS - Use FrontendUrl if set, otherwise use config AllowOrigins
	corsOrigins := m.Cfg.Middleware.Cors.AllowOrigins
	if m.Cfg.App.FrontendUrl != "" {
		corsOrigins = m.Cfg.App.FrontendUrl
	}

	m.App.Use(cors.New(cors.Config{
		Next:             utils.IsEnabled(m.Cfg.Middleware.Cors.Enable),
		AllowOrigins:     corsOrigins,
		AllowHeaders:     m.Cfg.Middleware.Cors.AllowHeaders,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Add Extra Middlewares
	m.App.Use(limiter.New(limiter.Config{
		Next:       utils.IsEnabled(m.Cfg.Middleware.Limiter.Enable),
		Max:        m.Cfg.Middleware.Limiter.Max,
		Expiration: m.Cfg.Middleware.Limiter.Expiration * time.Second,
	}))

	m.App.Use(compress.New(compress.Config{
		Next:  utils.IsEnabled(m.Cfg.Middleware.Compress.Enable),
		Level: m.Cfg.Middleware.Compress.Level,
	}))

	m.App.Use(recover.New(recover.Config{
		Next: utils.IsEnabled(m.Cfg.Middleware.Recover.Enable),
	}))

	m.App.Use(pprof.New(pprof.Config{
		Next: utils.IsEnabled(m.Cfg.Middleware.Pprof.Enable),
	}))

	m.App.Use(filesystem.New(filesystem.Config{
		Next:   utils.IsEnabled(m.Cfg.Middleware.FileSystem.Enable),
		Root:   http.Dir(m.Cfg.Middleware.FileSystem.Root),
		Browse: m.Cfg.Middleware.FileSystem.Browse,
		MaxAge: m.Cfg.Middleware.FileSystem.MaxAge,
	}))

	m.App.Get(m.Cfg.Middleware.Monitor.Path, monitor.New(monitor.Config{
		Next: utils.IsEnabled(m.Cfg.Middleware.Monitor.Enable),
	}))

	// Static Storage for Local Driver
	if m.Cfg.Storage.Driver == "local" {
		m.App.Static("/storage", m.Cfg.Storage.Local.Path)
	}
}
