package mysql

import (
	"database/sql"
	"fmt"
	"context"
	"code.jogchat.internal/go-schemaless/utils"
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
func (i *Index) PutIndex(ctx context.Context, rowKey []byte, value interface{}) {
	res := i.execCtx(ctx, updateIndexSQL, value, rowKey)
	rowCnt, err := res.RowsAffected()
	utils.CheckErr(err)

	if rowCnt == 0 {
		i.execCtx(ctx, insertIndexSQL, rowKey, value)
	}
}

func (i *Index) QueryByField(ctx context.Context, value interface{}) ([][]byte) {
	rows, err := i.conn.QueryContext(ctx, queryIndexSQL, value)
	utils.CheckErr(err)
	var rowKeys [][]byte

	for rows.Next() {
		var field interface{}
		var rowKey []byte
		err = rows.Scan(&field, &rowKey)
		utils.CheckErr(err)
		rowKeys = append(rowKeys, rowKey)
	}
	return rowKeys
}

func (i *Index) execCtx(ctx context.Context, rawStmt string, args ...interface{}) sql.Result {
	stmt, err := i.conn.PrepareContext(ctx, fmt.Sprintf(rawStmt, indexTableName(i.Column, i.Field), i.Field))
	utils.CheckErr(err)
	res, err := stmt.Exec(args)
	utils.CheckErr(err)
	return res
}

func indexTableName(columnKey string, field string) string {
	return "index_" + columnKey + "_" + field
}
