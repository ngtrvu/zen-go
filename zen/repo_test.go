package zen

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	common_gorm "github.com/ngtrvu/zen-go/gorm"
)

// TestModel is a test model for testing
type TestModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TypeModelID uuid.UUID `gorm:"type:uuid"`
	TypeModel   TypeModel `gorm:"foreignKey:TypeModelID"`
}

type TypeModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&TestModel{})
	assert.NoError(t, err)

	return db
}

func TestNewRepo(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestWithScopes(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", "test")
	}

	newRepo := repo.WithScopes(scope)
	assert.NotNil(t, newRepo)
	assert.Len(t, newRepo.(*Repo).Scopes, 1)
}

func TestGetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	// Create test data
	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}
	err := db.Create(testModel).Error
	assert.NoError(t, err)

	// Test GetAll
	query := &common_gorm.Query{
		Filter: common_gorm.Filter{
			Filters: []*common_gorm.FilterAttribute{
				{
					Field:    "name",
					Type:     "string",
					Operator: "=",
					Value:    "test",
				},
			},
		},
	}

	var results []TestModel
	err = repo.GetAll(context.Background(), query, &results)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, testModel.Name, results[0].Name)
}

func TestGetCount(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	// Create test data
	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}
	err := db.Create(testModel).Error
	assert.NoError(t, err)

	query := &common_gorm.Query{
		Filter: common_gorm.Filter{
			Filters: []*common_gorm.FilterAttribute{
				{
					Field:    "name",
					Type:     "string",
					Operator: "=",
					Value:    "test",
				},
			},
		},
	}

	count, err := repo.GetCount(context.Background(), query, &TestModel{})
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestGet(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	// Create test data
	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}
	err := db.Create(testModel).Error
	assert.NoError(t, err)

	var result TestModel
	err = repo.Get(context.Background(), map[string]interface{}{"name": "test"}, &result)
	assert.NoError(t, err)
	assert.Equal(t, testModel.Name, result.Name)
}

func TestGetByUUID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	// Create test data
	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}
	err := db.Create(testModel).Error
	assert.NoError(t, err)

	var result TestModel
	err = repo.GetByUUID(context.Background(), testModel.ID, &result)
	assert.NoError(t, err)
	assert.Equal(t, testModel.ID, result.ID)
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}

	err := repo.Create(context.Background(), testModel)
	assert.NoError(t, err)

	var result TestModel
	err = db.First(&result, "id = ?", testModel.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, testModel.Name, result.Name)
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	// Create test data
	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}
	err := db.Create(testModel).Error
	assert.NoError(t, err)

	// Update
	testModel.Name = "updated"
	err = repo.Update(context.Background(), testModel)
	assert.NoError(t, err)

	var result TestModel
	err = db.First(&result, "id = ?", testModel.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "updated", result.Name)
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	// Create test data
	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}
	err := db.Create(testModel).Error
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), testModel)
	assert.NoError(t, err)

	var result TestModel
	err = db.First(&result, "id = ?", testModel.ID).Error
	assert.Error(t, err)
}

func TestTransaction(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	ctx := context.Background()
	tx := repo.StartTransaction(ctx)
	assert.NotNil(t, tx)

	// Create a new context with the transaction
	txCtx := context.WithValue(ctx, "tx", tx)

	// Create test data in transaction
	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}
	err := tx.Create(testModel).Error
	assert.NoError(t, err)

	// Commit transaction using the transaction context
	tx = repo.CommitTransaction(txCtx)
	assert.NotNil(t, tx)

	// Verify data was committed
	var result TestModel
	err = db.First(&result, "id = ?", testModel.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, testModel.Name, result.Name)
}

func TestUpdateOrCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}

	// Test create
	err := repo.UpdateOrCreate(
		context.Background(),
		map[string]interface{}{"name": "test"},
		testModel,
		map[string]interface{}{"name": "test"},
	)
	assert.NoError(t, err)

	// Test update
	testModel.Name = "updated"
	err = repo.UpdateOrCreate(
		context.Background(),
		map[string]interface{}{"name": "test"},
		testModel,
		map[string]interface{}{"name": "updated"},
	)
	assert.NoError(t, err)

	var result TestModel
	err = db.First(&result, "id = ?", testModel.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "updated", result.Name)
}

func TestBulkUpdateOrCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	items := []TestModel{
		{ID: uuid.New(), Name: "test1"},
		{ID: uuid.New(), Name: "test2"},
	}

	columns := []clause.Column{{Name: "id"}}
	updateColumns := []string{"name"}

	err := repo.BulkUpdateOrCreate(context.Background(), items, columns, updateColumns)
	assert.NoError(t, err)

	var results []TestModel
	err = db.Find(&results).Error
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestPartialUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	// Create test data
	modelType := &TypeModel{
		ID:   uuid.New(),
		Name: "my model type",
	}
	modelType2 := &TypeModel{
		ID:   uuid.New(),
		Name: "my model type 2",
	}
	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}
	err := db.Create(modelType).Error
	assert.NoError(t, err)

	err = db.Create(testModel).Error
	assert.NoError(t, err)

	// Update partial
	err = repo.UpdatePartial(context.Background(), testModel, map[string]interface{}{"name": "updated"})
	assert.NoError(t, err)

	var result TestModel
	err = db.First(&result, "id = ?", testModel.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "updated", result.Name)

	// Update partial with foreign key
	err = repo.UpdatePartial(context.Background(), testModel, map[string]interface{}{"type_model_id": modelType.ID})
	assert.NoError(t, err)

	var result2 TestModel
	err = db.First(&result2, "id = ?", testModel.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, modelType.ID, result2.TypeModelID)

	// Update partial with foreign key
	err = repo.UpdatePartial(context.Background(), testModel, map[string]interface{}{"type_model_id": modelType2.ID})
	assert.NoError(t, err)

	var result3 TestModel
	err = db.First(&result3, "id = ?", testModel.ID).Error
	assert.NoError(t, err)
}

func TestUpdateLocking(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)

	// Create test data
	testModel := &TestModel{
		ID:   uuid.New(),
		Name: "test",
	}
	err := db.Create(testModel).Error
	assert.NoError(t, err)

	// Update with locking
	err = repo.UpdateLocking(context.Background(), testModel, map[string]interface{}{"name": "updated"})
	assert.NoError(t, err)

	var result TestModel
	err = db.First(&result, "id = ?", testModel.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "updated", result.Name)
}
