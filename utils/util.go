package utils

import "github.com/satori/go.uuid"

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NewUUID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

// Convert column key (equivalent to SQL table name) and field into index table name
func IndexTableName(columnKey string, field string) string {
	return "index_" + columnKey + "_" + field
}
