package mysql

import (
	"database/sql"
	"fmt"
	"context"
	"code.jogchat.internal/go-schemaless/utils"
)

// PutIndex updates all Index tables relevant to the current cell, If entry does not exist, insert into Index table instead
func PutIndex(ctx context.Context, conn *sql.DB, column string, field string, rowKey []byte, value interface{}) {
	stmt, err := conn.PrepareContext(ctx, fmt.Sprintf(insertIndexSQL, utils.IndexTableName(column, field), field, field))
	utils.CheckErr(err)
	_, err = stmt.Exec(rowKey, value, value)
	utils.CheckErr(err)
}

func QueryByField(ctx context.Context, conn *sql.DB, column string, field string, value interface{}, operator string) [][]byte {
	stmt := fmt.Sprintf(queryIndexSQL, utils.IndexTableName(column, field), field, operator)
	rows, err := conn.QueryContext(ctx, stmt, value)
	utils.CheckErr(err)
	return extractRowKeys(rows)
}

func QueryAll(ctx context.Context, conn *sql.DB, column string, field string) [][]byte {
	stmt := fmt.Sprintf(queryIndexAllSQL, utils.IndexTableName(column, field))
	rows, err := conn.QueryContext(ctx, stmt)
	utils.CheckErr(err)
	return extractRowKeys(rows)
}

// Check if value exist in index table, return true if value already exist
func CheckValueExist(ctx context.Context, conn *sql.DB, column string, field string, value interface{}) bool {
	stmt := fmt.Sprintf(queryIndexSQL, utils.IndexTableName(column, field), field, "=")
	results, err := conn.QueryContext(ctx, stmt, value)
	utils.CheckErr(err)
	return results.Next()
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
