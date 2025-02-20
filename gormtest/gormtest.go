package gormtest

import (
	"context"
	"fmt"
	"strings"

	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/google/uuid"
	common_gorm "github.com/ngtrvu/zen-go/gorm"
	"github.com/ngtrvu/zen-go/gorm/migrator"
	"github.com/ngtrvu/zen-go/log"
	"github.com/ngtrvu/zen-go/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// pgConn represents the postgres database connection.
// This is used to create the test database.
// Or drop the test database after all test cases are done.
var pgConn *gorm.DB
var dbConn = make(map[string]*gorm.DB, 0)

type GormTestSuite struct {
	suite.Suite

	// Db represents the test database connection.
	Db *gorm.DB

	// dbManager represents the database manager connection. It is used to create the test database. Run up/down migration only.
	dbManager *common_gorm.Database

	// Db represents the transaction connection. Each test cases will have a separated transaction.
	Tx *gorm.DB

	// TestCaseName is the postfix of the test database name. if not set will auto generate from test suite name.
	TestCaseName string

	// TestDbName is the name of the test database.
	TestDbName string

	// Data represents test data setup.
	Data map[string]interface{}

	// Context represents the context of the test.
	Context context.Context

	// Migration Bin Data for test
	MigrationData *bindata.AssetSource

	// Migrator
	Migrator *migrator.Migrator
}

func (t *GormTestSuite) SetTestCaseName(testCaseName string) {
	t.TestCaseName = testCaseName
}

// DatabaseName returns the name of the test database.
func (t *GormTestSuite) DatabaseName() string {
	dbName := utils.GetEnvDefault("POSTGRES_DB", "stagtest")
	uuidKey := strings.Replace(uuid.New().String(), "-", "", -1)
	testDbName := fmt.Sprintf("%s_%s", dbName, uuidKey)
	return strings.ToLower(testDbName)
}

func (t *GormTestSuite) SetupMigration(migrationData *bindata.AssetSource) {
	t.MigrationData = migrationData
}

func (t *GormTestSuite) SetupSuite() {
	t.Context = context.Background()

	// Init database
	if t.TestCaseName == "" && t != nil && t.T() != nil && t.T().Name() != "" {
		t.SetTestCaseName(t.T().Name())
	}

	t.TestDbName = t.DatabaseName()
	postgresDBName := t.TestDbName

	dbPassword := utils.GetEnvDefault("POSTGRES_PASSWORD", "postgres")
	postgresHost := utils.GetEnvDefault("POSTGRES_HOST", "localhost")
	postgresUser := utils.GetEnvDefault("POSTGRES_USER", "postgres")
	postgresPort := utils.GetEnvDefault("POSTGRES_PORT", "5432")
	dbConfig := common_gorm.DBConfig{
		DB_HOST:                      postgresHost,
		DB_USERNAME:                  postgresUser,
		DB_NAME:                      postgresDBName,
		DB_PORT:                      postgresPort,
		DB_PASSWORD:                  dbPassword,
		DB_POOL_MAX_OPEN_CONNECTIONS: 50,
		DB_POOL_MAX_IDLE_CONNECTIONS: 10,
		DB_METRICS_ENABLED:           false,

		SSLMode:  "disable",
		TimeZone: "Asia/Ho_Chi_Minh",
	}

	if pgConn == nil {
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/postgres?sslmode=%s",
			postgresUser,
			dbPassword,
			postgresHost,
			postgresPort,
			dbConfig.SSLMode,
		)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: log.NewGormLogger("info"),
		})
		if err != nil {
			t.Fail("failed to connect database")
		}
		pgConn = db
	}

	if dbConn != nil {
		if conn, ok := dbConn[postgresDBName]; ok {
			t.Db = conn
			return
		}
	}

	// Check if the database already exists
	sqlCheckExists := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s';", postgresDBName)
	var exists int
	err := pgConn.Raw(sqlCheckExists).Scan(&exists).Error
	if err != nil {
		t.Fail("failed to check if database exists: %v", err)
		return
	}

	// Create the database if it doesn't exist
	if exists == 0 {
		sqlStatement := fmt.Sprintf("CREATE DATABASE %s;", postgresDBName)
		if err := pgConn.Exec(sqlStatement).Error; err != nil {
			t.Fail("failed to create database: %v", err)
		}
	}

	// connect test db
	dbObj, err := common_gorm.NewDatabase(&dbConfig)
	if err != nil {
		t.Fail("error: %v", err)
		return
	}

	t.Migrator = dbObj.Migrator

	// run migration
	if t.MigrationData != nil {
		migrationErr := dbObj.RunMigration(context.Background(), t.MigrationData)
		require.NoError(t.T(), migrationErr)
	}

	t.dbManager = dbObj

	dbConn[postgresDBName] = dbObj.GormDB
	t.Db = dbConn[postgresDBName]
}

func (t *GormTestSuite) TearDownSuite() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			t.Fail("error: %v", err)
		}

		// Drop test database
		sqlStatement := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE);", t.TestDbName)
		if err := pgConn.Exec(sqlStatement).Error; err != nil {
			t.Fail("failed to drop database: %v", err)
		}
	}()

	// auto migration
	if t.MigrationData != nil {
		migrationErr := t.dbManager.DownMigration(context.Background(), t.MigrationData)
		if migrationErr != nil {
			t.Fail("down migration failed: %v", migrationErr)
			return
		}
	}

	sqlDB, _ := t.Db.DB()
	err := sqlDB.Close()
	if err != nil {
		t.Fail("cannot close test database connection: %v", err)
	}

	// Drop test database
	sqlStatement := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE);", t.TestDbName)
	if err := pgConn.Exec(sqlStatement).Error; err != nil {
		t.Fail("failed to drop database: %v", err)
	}
}

func (t *GormTestSuite) SetupTest() {
	t.Tx = t.Db.Begin()

	if t.Tx.Error != nil {
		t.Fail("failed to start transaction for creating test data: %v", t.Tx.Error)
	}

	t.Data = make(map[string]interface{})
	t.Context = context.WithValue(t.Context, "tx", t.Tx)
}

func (t *GormTestSuite) TearDownTest() {
	err := t.Tx.Rollback().Error
	if err != nil {
		t.Fail(fmt.Sprintf("failed to clean up test data %s", err.Error()))
	}
}
