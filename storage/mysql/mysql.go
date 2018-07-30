// Package mysql is a mysql-backed Schemaless store.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"code.jogchat.internal/go-schemaless/models"
	"go.uber.org/zap"
	"time"
	"encoding/json"
	"code.jogchat.internal/go-schemaless/utils"
	"github.com/pkg/errors"
)

// Storage is a MySQL-backed storage.
type Storage struct {
	user     string
	pass     string
	host     string
	port     string
	database string

	store	*sql.DB
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

	getCellLatestSQL    = "SELECT added_at, row_key, column_name, ref_key, body, created_at FROM cell WHERE row_key = ? AND column_name = ? ORDER BY ref_key DESC LIMIT 1"
	getCellsLatestSQL	= "SELECT added_at, row_key, column_name, ref_key, body, created_at FROM (SELECT * FROM %s ORDER BY %s, %s ASC) GROUP BY %s"
	joinTableSQL		= "%s INNER JOIN %s ON %s.row_key = %s.row_key"
	joinTableFilterSQL	= "%s INNER JOIN %s ON %s.row_key = %s.row_key WHERE %s = ?"
	putCellSQL          = "INSERT INTO cell (row_key, column_name, ref_key, body) VALUES(?, ?, ?, ?)"
	insertIndexSQL		= "INSERT INTO %s (row_key, %s) VALUES (?, ?) ON DUPLICATE KEY UPDATE %s = ?"
	queryIndexSQL		= "SELECT row_key FROM %s WHERE %s %s ?"
	queryIndexAllSQL	= "SELECT row_key FROM %s"
)

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

func (s *Storage) GetCellsByColumnLatest(ctx context.Context, columnKey string) (cells []models.Cell, found bool, err error) {
	// Add Index table if not exist
	rowKeys := QueryAll(ctx, s.store, columnKey, "id")
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

func (s *Storage) GetCellByUniqueFieldLatest(ctx context.Context, columnKey string, field string, value interface{}) (cell models.Cell, found bool, err error) {
	// Add Index table if not exist
	rowKeys := QueryByField(ctx, s.store, columnKey, field, value, "=")
	if len(rowKeys) == 0 {
		return cell, false, nil
	}
	if len(rowKeys) > 1 {
		panic(errors.New("field value not unique"))
	}

	return s.GetCellLatest(ctx, rowKeys[0], columnKey)
}

func (s *Storage) GetCellsByFieldLatest(ctx context.Context, columnKey string, field string, value interface{}, operator string) (cells []models.Cell, found bool, err error) {
	var (
		resAddedAt   int64
		resRowKey    []byte
		resColName   string
		resRefKey    int64
		resBody      []byte
		resCreatedAt *time.Time
		cell models.Cell
		rows         *sql.Rows
	)
	indexTable := utils.IndexTableName(columnKey, field)
	joinStmt := fmt.Sprintf(joinTableFilterSQL, "cell", indexTable, "cell", indexTable, field)
	queryStmt := fmt.Sprintf(getCellsLatestSQL, joinStmt, "row_key", "ref_key", "row_key")
	rows, err = s.store.QueryContext(ctx, queryStmt, value)
	if err != nil {
		return nil, false, err
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
		cells = append(cells, cell)
		found = true
	}
	return cells, found, nil
}

func (s *Storage) CheckValueExist(ctx context.Context, columnKey string, field string, value interface{}) (found bool, err error) {
	return CheckValueExist(ctx, s.store, columnKey, field, value), nil
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
			PutIndex(ctx, s.store, columnKey, field, rowKey, value)
		}
	}
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
