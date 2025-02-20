package migrator

import (
	"sort"

	"github.com/ngtrvu/zen-go/log"
	"gorm.io/gorm"
)

type Migrator struct {
	gorm.Migrator

	DB *gorm.DB
}

type DbTable struct {
	TableCatalog string `json:"table_catalog"`
	TableSchema  string `json:"table_schema"`
	TableName    string `json:"table_name"`
	TableType    string `json:"table_type"`
}

type DbColumn struct {
	TableCatalog           string `json:"table_catalog"`
	TableSchema            string `json:"table_schema"`
	TableName              string `json:"table_name"`
	OridinalPosition       string `json:"ordinal_position"`
	ColumnName             string `json:"column_name"`
	ColumnDefault          string `json:"column_default"`
	IsNullable             string `json:"is_nullable"`
	DataType               string `json:"data_type"`
	UdtName                string `json:"udt_name"`
	CharacterMaximumLength string `json:"character_maximum_length"`
	CharacterOctetLength   string `json:"character_octet_length"`
	NumericPrecision       string `json:"numeric_precision"`
	NumericScale           string `json:"numeric_scale"`
	DatetimePrecision      string `json:"datetime_precision"`
}

// Get column type of the column with specs, for example: character varying, timestamp with time zone, integer, smallint
func (t *DbColumn) ColumnType() string {
	switch t.DataType {
	case "character varying":
		return t.DataType + "(" + t.CharacterMaximumLength + ")"
	case "numeric":
		if t.NumericPrecision == "" || t.NumericScale == "" {
			return t.DataType
		}

		return t.DataType + "(" + t.NumericPrecision + "," + t.NumericScale + ")"
	}

	return t.DataType
}

func SortDbTables(tables []DbTable) {
	sort.Slice(tables, func(i, j int) bool {
		return tables[i].TableName < tables[j].TableName
	})
}

func SortDbColumns(columns []DbColumn) {
	sort.Slice(columns, func(i, j int) bool {
		return columns[i].ColumnName < columns[j].ColumnName
	})
}

func FindColumn(columns []gorm.ColumnType, columnName string) *gorm.ColumnType {
	for _, column := range columns {
		if column.Name() == columnName {
			return &column
		}
	}

	return nil
}

func FindTable(tables []DbTable, tableName string) *DbTable {
	for _, table := range tables {
		if table.TableName == tableName {
			return &table
		}
	}

	return nil
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		Migrator: db.Migrator(),
		DB:       db,
	}
}

func (m *Migrator) GetColumns(tableName string) ([]DbColumn, error) {
	var columns []DbColumn
	err := m.DB.Raw(`
		SELECT * 
		FROM information_schema.columns 
		WHERE table_name = $1 
		ORDER BY column_name asc
		`, tableName).
		Scan(&columns).
		Error
	if err != nil {
		log.Error("Failed to get all columns: %v", err)
		return nil, err
	}

	SortDbColumns(columns)
	return columns, nil
}

func (m *Migrator) GetTables() ([]DbTable, error) {
	var tables []DbTable
	err := m.DB.Raw(`
		SELECT * 
		FROM information_schema.tables 
		WHERE 
			table_schema = 'public' and table_name != 'schema_migrations'
		ORDER BY table_name asc
		`).
		Scan(&tables).
		Error
	if err != nil {
		log.Error("Failed to get all tables: %v", err)
		return nil, err
	}

	SortDbTables(tables)
	return tables, nil
}
