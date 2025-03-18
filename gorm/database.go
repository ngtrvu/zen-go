package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate"
	migrate_postgres "github.com/golang-migrate/migrate/database/postgres"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/ngtrvu/zen-go/gorm/migrator"
	"github.com/ngtrvu/zen-go/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

type Database struct {
	GormDB   *gorm.DB
	Migrator *migrator.Migrator
}

type DBConfig struct {
	DB_HOST                      string `config:"DB_HOST"`
	DB_USERNAME                  string `config:"DB_USERNAME"`
	DB_PASSWORD                  string `config:"DB_PASSWORD"`
	DB_NAME                      string `config:"DB_NAME"`
	DB_PORT                      string `config:"DB_PORT"`
	SSLMode                      string `config:"DB_SSL_MODE"`
	TimeZone                     string `config:"DB_TIME_ZONE"`
	DB_POOL_MAX_OPEN_CONNECTIONS int    `config:"DB_POOL_MAX_OPEN_CONNECTIONS"`
	DB_POOL_MAX_IDLE_CONNECTIONS int    `config:"DB_POOL_MAX_IDLE_CONNECTIONS"`
	DB_METRICS_ENABLED           bool   `config:"DB_METRICS_ENABLED"`
}

func NewDatabase(cfg *DBConfig) (db *Database, err error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DB_HOST,
		cfg.DB_USERNAME,
		cfg.DB_PASSWORD,
		cfg.DB_NAME,
		cfg.DB_PORT,
		cfg.SSLMode,
		cfg.TimeZone,
	)
	pgCfg := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}
	pg := postgres.New(pgCfg)
	gormDB, err := gorm.Open(pg, &gorm.Config{
		Logger: log.NewGormLogger("info"),
	})
	if err != nil {
		log.Error("failed to connect database: %v", err)
		return
	}
	if cfg.DB_METRICS_ENABLED {
		gormDB.Use(prometheus.New(prometheus.Config{
			DBName:          cfg.DB_NAME, // `DBName` as metrics label
			RefreshInterval: 15,          // refresh metrics interval (default 15 seconds)
			MetricsCollector: []prometheus.MetricsCollector{
				&prometheus.Postgres{VariableNames: []string{"Threads_running"}},
			},
			// Labels: map[string]string{
			// 	"instance": "127.0.0.1", // config custom labels if necessary
			// },
		}))
	}

	// Set connection pool settings
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Error("failed to get sql.DB from gorm.DB: %v", err)
		return
	}

	// Set the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.DB_POOL_MAX_OPEN_CONNECTIONS)

	// Set the maximum number of idle connections in the pool.
	sqlDB.SetMaxIdleConns(cfg.DB_POOL_MAX_IDLE_CONNECTIONS)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	db = &Database{GormDB: gormDB, Migrator: migrator.NewMigrator(gormDB)}

	return
}

func (db *Database) CheckMigration(ctx context.Context, migrationData *bindata.AssetSource) (*migrate.Migrate, error) {
	sqlDB, _ := db.GormDB.DB()
	sourceInstance, err := bindata.WithInstance(migrationData)
	if err != nil {
		log.Error("read source instance from bindata error: %v", err)
		return nil, err
	}
	targetInstance, err := migrate_postgres.WithInstance(sqlDB, &migrate_postgres.Config{})
	if err != nil {
		log.Error("init target instance error: %v", err)
		return nil, err
	}
	m, err := migrate.NewWithInstance("go-bindata", sourceInstance, "postgres", targetInstance)
	if err != nil {
		log.Error("init migration instance error: %v", err)
		return nil, err
	}
	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Info("checking current version completed, there is no mogration")
		} else {
			log.Error("checking current verion error: %v", err)
			return nil, err
		}
	} else {
		log.Info("schema version before running migration: %d", version)
		if dirty {
			log.Warn("current version is dirty, running migrations failed")
			return nil, err
		}
	}
	return m, nil
}

func (db *Database) RunMigration(ctx context.Context, migrationData *bindata.AssetSource) error {
	m, err := db.CheckMigration(ctx, migrationData)
	if err != nil {
		return err
	}

	// running migration
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("run migrations completed, there is no change")
			return nil
		} else {
			log.Error("migration error: %v", err)
			return err
		}
	}

	log.Info("run migrations successfully")
	return nil
}

func (db *Database) DownMigration(ctx context.Context, migrationData *bindata.AssetSource) error {
	sqlDB, _ := db.GormDB.DB()
	sourceInstance, err := bindata.WithInstance(migrationData)
	if err != nil {
		log.Error("Read source instance from bindata error %v", err)
		return err
	}
	targetInstance, err := migrate_postgres.WithInstance(sqlDB, &migrate_postgres.Config{})
	if err != nil {
		log.Error("Init target instance error %v", err)
		return err
	}
	m, err := migrate.NewWithInstance("go-bindata", sourceInstance, "postgres", targetInstance)
	if err != nil {
		log.Error("Init migration instance error %v", err)
		return err
	}
	_, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
		} else {
			return err
		}
	} else {
		if dirty {
			return err
		}
	}

	// running migration
	err = m.Down()
	if err != nil {
		fmt.Println("err down migration", err)
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
