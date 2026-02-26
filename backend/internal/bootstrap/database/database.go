package database

import (
	"git.dev.siap.id/kukuhkkh/app-diagram/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
	ulog "git.dev.siap.id/kukuhkkh/app-diagram/utils/logger"
	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database setup database with gorm
type Database struct {
	DB  *gorm.DB
	Log zerolog.Logger
	Cfg *config.Config
}

type Seeder interface {
	Seed(*gorm.DB) error
	Count() (int, error)
}

func NewDatabase(cfg *config.Config, log zerolog.Logger) *Database {
	db := &Database{
		Cfg: cfg,
		Log: log,
	}

	return db
}

// ConnectDatabase connect database
func (_db *Database) ConnectDatabase() {
	var dialector gorm.Dialector

	if _db.Cfg.DB.Postgres.DSN != "" {
		dialector = postgres.Open(_db.Cfg.DB.Postgres.DSN)
		_db.Log.Info().Msg("Connecting to PostgreSQL...")
	} else if _db.Cfg.DB.Mysql.DSN != "" {
		dialector = mysql.Open(_db.Cfg.DB.Mysql.DSN)
		_db.Log.Info().Msg("Connecting to MySQL...")
	} else {
		_db.Log.Panic().Msg("No database configuration found (MySQL or PostgreSQL)")
	}

	conn, err := gorm.Open(dialector, &gorm.Config{
		Logger: &ulog.GormLogger{Log: _db.Log},
	})

	if err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to connect the database!")
	} else {
		_db.Log.Info().Msg("Connected the database succesfully!")
	}

	_db.DB = conn
}

// ShutdownDatabase shutdown database
func (_db *Database) ShutdownDatabase() {
	sqlDB, err := _db.DB.DB()
	if err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to shutdown the database!")
	} else {
		_db.Log.Info().Msg("Shutdown the database succesfully!")
	}

	err = sqlDB.Close()
	if err != nil {
		return
	}
}

// MigrateModels migrate models
func (_db *Database) MigrateModels() {
	if err := _db.DB.AutoMigrate(
		Models()...,
	); err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to migrate the database!")
	}
}

// Models list of models for migration
func Models() []interface{} {
	return []interface{}{
		schema.User{},
		schema.Workspace{},
		schema.Document{},
		schema.DocumentVersion{},
		schema.SharedAccess{},
	}
}

// SeedModels seed data
func (_db *Database) SeedModels(seeder ...Seeder) {
	for _, seed := range seeder {
		count, err := seed.Count()
		if err != nil {
			_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the database!")
		}

		if count == 0 {
			if err := seed.Seed(_db.DB); err != nil {
				_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the database!")
			}

			_db.Log.Info().Msg("Seeded the database succesfully!")
		} else {
			_db.Log.Info().Msg("Database is already seeded!")
		}
	}

	_db.Log.Info().Msg("Seeded the database succesfully!")
}
