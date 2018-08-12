package schemaless

import (
	"context"
	"code.jogchat.internal/go-schemaless/utils"
	"github.com/satori/go.uuid"
	"encoding/json"
	"fmt"
	"code.jogchat.internal/go-schemaless/models"
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)


func newBusiness(id uuid.UUID, colKey string, domain string, name string) models.Cell {
	rowKey := utils.NewUUID().Bytes()
	refKey := time.Now().UnixNano()
	blob, err := json.Marshal(map[string]interface{} {
		"id": id,
		"domain": domain,
		"name": name,
	})
	utils.CheckErr(err)
	return models.Cell{RowKey: rowKey, ColumnName: colKey, RefKey: refKey, Body: blob}
}

// run this to make sure nothing breaks
func TestSchemaless(t *testing.T) {
	assert := assert.New(t)

	dataStore := InitDataStore()
	defer dataStore.Destroy(context.TODO())

	UIUC := newBusiness(utils.NewUUID(), "schools", "illinois.edu", "UIUC")
	err := dataStore.PutCell(context.TODO(), UIUC.RowKey, UIUC.ColumnName, UIUC.RefKey, UIUC)
	utils.CheckErr(err)

	CMU := newBusiness(utils.NewUUID(), "schools", "andrew.cmu.edu", "CMU")
	err = dataStore.PutCell(context.TODO(), CMU.RowKey, CMU.ColumnName, CMU.RefKey, CMU)

	Sift := newBusiness(utils.NewUUID(), "companies", "siftscience.com", "Sift Science")
	err = dataStore.PutCell(context.TODO(), Sift.RowKey, Sift.ColumnName, Sift.RefKey, Sift)
	utils.CheckErr(err)

	Yahoo := newBusiness(utils.NewUUID(), "companies", "yahoo-inc.com", "Yahoo!")
	err = dataStore.PutCell(context.TODO(), Yahoo.RowKey, Yahoo.ColumnName, Yahoo.RefKey, Yahoo)
	utils.CheckErr(err)

	var body map[string]interface{}

	cells, _, err := dataStore.GetCellsByFieldLatest(context.TODO(), "schools", "domain", "illinois.edu", "=")
	utils.CheckErr(err)
	assert.Equal(len(cells), 1)
	for _, cell := range cells {
		err := json.Unmarshal(cell.Body, &body)
		utils.CheckErr(err)
		assert.Equal(body["name"], "UIUC")
		assert.Equal(body["domain"], "illinois.edu")
		fmt.Println(cell.String())
	}

	cells, _, err = dataStore.GetCellsByColumnLatest(context.TODO(), "companies")
	utils.CheckErr(err)
	assert.Equal(len(cells), 2)
	for _, cell := range cells {
		assert.Equal(cell.ColumnName, "companies")
		fmt.Println(cell.String())
	}
}
