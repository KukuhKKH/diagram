package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"git.dev.siap.id/kukuhkkh/app-diagram/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/auth/request"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/auth/response"
	user_repo "git.dev.siap.id/kukuhkkh/app-diagram/app/module/user/repository"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/helpers"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// AuthService
type AuthService interface {
	Login(req request.LoginRequest) (res response.LoginResponse, err error)
	Register(req request.RegisterRequest) (res response.RegisterResponse, err error)
	Me(userID uint64) (res response.UserResponse, err error)

	// OIDC methods
	GetAuthURL() (url, state, verifier string, err error)
	HandleCallback(code, state, sessionState, sessionVerifier string) (userID uint64, err error)
	GetLogoutURL() string
}

type userService struct {
	userRepo user_repo.UserRepository
	cfg      *config.Config
	authMw   *middleware.AuthMiddleware
}

// init AuthService
func NewAuthService(userRepo user_repo.UserRepository, cfg *config.Config, authMw *middleware.AuthMiddleware) AuthService {
	return &userService{
		userRepo: userRepo,
		cfg:      cfg,
		authMw:   authMw,
	}
}

func (_i *userService) Login(req request.LoginRequest) (res response.LoginResponse, err error) {
	// check user by email
	user, err := _i.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		err = errors.New("Email or password is incorrect")
		return
	}

	if user == nil {
		err = errors.New("Email or password is incorrect")
		return
	}

	// check password
	if !user.ComparePassword(req.Password) {
		err = errors.New("Email or password is incorrect")
		return
	}

	// do create token
	claims, err := _i.authMw.GenerateTokenAccess(user.ID)
	if err != nil {
		return
	}

	res.Token = claims.Token
	res.Type = claims.Type
	res.ExpiresAt = claims.ExpiresAt.Unix()

	return
}

func (_i *userService) Register(req request.RegisterRequest) (res response.RegisterResponse, err error) {
	// check user by email
	user, err := _i.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return
	}

	if user != nil {
		err = errors.New("email already exists")
		return
	}

	// do create user
	user = &schema.User{
		Email:    req.Email,
		Password: &req.Password,
	}

	user, err = _i.userRepo.CreateUser(user)
	if err != nil {
		return
	}

	res.ID = user.ID
	res.Email = user.Email

	return
}

func (_i *userService) Me(userID uint64) (res response.UserResponse, err error) {
	// check user by id
	user, err := _i.userRepo.FindUserByID(userID)
	if err != nil {
		return
	}

	if user == nil {
		err = errors.New("user not found")
		return
	}

	res.ID = user.ID
	res.Name = user.Name
	res.Email = user.Email

	return
}

func (s *userService) GetAuthURL() (url, state, verifier string, err error) {
	state = helpers.GenerateRandomString(32)
	verifier = helpers.GenerateRandomString(64)
	challenge := helpers.GenerateCodeChallenge(verifier)

	conf := s.getOauth2Config()
	url = conf.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	return
}

func (s *userService) HandleCallback(code, state, sessionState, sessionVerifier string) (userID uint64, err error) {
	if state != sessionState {
		return 0, errors.New("invalid state")
	}

	ctx := context.Background()
	endpoint := strings.TrimSuffix(s.cfg.Sso.Logto.Endpoint, "/")
	if !strings.HasSuffix(endpoint, "/oidc") {
		endpoint = endpoint + "/oidc"
	}

	provider, err := oidc.NewProvider(ctx, endpoint)
	if err != nil {
		return 0, err
	}

	conf := s.getOauth2Config()
	token, err := conf.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", sessionVerifier))
	if err != nil {
		return 0, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return 0, errors.New("no id_token in token response")
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: s.cfg.Sso.Logto.ClientId})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return 0, err
	}

	var claims struct {
		Subject string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return 0, err
	}

	// Provision user
	user, err := s.userRepo.FindUserByLogtoSub(claims.Subject)
	now := time.Now()
	if err != nil {
		// Create new user
		user = &schema.User{
			LogtoSub:  claims.Subject,
			Email:     claims.Email,
			Name:      claims.Name,
			LastLogin: &now,
		}
		user, err = s.userRepo.CreateUser(user)
		if err != nil {
			return 0, err
		}
	} else {
		// Update last login
		user.LastLogin = &now
		if err := s.userRepo.UpdateUser(user); err != nil {
			return 0, err
		}
	}

	return user.ID, nil
}

func (s *userService) GetLogoutURL() string {
	endpoint := strings.TrimSuffix(s.cfg.Sso.Logto.Endpoint, "/")
	basePath := endpoint
	if !strings.HasSuffix(endpoint, "/oidc") {
		basePath = endpoint + "/oidc"
	}

	// Logto end session endpoint
	// format: <base>/session/end?post_logout_redirect_uri=<uri>&client_id=<id>
	logoutURL := fmt.Sprintf("%s/session/end?client_id=%s&post_logout_redirect_uri=%s",
		basePath,
		s.cfg.Sso.Logto.ClientId,
		s.cfg.Sso.Logto.PostLogoutRedirectUri,
	)

	return logoutURL
}

func (s *userService) getOauth2Config() *oauth2.Config {
	endpoint := strings.TrimSuffix(s.cfg.Sso.Logto.Endpoint, "/")

	basePath := endpoint
	if !strings.HasSuffix(endpoint, "/oidc") {
		basePath = endpoint + "/oidc"
	}

	return &oauth2.Config{
		ClientID:     s.cfg.Sso.Logto.ClientId,
		ClientSecret: s.cfg.Sso.Logto.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  basePath + "/auth",
			TokenURL: basePath + "/token",
		},
		RedirectURL: s.cfg.Sso.Logto.CallbackUrl,
		Scopes:      []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile", "email"},
	}
}
