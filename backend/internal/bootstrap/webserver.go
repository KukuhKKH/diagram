package bootstrap

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"git.dev.siap.id/kukuhkkh/app-diagram/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/router"
	"git.dev.siap.id/kukuhkkh/app-diagram/internal/bootstrap/database"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	futils "github.com/gofiber/fiber/v2/utils"
)

// NewFiber initialize the webserver
func NewFiber(cfg *config.Config) *fiber.App {
	// setup
	bodyLimit := cfg.App.BodyLimit
	if bodyLimit <= 0 {
		bodyLimit = 100 * 1024 * 1024 // Default 100MB
	}

	app := fiber.New(fiber.Config{
		ServerHeader:          cfg.App.Name,
		AppName:               cfg.App.Name,
		Prefork:               cfg.App.Prefork,
		ErrorHandler:          response.ErrorHandler,
		IdleTimeout:           cfg.App.IdleTimeout * time.Second,
		EnablePrintRoutes:     cfg.App.PrintRoutes,
		BodyLimit:             bodyLimit,
		DisableStartupMessage: true,
		StreamRequestBody:     true,
	})

	// pass production config to check it
	response.IsProduction = cfg.App.Production

	return app
}

// Start function to start webserver
func Start(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	fiberApp *fiber.App,
	router *router.Router,
	middlewares *middleware.Middleware,
	db *database.Database,
	log zerolog.Logger,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				// ---------------------------------------------------------
				// 1. Register Middlewares & Routes
				// ---------------------------------------------------------
				middlewares.Register()
				router.Register()

				// ---------------------------------------------------------
				// 2. Prepare Host & Port Info
				// ---------------------------------------------------------
				host, port := config.ParseAddress(cfg.App.Port)
				if host == "" {
					if fiberApp.Config().Network == "tcp6" {
						host = "[::1]"
					} else {
						host = "0.0.0.0"
					}
				}

				// Gabungkan host dan port agar aman untuk Listen
				// Jika port hanya berisi angka (misal "8080"), tambahkan ":" di depan
				addr := cfg.App.Port
				if !strings.Contains(addr, ":") {
					addr = fmt.Sprintf(":%s", addr)
				}

				// ---------------------------------------------------------
				// 3. Print Startup Messages (ASCII & Info)
				// ---------------------------------------------------------
				printStartupMessage(cfg, fiberApp, log, host, port)

				// ---------------------------------------------------------
				// 4. Start Server (NON-BLOCKING / Goroutine)
				// ---------------------------------------------------------
				go func() {
					var err error

					// Cek apakah TLS di-enable
					if cfg.App.TLS.Enable {
						log.Info().Msg("🔒 TLS support was enabled.")
						err = fiberApp.ListenTLS(addr, cfg.App.TLS.CertFile, cfg.App.TLS.KeyFile)
					} else {
						err = fiberApp.Listen(addr)
					}

					if err != nil {
						log.Error().Err(err).Msg("❌ An unknown error occurred when to run server!")
					}
				}()

				// ---------------------------------------------------------
				// 5. Database Connection & Operations
				// ---------------------------------------------------------
				db.ConnectDatabase()

				if hasFlag("migrate") {
					log.Info().Msg("🛠️  Migrate flag detected. Running migration...")
					db.MigrateModels()
				}

				if hasFlag("seed") {
					log.Info().Msg("🌱 Seed flag detected. Running seeder...")
					db.SeedModels()
				}

				// Return nil agar FX tahu aplikasi berhasil start
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Info().Msg("🛑 Shutting down the app...")

				if err := fiberApp.ShutdownWithContext(ctx); err != nil {
					log.Panic().Err(err).Msg("Fiber shutdown failed")
				}

				log.Info().Msg("Running cleanup tasks...")
				log.Info().Msg("1- Shutdown the database")
				db.ShutdownDatabase()

				log.Info().Msgf("%s was successful shutdown.", cfg.App.Name)
				log.Info().Msg("\u001b[96msee you again👋\u001b[0m")

				return nil
			},
		},
	)
}

// ---------------------------------------------------------
// Helper Functions
// ---------------------------------------------------------

// hasFlag
func hasFlag(name string) bool {
	for _, arg := range os.Args {
		if arg == "-"+name || arg == "--"+name {
			return true
		}
	}

	if f := flag.Lookup(name); f != nil {
		return f.Value.String() == "true"
	}

	return false
}

func printStartupMessage(cfg *config.Config, fiberApp *fiber.App, log zerolog.Logger, host, port string) {
	// ASCII Art
	if !fiber.IsChild() {
		ascii, err := os.ReadFile("./storage/ascii_art.txt")
		if err != nil {
			log.Debug().Err(err).Msg("An unknown error occurred when to print ASCII art!")
		} else {
			for _, line := range strings.Split(futils.UnsafeString(ascii), "\n") {
				// Menggunakan fmt.Println agar format ASCII terjaga rapi
				fmt.Println(line)
			}
		}
	}

	// Information message
	log.Info().Msg(fiberApp.Config().AppName + " is running at the moment!")

	// Debug information
	if !cfg.App.Production {
		prefork := "Enabled"
		procs := runtime.GOMAXPROCS(0)
		if !cfg.App.Prefork {
			procs = 1
			prefork = "Disabled"
		}

		log.Debug().Msg("--------------------------------------------------")
		log.Debug().Msgf("Version: %s", "-")
		log.Debug().Msgf("Host: %s", host)
		log.Debug().Msgf("Port: %s", port)
		log.Debug().Msgf("Prefork: %s", prefork)
		log.Debug().Msgf("Handlers: %d", fiberApp.HandlersCount())
		log.Debug().Msgf("Procs:    %d", procs)
		log.Debug().Msgf("PID: %d", os.Getpid())
		log.Debug().Msg("--------------------------------------------------")
	}
}
