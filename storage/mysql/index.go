package mysql

import (
	"database/sql"
	"fmt"
	"context"
	"log"
)

type Index struct {
	Column	string
	Field	string
	conn	*sql.DB
}

func NewIndex(col string, field string, conn *sql.DB) *Index {
	i := &Index{Column: col, Field: field, conn: conn}
	return i
}

// PutIndex updates all Index tables relevant to the currnet cell, If entry does not exist, insert into Index table instead
func (i *Index) PutIndex(ctx context.Context, rowKey []byte, columnKey string, field string, value interface{}) {
	res := execContext(ctx, i.conn, updateIndexSQL, rowKey, columnKey, field, value)
	rowCnt, err := res.RowsAffected()
	log.Fatal(err)

	if rowCnt == 0 {
		execContext(ctx, i.conn, insertIndexSQL, rowKey, columnKey, field, value)
	}
}

func execContext(ctx context.Context, conn *sql.DB, rawStmt string, rowKey []byte, columnKey string, field string, value interface{}) (sql.Result) {
	stmt, err := conn.PrepareContext(ctx, fmt.Sprintf(rawStmt, indexTableName(columnKey, field), field))
	log.Fatal(err)
	res, err := stmt.Exec(rowKey, value)
	log.Fatal(err)
	return res
}

func indexTableName(columnKey string, field string) string {
	return "index_" + columnKey + "_" + field
}
