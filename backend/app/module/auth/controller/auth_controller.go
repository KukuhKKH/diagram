package controller

import (
	"time"

	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/auth/service"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/response"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

type authController struct {
	authService    service.AuthService
	cfg            *config.Config
	cookieTemplate *fiber.Cookie
	sessStore      *fsession.Store
}

type AuthController interface {
	Login(c *fiber.Ctx) error
	Callback(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

func NewAuthController(authService service.AuthService, cfg *config.Config, sessStore *fsession.Store) AuthController {
	tmpl := &fiber.Cookie{
		Name:     cfg.Cookie.Name,
		HTTPOnly: cfg.Cookie.HTTPOnly,
		Secure:   cfg.Cookie.Secure,
		SameSite: cfg.Cookie.SameSite,
		Path:     "/",
	}

	return &authController{
		authService:    authService,
		cfg:            cfg,
		cookieTemplate: tmpl,
		sessStore:      sessStore,
	}
}

func (_i *authController) Login(c *fiber.Ctx) error {
	url, state, verifier, err := _i.authService.GetAuthURL()
	if err != nil {
		return err
	}

	sess, err := _i.sessStore.Get(c)
	if err != nil {
		return err
	}

	if err := sess.Regenerate(); err != nil {
		return err
	}

	sess.Set("oidc_state", state)
	sess.Set("oidc_verifier", verifier)

	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect(url)
}

func (_i *authController) Callback(c *fiber.Ctx) error {
	sess, err := _i.sessStore.Get(c)
	if err != nil {
		return err
	}

	sessionState, _ := sess.Get("oidc_state").(string)
	sessionVerifier, _ := sess.Get("oidc_verifier").(string)

	userID, err := _i.authService.HandleCallback(c.Query("code"), c.Query("state"), sessionState, sessionVerifier)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	sess.Set("user_id", userID)
	sess.Delete("oidc_state")
	sess.Delete("oidc_verifier")

	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect(_i.cfg.App.FrontendUrl)
}

func (_i *authController) Me(c *fiber.Ctx) error {
	sess, err := _i.sessStore.Get(c)
	if err != nil {
		return err
	}

	id := sess.Get("user_id")
	if id == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Handle type assertion safely
	var userID uint64
	switch v := id.(type) {
	case uint64:
		userID = v
	case float64:
		userID = uint64(v)
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid session data"})
	}

	res, err := _i.authService.Me(userID)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data:     res,
		Messages: response.Messages{"Get me success"},
		Code:     fiber.StatusOK,
	})
}

func (_i *authController) Logout(c *fiber.Ctx) error {
	sess, err := _i.sessStore.Get(c)
	if err != nil {
		return err
	}

	if err := sess.Destroy(); err != nil {
		return err
	}

	logoutURL := _i.authService.GetLogoutURL()

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Logout success"},
		Data:     fiber.Map{"logout_url": logoutURL},
		Code:     fiber.StatusOK,
	})
}

func (_i *authController) makeCookie(value string, expires time.Time) *fiber.Cookie {
	c := *(_i.cookieTemplate)
	c.Value = value
	c.Expires = expires
	return &c
}
