package zen

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	common_gorm "github.com/ngtrvu/zen-go/gorm"
	"github.com/ngtrvu/zen-go/log"
	"github.com/ngtrvu/zen-go/utils"
	"gorm.io/gorm"
)

const (
	IgnoreConflictCtx = "ignoreConflict"

	DefaultTimeout = 5 * time.Second
)

type RepoInterface interface {
	WithScopes(scopes ...Scope) RepoInterface
	GetAll(ctx context.Context, query *common_gorm.Query, items interface{}) error
	GetCount(ctx context.Context, query *common_gorm.Query, items interface{}) (int, error)
	Get(ctx context.Context, query interface{}, item interface{}) error
	GetByUUID(ctx context.Context, id uuid.UUID, item interface{}) error
	Create(ctx context.Context, item interface{}) error
	FirstOrCreate(ctx context.Context, item interface{}, conditions interface{}) error
	CreateMany(ctx context.Context, items interface{}) error
	Update(ctx context.Context, item interface{}) error
	UpdateOrCreate(ctx context.Context, query interface{}, getItem interface{}, updateItem interface{}) error
	BulkUpdateOrCreate(ctx context.Context, items interface{}, columns []clause.Column, updateColumns []string) error
	UpdatePartial(ctx context.Context, item interface{}, params map[string]interface{}) error
	UpdateLocking(ctx context.Context, item interface{}, params map[string]interface{}) error
	Delete(ctx context.Context, item interface{}) error
	GroupByField(ctx context.Context, db *gorm.DB, field string, fieldCount string, result interface{}) error
	GetDB(ctx context.Context) *gorm.DB
	StartTransaction(ctx context.Context) *gorm.DB
	RollbackTransaction(ctx context.Context) *gorm.DB
	CommitTransaction(ctx context.Context) *gorm.DB
	DeleteMany(ctx context.Context, condition interface{}) error
	BuildJoins(items interface{}, query *common_gorm.Query) []string
}

type Repo struct {
	db        *gorm.DB
	Scopes    []Scope
	SavePoint string
}

// Repo creates a new Repo instance.
func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) WithScopes(scopes ...Scope) RepoInterface {
	newRepo := &Repo{db: r.db, Scopes: r.Scopes}
	newRepo.Scopes = append(newRepo.Scopes, scopes...)
	return newRepo
}

func (r *Repo) GetAll(ctx context.Context, query *common_gorm.Query, items interface{}) (err error) {
	if query == nil {
		panic("query is required")
	}

	// Filter out invalid fields from query.Filters
	tableColumns := common_gorm.GetGormColumnNames(items)
	tableName := common_gorm.GetTableName(items)
	if len(tableColumns) > 0 {
		validFilters := make([]*common_gorm.FilterAttribute, 0)
		for _, filterItem := range query.Filter.Filters {
			fieldNameWithoutPrefix := strings.Replace(filterItem.Field, fmt.Sprintf("%s.", tableName), "", 1)
			if utils.Contains(tableColumns, fieldNameWithoutPrefix) {
				validFilters = append(validFilters, filterItem)
			}
		}

		query.Filter.Filters = validFilters
	}

	// Set a timeout for the GORM query
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	for _, s := range r.Scopes {
		db = db.Scopes(s)
	}

	q := db.Offset(query.Offset)

	// add limit if any
	if query.Limit > 0 {
		q = q.Limit(query.Limit)
	}

	// add order if any
	q = q.Order(query.SortStatement())

	// join fk tables
	joinQueries := r.BuildJoins(items, query)
	for _, joinQuery := range joinQueries {
		q = q.Joins(joinQuery)
	}

	// join query
	q = attachJoinQuery(items, query, q)
	attachPrefixTableNameToQueryField(items, query)

	queryStr, queryParams := query.Filter.QueryStatement()
	q = q.Where(queryStr, queryParams...)

	queryStr, queryParams = query.Search.QueryStatement()
	q = q.Where(queryStr, queryParams...)

	err = q.Find(items).Error

	return err
}

func (r *Repo) GetCount(ctx context.Context, query *common_gorm.Query, items interface{}) (int, error) {
	if query == nil {
		panic("query is required")
	}

	// Filter out invalid fields from query.Filters
	tableColumns := common_gorm.GetGormColumnNames(items)
	tableName := common_gorm.GetTableName(items)
	if len(tableColumns) > 0 {
		validFilters := make([]*common_gorm.FilterAttribute, 0)
		for _, filterItem := range query.Filter.Filters {
			fieldNameWithoutPrefix := strings.Replace(filterItem.Field, fmt.Sprintf("%s.", tableName), "", 1)
			if utils.Contains(tableColumns, fieldNameWithoutPrefix) {
				validFilters = append(validFilters, filterItem)
			}
		}

		query.Filter.Filters = validFilters
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	q := r.db.WithContext(timeoutCtx)

	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		q = tx
	}

	// join fk tables
	joinedTables := r.BuildJoins(items, query)
	for _, joinQuery := range joinedTables {
		q = q.Joins(joinQuery)
	}

	// join query
	q = attachJoinQuery(items, query, q)
	attachPrefixTableNameToQueryField(items, query)

	// join fk tables
	queryStr, queryParams := query.Filter.QueryStatement()
	q = q.Where(queryStr, queryParams...)

	queryStr, queryParams = query.Search.QueryStatement()
	q = q.Where(queryStr, queryParams...)

	var count int64
	q.Model(items).Count(&count)
	return int(count), nil
}

func (r *Repo) BuildJoins(items interface{}, query *common_gorm.Query) []string {
	tableName := common_gorm.GetTableName(items)
	if tableName == "" {
		log.Error("table name is empty")
		return nil
	}

	joinedTables := make(map[string]bool, 0)
	queries := make([]string, 0)

	for _, searchField := range query.Search.SearchFields {
		if strings.Contains(searchField.Field, ".") {
			parts := strings.Split(searchField.Field, ".")
			if len(parts) == 2 {
				joinTable := parts[0]
				if joinTable == tableName {
					continue
				}

				if searchField.JoinedColumn == "" {
					log.Warn("joined column is empty")
					continue
				}

				if _, ok := joinedTables[joinTable]; ok {
					continue
				}

				queries = append(
					queries,
					fmt.Sprintf(
						"LEFT JOIN %s ON %s.id = %s.%s",
						joinTable,
						joinTable,
						tableName,
						searchField.JoinedColumn,
					),
				)
				joinedTables[joinTable] = true
			}
		}
	}

	return queries
}

func attachJoinQuery(items interface{}, query *common_gorm.Query, q *gorm.DB) *gorm.DB {
	tableName := common_gorm.GetTableName(items)
	if tableName == "" {
		return q
	}

	joinedTables := make(map[string]bool, 0)
	if len(query.Filter.Filters) > 0 {
		for _, filterField := range query.Filter.Filters {
			if strings.Contains(filterField.Field, ".") {
				parts := strings.Split(filterField.Field, ".")
				if len(parts) == 2 {
					joinTable := parts[0]
					fkField := common_gorm.GetForeignKeyField(items, joinTable)

					if _, ok := joinedTables[joinTable]; ok {
						continue
					}
					joinedTables[joinTable] = true
					if fkField != "" {
						q = q.Joins(
							fmt.Sprintf("LEFT JOIN %s ON %s.id = %s.%s", joinTable, joinTable, tableName, fkField),
						)
					}
				}
			}
		}
	}

	return q
}

func attachPrefixTableNameToQueryField(
	items interface{},
	query *common_gorm.Query,
) error {
	newFilter := common_gorm.Filter{}
	newSearch := common_gorm.Search{}

	tableName := common_gorm.GetTableName(items)
	if tableName == "" {
		return errors.New("table name is empty")
	}

	if len(query.Filter.Filters) > 0 {
		for _, filterField := range query.Filter.Filters {
			if strings.Contains(filterField.Field, ".") {
				newFilter.AddFilter(
					&common_gorm.FilterAttribute{
						Field:    filterField.Field,
						Type:     filterField.Type,
						Operator: filterField.Operator,
						Value:    filterField.Value,
					},
				)
			} else {
				newFilter.AddFilter(
					&common_gorm.FilterAttribute{
						Field:    fmt.Sprintf("%s.%s", tableName, filterField.Field),
						Type:     filterField.Type,
						Operator: filterField.Operator,
						Value:    filterField.Value,
					},
				)
			}
		}
	}

	query.Filter = newFilter

	for _, searchField := range query.Search.SearchFields {
		if strings.Contains(searchField.Field, ".") {
			newSearch.AddSearchField(
				&common_gorm.SearchAttribute{
					Field:        searchField.Field,
					JoinedColumn: searchField.JoinedColumn,
					Type:         searchField.Type,
					Operator:     searchField.Operator,
					Value:        searchField.Value,
				},
			)
		} else {
			newSearch.AddSearchField(
				&common_gorm.SearchAttribute{
					Field:    fmt.Sprintf("%s.%s", tableName, searchField.Field),
					Type:     searchField.Type,
					Operator: searchField.Operator,
					Value:    searchField.Value,
				},
			)
		}
	}

	query.Search = newSearch

	return nil
}

func (r *Repo) Get(ctx context.Context, query interface{}, item interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	for _, s := range r.Scopes {
		db = db.Scopes(s)
	}

	return db.Where(query).First(item).Error
}

func (r *Repo) GetByUUID(ctx context.Context, id uuid.UUID, item interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	for _, s := range r.Scopes {
		db = db.Scopes(s)
	}

	return db.Where("id = ?", id).First(item).Error
}

func (r *Repo) Create(ctx context.Context, item interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	ignoreConflict, _ := ctx.Value("ignoreConflict").(bool)
	if ignoreConflict {
		return db.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error
	}

	return db.Create(item).Error
}

func (r *Repo) FirstOrCreate(ctx context.Context, item interface{}, conditions interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	return db.FirstOrCreate(item, conditions).Error
}

func (r *Repo) CreateMany(ctx context.Context, items interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	ignoreConflict, _ := ctx.Value("ignoreConflict").(bool)
	if ignoreConflict {
		return db.Clauses(clause.OnConflict{DoNothing: true}).Create(items).Error
	}

	return db.Create(items).Error
}

func (r *Repo) Update(ctx context.Context, item interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	return db.Save(item).Error
}

func (r *Repo) UpdatePartial(ctx context.Context, item interface{}, params map[string]interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	return db.Model(item).Updates(params).Error
}

func (r *Repo) UpdateLocking(ctx context.Context, item interface{}, params map[string]interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	return db.Clauses(clause.Locking{Strength: "UPDATE"}).Model(item).Updates(params).Error
}

func (r *Repo) Delete(ctx context.Context, item interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	q := db.Delete(item)
	if q.Error != nil {
		return q.Error
	}
	return nil
}

// Implement DeleteMany method
func (r *Repo) DeleteMany(ctx context.Context, condition interface{}) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	q := db.Where(condition).Delete(nil)
	if q.Error != nil {
		return q.Error
	}

	return nil
}

func (r *Repo) GroupByField(
	ctx context.Context,
	db *gorm.DB,
	field string,
	fieldCount string,
	result interface{},
) error {
	return db.Select(fmt.Sprintf("%s, count(distinct %s)", field, fieldCount)).Group(field).Find(result).Error
}

func (r *Repo) GetDB(ctx context.Context) *gorm.DB {
	db := r.db
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	return db
}

func (r *Repo) StartTransaction(ctx context.Context) *gorm.DB {
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		r.SavePoint = "xyz"
		return tx.SavePoint(r.SavePoint)
	}
	return r.db.Begin()
}

func (r *Repo) RollbackTransaction(ctx context.Context) *gorm.DB {
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		if r.SavePoint != "" {
			tx = tx.RollbackTo(r.SavePoint)
			r.SavePoint = ""
			return tx
		} else {
			return tx.Rollback()
		}
	}

	panic("RollbackTransaction failed, not in any transaction")
}

func (r *Repo) CommitTransaction(ctx context.Context) *gorm.DB {
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		if r.SavePoint == "" {
			return tx.Commit()
		}

		r.SavePoint = ""
		return tx
	}

	panic("CommitTransaction failed, not in any transaction")
}

func (r *Repo) UpdateOrCreate(
	ctx context.Context,
	query interface{},
	getItem interface{},
	updateItem interface{},
) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	return db.Where(query).Assign(updateItem).FirstOrCreate(getItem).Error
}

func (r *Repo) BulkUpdateOrCreate(
	ctx context.Context,
	items interface{},
	columns []clause.Column,
	updateColumns []string,
) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	db := r.db.WithContext(timeoutCtx)
	tx, inTransaction := ctx.Value("tx").(*gorm.DB)
	if inTransaction {
		db = tx
	}

	return db.Clauses(clause.OnConflict{
		Columns:   columns,
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(items).Error
}
