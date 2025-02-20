package zen_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	common_gorm "github.com/ngtrvu/zen-go/gorm"
	"github.com/ngtrvu/zen-go/gormtest"
	"github.com/ngtrvu/zen-go/utils"
	"github.com/ngtrvu/zen-go/zen"
)

type DBConfig struct {
	DbHost             string `config:"POSTGRES_HOST"`
	DbPort             string `config:"POSTGRES_PORT"`
	DbUser             string `config:"POSTGRES_USER"`
	DbPassword         string `config:"POSTGRES_PASSWORD"`
	DbName             string `config:"POSTGRES_DB"`
	DbNameTest         string `config:"POSTGRES_DB_TEST"`
	PoolMaxConnections int    `config:"DB_POOL_MAX_CONNECTIONS"`
	SSLMode            string `config:"DB_SSL_MODE"`
	TimeZone           string `config:"DB_TIME_ZONE"`
}

type PostItem struct {
	ID          int       `gorm:"column:id"                           json:"id"`
	Name        string    `gorm:"column:name"                         json:"name"`
	Email       string    `gorm:"uniqueIndex:uidx_email,column:email" json:"email"`
	Age         int       `gorm:"column:age"                          json:"age"`
	CategoryID  int       `gorm:"column:categories_id"                json:"categories_id"`
	Category    Category  `gorm:"foreignKey:CategoryID"               json:"categories"`
	Category2ID int       `gorm:"column:categories_2_id"              json:"categories_2_id"`
	Category2   Category  `gorm:"foreignKey:Category2ID"              json:"categories_2"`
	TagID       int       `gorm:"column:tag_id"                       json:"tag_id"`
	Tag         Tag       `gorm:"foreignKey:TagID"                    json:"tags"`
	CreatedAt   time.Time `gorm:"column:created_at"                   json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"                   json:"updated_at"`
}

type Category struct {
	ID        int       `gorm:"column:id"         json:"id"`
	Name      string    `gorm:"column:name"       json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Tag struct {
	ID        int       `gorm:"column:id"         json:"id"`
	Name      string    `gorm:"column:name"       json:"name"`
	Title     string    `gorm:"column:title"      json:"title"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PostItem) TableName() string {
	return "post_items"
}

func (Category) TableName() string {
	return "categories"
}

func (Tag) TableName() string {
	return "tags"
}

type RepoTestSuite struct {
	gormtest.GormTestSuite
}

func TestRepoTestSuite(t *testing.T) {
	s := new(RepoTestSuite)
	suite.Run(t, s)
}

func createTestDb(dbConfig *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		dbConfig.DbUser,
		dbConfig.DbPassword,
		dbConfig.DbHost,
		dbConfig.DbPort,
		dbConfig.DbName,
	)
	// Connect to the default database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Create the test database if it does not exist
	sql := fmt.Sprintf(
		"SELECT 'CREATE DATABASE %s' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s');",
		dbConfig.DbName,
		dbConfig.DbName,
	)
	err = db.Exec(sql).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}

	// Connect to the test database
	dsn = fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		dbConfig.DbUser,
		dbConfig.DbPassword,
		dbConfig.DbHost,
		dbConfig.DbPort,
		dbConfig.DbNameTest,
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %v", err)
	}

	// Drop table post_items
	sql = `DROP TABLE IF EXISTS post_items;`
	err = db.Exec(sql).Error
	if err != nil {
		return nil, fmt.Errorf("failed to drop table post_items: %v", err)
	}

	// Create table post_items
	sql = `CREATE TABLE IF NOT EXISTS post_items (id SERIAL PRIMARY KEY, name VARCHAR(255),
	age INT, created_at timestamp without time zone default (now() at time zone 'utc'),
	updated_at timestamp without time zone default (now() at time zone 'utc')); 
	CREATE UNIQUE INDEX uidx_email ON post_items (email);
	`
	err = db.Exec(sql).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create table post_items: %v", err)
	}

	return db, nil
}

func (t *RepoTestSuite) CreateTables() {
	db := t.Tx

	// Drop table post_items
	sql := `DROP TABLE IF EXISTS post_items;`
	err := db.Exec(sql).Error
	if err != nil {
		panic(err)
	}

	// Create table post_items
	sql = `
	CREATE EXTENSION IF NOT EXISTS unaccent;
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		created_at timestamp without time zone default (now() at time zone 'utc'),
		updated_at timestamp without time zone default (now() at time zone 'utc')
	);

	CREATE TABLE IF NOT EXISTS tags (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		title VARCHAR(255),
		created_at timestamp without time zone default (now() at time zone 'utc'),
		updated_at timestamp without time zone default (now() at time zone 'utc')
	);

	CREATE TABLE IF NOT EXISTS post_items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		email VARCHAR(255),
		age INT,
		categories_id INT,
		categories_2_id INT,
		tag_id INT,
		created_at timestamp without time zone default (now() at time zone 'utc'),
		updated_at timestamp without time zone default (now() at time zone 'utc')
	);

	CREATE UNIQUE INDEX uidx_email ON post_items (email);
	`

	err = db.Exec(sql).Error
	if err != nil {
		panic(err)
	}
}

func initTestDb() (*gorm.DB, error) {
	// Init database
	dbPassword := utils.GetEnvDefault("POSTGRES_PASSWORD", "postgres")
	postgresHost := utils.GetEnvDefault("POSTGRES_HOST", "localhost")
	postgresUser := utils.GetEnvDefault("POSTGRES_USER", "postgres")
	postgresDBName := utils.GetEnvDefault("POSTGRES_DB", "test_stag_1")
	postgresPort := utils.GetEnvDefault("POSTGRES_PORT", "5432")

	dbConfig := DBConfig{
		DbName:             postgresDBName,
		DbHost:             postgresHost,
		DbPort:             postgresPort,
		DbUser:             postgresUser,
		DbPassword:         dbPassword,
		PoolMaxConnections: 10,
		SSLMode:            "disable",
	}

	db, err := createTestDb(&dbConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestRepo_Create(t *testing.T) {
	t.Skip("skipping testing in CI environment")

	db, err := initTestDb()
	if err != nil {
		t.Fatalf("create database test error: %v", err)
		return
	}

	// Create a new Repo instance
	repo := zen.NewRepo(db)

	// Create a new PostItem instance
	item := &PostItem{
		Name: "John",
		Age:  30,
	}

	// Insert the PostItem instance into the database
	err = repo.Create(context.Background(), item)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	// Check if the record was created
	var item2 PostItem
	err = db.First(&item2).Error
	if err != nil || item2.ID == 0 {
		t.Fatalf("failed to get record: %v", err)
	}
}

func TestRepo_CreateMany(t *testing.T) {
	t.Skip("skipping testing in CI environment")

	db, err := initTestDb()
	if err != nil {
		t.Fatalf("create database test error: %v", err)
		return
	}

	// Create a new Repo instance
	repo := zen.NewRepo(db)

	// Create a new PostItem instance
	myTables := []*PostItem{
		{
			Name: "John",
			Age:  30,
		},
		{
			Name: "Jenny",
			Age:  20,
		},
		{
			Name: "Tom",
			Age:  35,
		},
	}

	// Insert the PostItem instance into the database
	err = repo.CreateMany(context.Background(), myTables)
	if err != nil {
		t.Fatalf("failed to create records: %v", err)
	}

	// Check if the record was created
	var myTables2 []*PostItem
	err = db.Find(&myTables2).Error
	if err != nil || len(myTables2) == 0 {
		t.Fatalf("failed to get records: %v", err)
	}
}

func TestRepo_UpdatePartial(t *testing.T) {
	t.Skip("skipping testing in CI environment")

	db, err := initTestDb()
	if err != nil {
		t.Fatalf("create database test error: %v", err)
		return
	}

	// Create a new Repo instance
	repo := zen.NewRepo(db)

	// Create a new PostItem instance
	item := &PostItem{
		Name: "John",
		Age:  30,
	}

	// Insert the PostItem instance into the database
	err = repo.Create(context.Background(), item)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	// Check if the record was created
	var item2 PostItem
	err = db.First(&item2).Error
	if err != nil || item2.ID == 0 {
		t.Fatalf("failed to get record: %v", err)
	}

	// Update the record
	err = repo.UpdatePartial(context.Background(), &item2, map[string]interface{}{
		"name": "John Doe",
		"age":  40,
	})
	if err != nil {
		t.Fatalf("failed to update record: %v", err)
	}

	// Load the record again
	var item3 PostItem
	err = db.First(&item3).Error
	if err != nil || item3.ID == 0 {
		t.Fatalf("failed to get record: %v", err)
	}

	assert.Equal(t, "John Doe", item3.Name)
	assert.Equal(t, 40, item3.Age)
}

func TestRepo_Get(t *testing.T) {
	t.Skip("skipping testing in CI environment")

	db, err := initTestDb()
	if err != nil {
		t.Fatalf("create database test error: %v", err)
		return
	}

	// Create a new Repo instance
	repo := zen.NewRepo(db)

	// Create a new PostItem instance
	item := &PostItem{
		Name: "John",
		Age:  30,
	}

	// Insert the PostItem instance into the database
	err = repo.Create(context.Background(), item)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	// Check if the record was created
	var item2 PostItem
	err = repo.Get(context.Background(), &PostItem{Name: "John2"}, &item2)
	assert.Error(t, err, "failed to get record")

	err = repo.Get(context.Background(), &PostItem{Name: "John"}, &item2)
	assert.NoError(t, err)
	assert.Equal(t, 30, item2.Age)
}

func (t *RepoTestSuite) TestRepo_GetAll() {
	ctx := t.Context
	t.CreateTables()

	// Create a new Repo instance
	repo := zen.NewRepo(t.Tx)

	tags := []*Tag{
		{
			Name:  "X1",
			Title: "Title 1",
		},
		{
			Name:  "X2",
			Title: "Title 2",
		},
		{
			Name:  "X3",
			Title: "Title 3",
		},
	}
	err := repo.CreateMany(ctx, tags)
	assert.NoError(t.T(), err)

	// Create a new PostItem instance
	post_items := []*PostItem{
		{
			Name:  "John",
			Email: "jon@dev",
			TagID: tags[0].ID,
		},
		{
			Name:  "John Depp",
			Email: "jon1@dev",
			TagID: tags[0].ID,
		},
		{
			Name:  "Jenny",
			Email: "jen@stagg",
			TagID: tags[1].ID,
		},
		{
			Name:  "Tom",
			Email: "tom@stagg",
			TagID: tags[2].ID,
		},
	}

	err = repo.CreateMany(ctx, post_items)
	assert.NoError(t.T(), err)

	// Check if the record was created
	search := &common_gorm.Search{
		SearchFields: []*common_gorm.SearchAttribute{
			{
				Field:        "tags.name",
				JoinedColumn: "tag_id",
				Type:         common_gorm.FieldTypeString,
				Operator:     common_gorm.OperatorLike,
				Value:        "X1",
			},
			{
				Field:        "tags.title",
				JoinedColumn: "tag_id",
				Type:         common_gorm.FieldTypeString,
				Operator:     common_gorm.OperatorLike,
				Value:        "Title 1",
			},
			{
				Field:    "email",
				Type:     common_gorm.FieldTypeString,
				Operator: common_gorm.OperatorLike,
				Value:    "@dev",
			},
		},
	}

	limit := 24

	var items []*PostItem
	query := &common_gorm.Query{Limit: limit, Search: *search}
	err = repo.GetAll(ctx, query, &items)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), 2, len(items))

	count, err := repo.GetCount(ctx, query, &items)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), 2, count)

	filter := &common_gorm.Filter{}
	filter.AddFilter(&common_gorm.FilterAttribute{
		Field:    "email",
		Type:     common_gorm.FieldTypeString,
		Operator: common_gorm.OperatorEqual,
		Value:    "jon1@dev",
	})

	query = &common_gorm.Query{Limit: limit, Filter: *filter}
	err = repo.GetAll(ctx, query, &items)
	assert.NoError(t.T(), err)

	assert.Equal(t.T(), 1, len(items))
}

func (t *RepoTestSuite) TestRepo_BuildJoins() {
	query := common_gorm.Query{
		Search: common_gorm.Search{
			SearchFields: []*common_gorm.SearchAttribute{
				{
					Field:        "categories.name",
					JoinedColumn: "categories_id",
					Type:         common_gorm.FieldTypeString,
					Operator:     common_gorm.OperatorLike,
					Value:        "xyz",
				},
				{
					Field:        "categories.email",
					JoinedColumn: "categories_id",
					Type:         common_gorm.FieldTypeString,
					Operator:     common_gorm.OperatorLike,
					Value:        "xyz",
				},
				{
					Field:        "tags.name",
					JoinedColumn: "tag_id",
					Type:         common_gorm.FieldTypeString,
					Operator:     common_gorm.OperatorLike,
					Value:        "xyz",
				},
			},
		},
	}

	items := PostItem{}
	repo := zen.NewRepo(t.Tx)
	result := repo.BuildJoins(items, &query)
	assert.Equal(t.T(), []string{
		"LEFT JOIN categories ON categories.id = post_items.categories_id",
		"LEFT JOIN tags ON tags.id = post_items.tag_id",
	}, result)

	result = repo.BuildJoins(&items, &query)
	assert.Equal(t.T(), []string{
		"LEFT JOIN categories ON categories.id = post_items.categories_id",
		"LEFT JOIN tags ON tags.id = post_items.tag_id",
	}, result)

	items2 := []*PostItem{}
	result = repo.BuildJoins(&items2, &query)
	assert.Equal(t.T(), []string{
		"LEFT JOIN categories ON categories.id = post_items.categories_id",
		"LEFT JOIN tags ON tags.id = post_items.tag_id",
	}, result)

	items3 := utils.CreateArrayFromObject(items2)
	result = repo.BuildJoins(items3, &query)
	assert.Equal(t.T(), []string{
		"LEFT JOIN categories ON categories.id = post_items.categories_id",
		"LEFT JOIN tags ON tags.id = post_items.tag_id",
	}, result)

	result = repo.BuildJoins(&items3, &query)
	assert.Equal(t.T(), []string{
		"LEFT JOIN categories ON categories.id = post_items.categories_id",
		"LEFT JOIN tags ON tags.id = post_items.tag_id",
	}, result)

	// TODO: support for multiple FK to a same table
}

func (t *RepoTestSuite) TestRepo_StartTransaction() {
	ctx := t.Context
	t.CreateTables()

	// Create a new Repo instance
	repo := zen.NewRepo(t.Tx)

	// Create data
	t.createDataWithTransaction(ctx, repo)

	// Check if the record was created
	search := &common_gorm.Search{
		SearchFields: []*common_gorm.SearchAttribute{
			{
				Field:        "tags.name",
				JoinedColumn: "tag_id",
				Type:         common_gorm.FieldTypeString,
				Operator:     common_gorm.OperatorLike,
				Value:        "X1",
			},
			{
				Field:        "tags.title",
				JoinedColumn: "tag_id",
				Type:         common_gorm.FieldTypeString,
				Operator:     common_gorm.OperatorLike,
				Value:        "Title 1",
			},
			{
				Field:    "email",
				Type:     common_gorm.FieldTypeString,
				Operator: common_gorm.OperatorLike,
				Value:    "@dev",
			},
		},
	}

	limit := 24

	var items []*PostItem
	query := &common_gorm.Query{Limit: limit, Search: *search}
	err := repo.GetAll(ctx, query, &items)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), 2, len(items))

	count, err := repo.GetCount(ctx, query, &items)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), 2, count)

	filter := &common_gorm.Filter{}
	filter.AddFilter(&common_gorm.FilterAttribute{
		Field:    "email",
		Type:     common_gorm.FieldTypeString,
		Operator: common_gorm.OperatorEqual,
		Value:    "jon1@dev",
	})

	query = &common_gorm.Query{Filter: *filter}
	err = repo.GetAll(ctx, query, &items)
	assert.NoError(t.T(), err)

	assert.Equal(t.T(), 1, len(items))
}

func (t *RepoTestSuite) createDataWithTransaction(ctx context.Context, repo *zen.Repo) (err error) {
	tags, err := t.createTag(ctx, repo)
	assert.NoError(t.T(), err)

	_, err = t.createMyTable(ctx, repo, tags)
	assert.NoError(t.T(), err)
	return
}

func (t *RepoTestSuite) createTag(ctx context.Context, repo *zen.Repo) (tags []*Tag, err error) {
	// Start a transaction
	tx := repo.StartTransaction(ctx)
	ctx = context.WithValue(ctx, "tx", tx)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}

		if err != nil {
			repo.RollbackTransaction(ctx)
		} else {
			repo.CommitTransaction(ctx)
		}

		// Savepoint should be empty after the transaction is committed/rollbacked
		assert.Equal(t.T(), "", repo.SavePoint)
	}()

	tags = []*Tag{
		{
			Name:  "X1",
			Title: "Title 1",
		},
		{
			Name:  "X2",
			Title: "Title 2",
		},
		{
			Name:  "X3",
			Title: "Title 3",
		},
	}

	err = repo.CreateMany(ctx, tags)
	return
}

func (t *RepoTestSuite) createMyTable(
	ctx context.Context,
	repo *zen.Repo,
	tags []*Tag,
) (items []*PostItem, err error) {
	// Start a transaction
	tx := repo.StartTransaction(ctx)
	ctx = context.WithValue(ctx, "tx", tx)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}

		if err != nil {
			repo.RollbackTransaction(ctx)
		} else {
			repo.CommitTransaction(ctx)
		}

		// Savepoint should be empty after the transaction is committed/rollbacked
		assert.Equal(t.T(), "", repo.SavePoint)
	}()

	// Create a new PostItem instance
	items = []*PostItem{
		{
			Name:  "John",
			Email: "jon@dev",
			TagID: tags[0].ID,
		},
		{
			Name:  "John Depp",
			Email: "jon1@dev",
			TagID: tags[0].ID,
		},
		{
			Name:  "Jenny",
			Email: "jen@stagg",
			TagID: tags[1].ID,
		},
		{
			Name:  "Tom",
			Email: "tom@stagg",
			TagID: tags[2].ID,
		},
	}

	err = repo.CreateMany(ctx, items)
	assert.NoError(t.T(), err)
	return
}

func (t *RepoTestSuite) TestBulkUpdateOrCreate() {
	ctx := t.Context
	t.CreateTables()

	// Create a new Repo instance
	repo := zen.NewRepo(t.Tx)

	// Create data
	t.createDataWithTransaction(ctx, repo)

	// Check if the record was created
	var items []*PostItem
	err := repo.GetAll(ctx, &common_gorm.Query{}, &items)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), 4, len(items))

	// Bulk update or create
	newItems := []*PostItem{
		{
			ID:    items[0].ID,
			Name:  "John",
			Email: "jon555@dev",
		},
		{
			ID:    99999,
			Name:  "John",
			Email: "jon99999@dev",
		},
	}

	err = repo.BulkUpdateOrCreate(ctx, newItems, []clause.Column{{Name: "id"}}, []string{"email", "name"})
	assert.NoError(t.T(), err)

	repo.GetAll(ctx, &common_gorm.Query{}, &items)
	assert.Equal(t.T(), 5, len(items))

	items = []*PostItem{{
		Name:  "John",
		Email: "jon@dev",
	}}
	err = repo.BulkUpdateOrCreate(ctx, items, []clause.Column{{Name: "email"}}, []string{"email", "name"})
	assert.NoError(t.T(), err)
}
