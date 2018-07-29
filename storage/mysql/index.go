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
	i := new(Index)
	i.Column = col
	i.Field = field
	i.conn = conn
	return i
}

// PutIndex updates all Index tables relevant to the current cell, If entry does not exist, insert into Index table instead
func (i *Index) PutIndex(ctx context.Context, rowKey []byte, value interface{}) {
	stmt, err := i.conn.PrepareContext(ctx, fmt.Sprintf(insertIndexSQL, utils.IndexTableName(i.Column, i.Field), i.Field, i.Field))
	utils.CheckErr(err)
	_, err = stmt.Exec(rowKey, value, value)
	utils.CheckErr(err)
}

func (i *Index) QueryByField(ctx context.Context, value interface{}, operator string) [][]byte {
	stmt := fmt.Sprintf(queryIndexSQL, utils.IndexTableName(i.Column, i.Field), i.Field, operator)
	rows, err := i.conn.QueryContext(ctx, stmt, value)
	utils.CheckErr(err)
	return extractRowKeys(rows)
}

func (i *Index) QueryAll(ctx context.Context) [][]byte {
	stmt := fmt.Sprintf(queryIndexAllSQL, utils.IndexTableName(i.Column, i.Field))
	rows, err := i.conn.QueryContext(ctx, stmt)
	utils.CheckErr(err)
	return extractRowKeys(rows)
}

func extractRowKeys(rows *sql.Rows) [][]byte {
	var rowKeys [][]byte
	for rows.Next() {
		var rowKey []byte
		err := rows.Scan(&rowKey)
		utils.CheckErr(err)
		rowKeys = append(rowKeys, rowKey)
	}
	return rowKeys
}

// Check if value exist in index table, return true if value already exist
func (i *Index) CheckValueExist(ctx context.Context, value interface{}) bool {
	stmt := fmt.Sprintf(queryIndexSQL, utils.IndexTableName(i.Column, i.Field), i.Field, "=")
	results, err := i.conn.QueryContext(ctx, stmt, value)
	utils.CheckErr(err)
	return results.Next()
}
