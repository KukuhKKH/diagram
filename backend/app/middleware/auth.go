package middleware

import (
	"time"

	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware holds dependencies for authentication middleware
type AuthMiddleware struct {
	cfg   *config.Config
	store *session.Store
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(cfg *config.Config, store *session.Store) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:   cfg,
		store: store,
	}
}

// RequireAuth is a middleware for session-based authentication
func (am *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if am.store == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Session store not initialized",
			})
		}

		sess, err := am.store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		userID := sess.Get("user_id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Set user_id in Locals for easy access in controllers
		c.Locals("user_id", userID)

		return c.Next()
	}
}

// GetUserID helper to get user_id from Locals
func GetUserID(c *fiber.Ctx) uint64 {
	id := c.Locals("user_id")
	if id == nil {
		return 0
	}

	switch v := id.(type) {
	case uint64:
		return v
	case int:
		return uint64(v)
	case float64:
		return uint64(v)
	}

	return 0
}

// Protected is JWT-based authentication middleware
func (am *AuthMiddleware) Protected() fiber.Handler {
	if am.cfg.Middleware.Jwt.Secret == "" {
		panic("JWT secret is not set")
	}

	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(am.cfg.Middleware.Jwt.Secret),
		ErrorHandler: jwtError,
		Claims:       &JWTClaims{},
		TokenLookup:  "header:Authorization,cookie:" + am.cfg.Cookie.Name,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}

	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

type JWTClaims struct {
	Token  string `json:"token"`
	Type   string `json:"type"`
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateTokenAccess generates a JWT token for a user
func (am *AuthMiddleware) GenerateTokenAccess(userID uint64) (*JWTClaims, error) {
	mySigningKey := []byte(am.cfg.Middleware.Jwt.Secret)

	expiration := am.cfg.Middleware.Jwt.Expiration
	if expiration < 1000000 {
		expiration = expiration * time.Second
	}

	// Create the Claims
	now := time.Now()
	exp := now.Add(expiration)
	claims := &JWTClaims{
		Type:   "Bearer",
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    am.cfg.Middleware.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}

	claims.Token = ss
	return claims, nil
}
