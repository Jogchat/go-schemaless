// Package mysql is a mysql-backed Schemaless store.
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"code.jogchat.internal/go-schemaless/models"
	"go.uber.org/zap"
	"time"
	"encoding/json"
	"code.jogchat.internal/go-schemaless/utils"
)

// Storage is a MySQL-backed storage.
type Storage struct {
	user     string
	pass     string
	host     string
	port     string
	database string

	store	*sql.DB
	indexes	map[string]*Index
	Sugar	*zap.SugaredLogger
}

const (
	//timeParseString = "2006-01-02T15:04:05Z"
	timeParseString  = "2006-01-02 15:04:05"
	driver = "mysql"
	// dsnFormat string parameters: username, password, host, port, database.
	// parseTime is for parsing and handling *time.Time properly
	dsnFormat = "%s:%s@tcp(%s:%s)/%s?parseTime=true"
	// This space intentionally left blank for facilitating vimdiff
	// acrosss storages.

	getCellSQL          = "SELECT added_at, row_key, column_name, ref_key, body,created_at FROM cell WHERE row_key = ? AND column_name = ? AND ref_key = ? LIMIT 1"
	getCellLatestSQL    = "SELECT added_at, row_key, column_name, ref_key, body, created_at FROM cell WHERE row_key = ? AND column_name = ? ORDER BY ref_key DESC LIMIT 1"
	putCellSQL          = "INSERT INTO cell ( row_key, column_name, ref_key, body ) VALUES(?, ?, ?, ?)"
	insertIndexSQL		= "INSERT INTO %s (row_key, %s) VALUES (?, ?) ON DUPLICATE KEY UPDATE %s = ?"
	queryIndexSQL		= "SELECT row_key FROM %s WHERE %s = ?"
)

func exec(db *sql.DB, sqlStr string) error {
	_, err := db.Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}

// New returns a new mysql-backed Storage
func New() *Storage {
	return new(Storage)
}

func (s *Storage) WithZap() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	sug := logger.Sugar()
	s.Sugar = sug
	return nil
}

func (s *Storage) Open() error {
	db, err := sql.Open(driver, fmt.Sprintf(dsnFormat, s.user, s.pass, s.host, s.port, s.database))
	if err != nil {
		return err
	}
	s.store = db

	s.indexes = make(map[string]*Index)
	return nil
}

func (s *Storage) WithUser(user string) *Storage {
	s.user = user
	return s
}

func (s *Storage) WithPass(pass string) *Storage {
	s.pass = pass
	return s
}

func (s *Storage) WithHost(host string) *Storage {
	s.host = host
	return s
}

func (s *Storage) WithPort(port string) *Storage {
	s.port = port
	return s
}

func (s *Storage) WithDatabase(database string) *Storage {
	s.database = database
	return s
}

func (s *Storage) GetCell(ctx context.Context, rowKey []byte, columnKey string, refKey int64) (cell models.Cell, found bool, err error) {
	var (
		resAddedAt   int64
		resRowKey    []byte
		resColName   string
		resRefKey    int64
		resBody      []byte
		resCreatedAt *time.Time
		rows         *sql.Rows
	)
	s.Sugar.Infow("GetCell", "query", getCellSQL, "rowKey", rowKey, "columnKey", columnKey, "refKey", refKey)
	rows, err = s.store.QueryContext(ctx, getCellSQL, rowKey, columnKey, refKey)
	if err != nil {
		return
	}
	defer rows.Close()

	found = false
	for rows.Next() {
		err = rows.Scan(&resAddedAt, &resRowKey, &resColName, &resRefKey, &resBody, &resCreatedAt)
		if err != nil {
			return
		}
		s.Sugar.Infow("GetCell scanned data", "AddedAt", resAddedAt, "RowKey", resRowKey, "ColName", resColName, "RefKey", resRefKey, "Body", resBody, "CreatedAt", resCreatedAt)

		cell.AddedAt = resAddedAt
		cell.RowKey = resRowKey
		cell.ColumnName = resColName
		cell.RefKey = resRefKey
		cell.Body = resBody
		cell.CreatedAt = resCreatedAt
		found = true
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return cell, found, nil
}

func (s *Storage) GetCellLatest(ctx context.Context, rowKey []byte, columnKey string) (cell models.Cell, found bool, err error) {
	var (
		resAddedAt   int64
		resRowKey    []byte
		resColName   string
		resRefKey    int64
		resBody      []byte
		resCreatedAt *time.Time
		rows         *sql.Rows
	)
	s.Sugar.Infow("GetCellLatest", "query before", getCellLatestSQL, "rowKey", rowKey, "columnKey", columnKey)
	rows, err = s.store.QueryContext(ctx, getCellLatestSQL, rowKey, columnKey)
	s.Sugar.Infow("GetCellLatest", "query after", getCellLatestSQL, "rowKey", rowKey, "columnKey", columnKey, "rows", rows, "error", err)
	if err != nil {
		return
	}
	defer rows.Close()

	found = false
	for rows.Next() {
		err = rows.Scan(&resAddedAt, &resRowKey, &resColName, &resRefKey, &resBody, &resCreatedAt)
		if err != nil {
			return
		}
		s.Sugar.Infow("GetCellLatest scanned data", "AddedAt", resAddedAt, "RowKey", resRowKey, "ColName", resColName, "RefKey", resRefKey, "Body", resBody, "CreatedAt", resCreatedAt)

		cell.AddedAt = resAddedAt
		cell.RowKey = resRowKey
		cell.ColumnName = resColName
		cell.RefKey = resRefKey
		cell.Body = resBody
		cell.CreatedAt = resCreatedAt
		found = true
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return cell, found, nil
}

func (s *Storage) GetCellsByFieldLatest(ctx context.Context, columnKey string, field string, value interface{}) (cells []models.Cell, found bool, err error) {
	// Add Index table if not exist
	table := s.checkAddIndex(columnKey, field)

	rowKeys := table.QueryByField(ctx, value)
	if len(rowKeys) == 0 {
		return cells, false, nil
	}

	for _, rowKey := range rowKeys {
		cell, _, err := s.GetCellLatest(ctx, rowKey, columnKey)
		utils.CheckErr(err)
		cells = append(cells, cell)
	}
	return cells, true, nil
}

func (s *Storage) CheckValueExist(ctx context.Context, columnKey string, field string, value interface{}) (found bool, err error) {
	table := s.getIndex(columnKey, field)
	if table == nil {
		return false, errors.New("invalid field")
	}
	return table.CheckValueExist(ctx, value), nil
}

func (s *Storage) putAllIndex(ctx context.Context, rowKey []byte, columnKey string, cell models.Cell, ignore_fields ...string) {
	var body map[string]interface{}
	err := json.Unmarshal(cell.Body, &body)
	utils.CheckErr(err)

	ignore_fields_ := make(map[string]bool)
	for _, field := range ignore_fields {
		ignore_fields_[field] = true
	}

	for field, value := range body {
		if _, ok := ignore_fields_[field]; !ok {
			table := s.checkAddIndex(columnKey, field)
			table.PutIndex(ctx, rowKey, value)
		}
	}
}

func (s *Storage) checkAddIndex(columnKey string, field string) *Index {
	tableName := utils.IndexTableName(columnKey, field)
	if _, ok := s.indexes[tableName]; !ok {
		s.indexes[tableName] = NewIndex(columnKey, field, s.store)
	}
	table, _ := s.indexes[tableName]
	return table
}

func (s *Storage) getIndex(columnKey string, field string) *Index {
	tableName := utils.IndexTableName(columnKey, field)
	table, _ := s.indexes[tableName]
	return table
}

func (s *Storage) PutCell(ctx context.Context, rowKey []byte, columnKey string, refKey int64, cell models.Cell, ignore_fileds ...string) (err error) {
	var stmt *sql.Stmt
	stmt, err = s.store.PrepareContext(ctx, putCellSQL)
	if err != nil {
		return
	}
	var res sql.Result
	s.Sugar.Infow("PutCell", "rowKey", rowKey, "columnKey", columnKey, "refKey", refKey, "Body", cell.Body)
	res, err = stmt.Exec(rowKey, columnKey, refKey, cell.Body)
	if err != nil {
		return
	}
	var lastID int64
	lastID, err = res.LastInsertId()
	if err != nil {
		return
	}
	var rowCnt int64
	rowCnt, err = res.RowsAffected()
	if err != nil {
		return
	}
	// TODO(rbastic): Should we side-affect the cell and record the AddedAt?
	s.Sugar.Infof("ID = %d, affected = %d\n", lastID, rowCnt)

	s.putAllIndex(ctx, rowKey, columnKey, cell, ignore_fileds...)
	return
}

// Destroy closes the in-memory store, and is a completely destructive operation.
func (s *Storage) Destroy(ctx context.Context) error {
	// TODO(rbastic): What do if there's an error in Sync()?
	s.Sugar.Sync()
	return s.store.Close()
}
