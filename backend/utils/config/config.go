package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// app struct config
type app = struct {
	Name        string        `toml:"name"`
	Port        string        `toml:"port"`
	FrontendUrl string        `toml:"frontend_url"`
	PrintRoutes bool          `toml:"print-routes"`
	Prefork     bool          `toml:"prefork"`
	Production  bool          `toml:"production"`
	BodyLimit   int           `toml:"body-limit"`
	IdleTimeout time.Duration `toml:"idle-timeout"`
	TLS         struct {
		Enable   bool   `toml:"enable"`
		CertFile string `toml:"cert-file"`
		KeyFile  string `toml:"key-file"`
	}
}

// db struct config
type db = struct {
	Mysql struct {
		DSN string `toml:"dsn"`
	}
	Postgres struct {
		DSN string `toml:"dsn"`
	}
}

// log struct config
type logger = struct {
	TimeFormat string        `toml:"time-format"`
	Level      zerolog.Level `toml:"level"`
	Prettier   bool          `toml:"prettier"`
}

// middleware
type middleware = struct {
	Compress struct {
		Enable bool
		Level  compress.Level
	}

	Recover struct {
		Enable bool
	}

	Monitor struct {
		Enable bool
		Path   string
	}

	Pprof struct {
		Enable bool
	}

	Limiter struct {
		Enable     bool
		Max        int
		Expiration time.Duration `toml:"expiration_seconds"`
	}

	FileSystem struct {
		Enable bool
		Browse bool
		MaxAge int `toml:"max_age"`
		Index  string
		Root   string
	}

	Jwt struct {
		Secret     string        `toml:"secret"`
		Issuer     string        `toml:"issuer"`
		Expiration time.Duration `toml:"expiration_seconds"`
	}

	Cors struct {
		Enable       bool   `toml:"enable"`
		AllowOrigins string `toml:"allow_origins"`
		AllowHeaders string `toml:"allow_headers"`
	}

	Session struct {
		Enable      bool          `toml:"enable"`
		Name        string        `toml:"name"`
		Driver      string        `toml:"driver"`     // "memory" or "redis"
		Expiration  time.Duration `toml:"expiration"` // in seconds
		RedisHost   string        `toml:"redis_host"`
		RedisPort   int           `toml:"redis_port"`
		RedisPasswd string        `toml:"redis_passwd"`
		RedisDB     int           `toml:"redis_db"`
	}
}

type cookie struct {
	Name     string `toml:"name"`
	HTTPOnly bool   `toml:"http_only"`
	Secure   bool   `toml:"secure"`
	SameSite string `toml:"same_site"`
}

type storage = struct {
	Driver string `toml:"driver"`

	Local struct {
		Path string `toml:"path"`
	} `toml:"local"`

	Ftp struct {
		Host      string `toml:"host"`
		Port      int    `toml:"port"`
		User      string `toml:"user"`
		Password  string `toml:"password"`
		BaseDir   string `toml:"base_dir"`
		PublicUrl string `toml:"public_url"`
	} `toml:"ftp"`

	S3 struct {
		Endpoint  string `toml:"endpoint"`
		AccessKey string `toml:"access_key"`
		SecretKey string `toml:"secret_key"`
		Bucket    string `toml:"bucket"`
		Region    string `toml:"region"`
		UseSsl    bool   `toml:"use_ssl"`
	} `toml:"s3"`
}

type Sso struct {
	Logto struct {
		Endpoint              string `toml:"endpoint"`
		ClientId              string `toml:"client_id"`
		ClientSecret          string `toml:"client_secret"`
		PostLogoutRedirectUri string `toml:"post_logout_redirect_uri"`
		CallbackUrl           string `toml:"callback_url"`
	} `toml:"logto"`
}

type Config struct {
	App        app
	DB         db
	Logger     logger
	Middleware middleware
	Cookie     cookie
	Storage    storage
	Sso        Sso
}

// Validate validates critical configuration fields
func (c *Config) Validate() error {
	var errs []string

	// Validate App
	if c.App.Name == "" {
		errs = append(errs, "app.name is required")
	}
	if c.App.Port == "" {
		errs = append(errs, "app.port is required")
	} else {
		// Validate port format
		portStr := c.App.Port
		if strings.Contains(portStr, ":") {
			portStr = strings.Split(portStr, ":")[1]
		}
		if _, err := strconv.Atoi(portStr); err != nil {
			errs = append(errs, fmt.Sprintf("app.port '%s' is not a valid port number", c.App.Port))
		}
	}

	// Validate Database - at least one must be configured
	if c.DB.Postgres.DSN == "" && c.DB.Mysql.DSN == "" {
		errs = append(errs, "database configuration is required (either postgres or mysql)")
	}

	// Validate JWT if session is not enabled
	if !c.Middleware.Session.Enable {
		if c.Middleware.Jwt.Secret == "" {
			errs = append(errs, "middleware.jwt.secret is required when session is disabled")
		}
		if len(c.Middleware.Jwt.Secret) < 32 {
			errs = append(errs, "middleware.jwt.secret should be at least 32 characters for security")
		}
	}

	// Validate TLS if enabled
	if c.App.TLS.Enable {
		if c.App.TLS.CertFile == "" {
			errs = append(errs, "app.tls.cert_file is required when TLS is enabled")
		}
		if c.App.TLS.KeyFile == "" {
			errs = append(errs, "app.tls.key_file is required when TLS is enabled")
		}
	}

	// Validate Session
	if c.Middleware.Session.Enable && c.Middleware.Session.Name == "" {
		errs = append(errs, "middleware.session.name is required when session is enabled")
	}

	// Validate Storage
	if c.Storage.Driver != "" && c.Storage.Driver != "local" && c.Storage.Driver != "ftp" && c.Storage.Driver != "s3" {
		errs = append(errs, fmt.Sprintf("storage.driver '%s' is not valid (must be: local, ftp, or s3)", c.Storage.Driver))
	}

	if c.Storage.Driver == "local" && c.Storage.Local.Path == "" {
		errs = append(errs, "storage.local.path is required when storage driver is 'local'")
	}

	if c.Storage.Driver == "s3" {
		if c.Storage.S3.Bucket == "" {
			errs = append(errs, "storage.s3.bucket is required when storage driver is 's3'")
		}
		if c.Storage.S3.AccessKey == "" || c.Storage.S3.SecretKey == "" {
			errs = append(errs, "storage.s3.access_key and secret_key are required when storage driver is 's3'")
		}
	}

	if c.Storage.Driver == "ftp" {
		if c.Storage.Ftp.Host == "" {
			errs = append(errs, "storage.ftp.host is required when storage driver is 'ftp'")
		}
	}

	// Validate CORS
	if c.Middleware.Cors.Enable && c.App.Production {
		if c.Middleware.Cors.AllowOrigins == "" || c.Middleware.Cors.AllowOrigins == "*" {
			log.Warn().Msg("⚠️  CORS is enabled in production with wildcard (*) origin - this is a security risk!")
		}
	}

	if len(errs) > 0 {
		return errors.New("configuration validation failed:\n  - " + strings.Join(errs, "\n  - "))
	}

	return nil
}

// ParseConfig func to parse config
func ParseConfig(name string) (contents *Config, err error) {
	filename := name + ".toml"
	candidates := []string{filename, filepath.Join("config", filename)}

	exePath, exeErr := os.Executable()
	if exeErr == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates, filepath.Join(exeDir, filename), filepath.Join(exeDir, "config", filename))
	}

	if wd, wdErr := os.Getwd(); wdErr == nil {
		candidates = append(candidates, filepath.Join(wd, filename), filepath.Join(wd, "config", filename))
	}

	var readErr error
	for _, p := range candidates {
		p = filepath.Clean(p)
		file, e := os.ReadFile(p)
		if e != nil {
			readErr = e
			continue
		}

		err = toml.Unmarshal(file, &contents)
		if err != nil {
			fmt.Printf("❌ CRITICAL: Gagal parse config di lokasi: %s. Error: %v\n", p, err)
			return nil, err
		}

		return contents, nil
	}

	fmt.Printf("❌ CRITICAL: Gagal baca config. Mencoba lokasi: %v. Last error: %v\n", candidates, readErr)
	if readErr == nil {
		readErr = fmt.Errorf("config file not found in candidates: %v", candidates)
	}

	return nil, readErr
}

// NewConfig initialize config
func NewConfig() *Config {
	config, err := ParseConfig("config")
	if err != nil {
		if !fiber.IsChild() {
			log.Panic().Err(err).Msg("config not found")
		} else {
			log.Error().Err(err).Msg("child process failed to load config")
		}
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		if !fiber.IsChild() {
			log.Panic().Err(err).Msg("configuration validation failed")
		} else {
			log.Error().Err(err).Msg("child process configuration validation failed")
		}
	}

	return config
}

// ParseAddress func to parse address
func ParseAddress(raw string) (host, port string) {
	if i := strings.LastIndex(raw, ":"); i > 0 {
		return raw[:i], raw[i+1:]
	}

	return raw, ""
}
