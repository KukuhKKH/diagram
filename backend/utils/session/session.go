package session

import (
	"fmt"
	"time"

	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	"github.com/gofiber/storage/redis/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NewStore(cfg *config.Config) *session.Store {
	var storage fiber.Storage

	// Determine session storage driver
	driver := cfg.Middleware.Session.Driver
	if driver == "" {
		driver = "memory" // default
	}

	switch driver {
	case "redis":
		// Redis storage configuration
		if cfg.Middleware.Session.RedisHost == "" {
			log.Warn().Msg("Redis host not configured, falling back to memory storage")
			storage = memory.New()
		} else {
			redisAddr := fmt.Sprintf("%s:%d", cfg.Middleware.Session.RedisHost, cfg.Middleware.Session.RedisPort)
			storage = redis.New(redis.Config{
				Host:     redisAddr,
				Password: cfg.Middleware.Session.RedisPasswd,
				Database: cfg.Middleware.Session.RedisDB,
				Reset:    false,
			})
			log.Info().Msgf("Session storage: Redis at %s", redisAddr)
		}
	case "memory":
		fallthrough
	default:
		storage = memory.New()
		log.Info().Msg("Session storage: In-memory (not recommended for production)")
	}

	// Session expiration
	expiration := cfg.Middleware.Session.Expiration
	if expiration == 0 {
		expiration = 24 * time.Hour // default 24 hours
	} else if expiration < 1000000 {
		// If expiration is small (likely in seconds), convert to duration
		expiration = expiration * time.Second
	}

	store := session.New(session.Config{
		KeyLookup:      fmt.Sprintf("cookie:%s", cfg.Middleware.Session.Name),
		CookieHTTPOnly: cfg.Cookie.HTTPOnly,
		CookieSecure:   cfg.Cookie.Secure,
		CookieSameSite: cfg.Cookie.SameSite,
		CookiePath:     "/",
		Expiration:     expiration,
		Storage:        storage,
		KeyGenerator: func() string {
			return uuid.New().String()
		},
	})

	return store
}

func Get(c *fiber.Ctx, store *session.Store) (*session.Session, error) {
	return store.Get(c)
}
