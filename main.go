package main

import (
	"context"

	"code.jogchat.internal/go-schemaless/utils"
	"code.jogchat.internal/go-schemaless/core"
	"code.jogchat.internal/go-schemaless/storage/mysql"
	"github.com/satori/go.uuid"
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"code.jogchat.internal/go-schemaless/models"
	"time"
)

func newBackend(user, pass, host, port, schemaName string) *mysql.Storage {
	m := mysql.New().WithUser(user).
		WithPass(pass).
		WithHost(host).
		WithPort(port).
		WithDatabase(schemaName)

	err := m.WithZap()
	utils.CheckErr(err)
	err = m.Open()
	utils.CheckErr(err)

	// TODO(rbastic): defer Sync() on all backend storage loggers
	return m
}

func getShards(config map[string][]map[string]string) []core.Shard {
	var shards []core.Shard
	hosts := config["hosts"]

	for _, host := range hosts {
		shard := core.Shard{
			Name: host["database"],
			Backend: newBackend(host["user"], host["password"], host["ip"], host["port"], host["database"])}
		shards = append(shards, shard)
	}

	return shards
}

func newUUID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func newBusiness(id uuid.UUID, colKey string, category string, domain string, name string) models.Cell {
	rowKey := newUUID().Bytes()
	refKey := time.Now().UnixNano()
	blob, err := json.Marshal(map[string]interface{} {
		"id": id,
		"category": category,
		"domain": domain,
		"name": name,
	})
	utils.CheckErr(err)
	return models.Cell{RowKey: rowKey, ColumnName: colKey, RefKey: refKey, Body: blob}
}

func main() {
	jsonFile, err := os.Open("config/config.json")
	utils.CheckErr(err)
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	utils.CheckErr(err)

	var config map[string][]map[string]string
	json.Unmarshal(bytes, &config)

	shards := getShards(config)

	dataStore := core.New(shards)
	defer dataStore.Destroy(context.TODO())

	UIUC := newBusiness(newUUID(), "schools", "university", "illinois.edu", "UIUC")
	err = dataStore.PutCell(context.TODO(), UIUC.RowKey, UIUC.ColumnName, UIUC.RefKey, UIUC)
	utils.CheckErr(err)

	CMU := newBusiness(newUUID(), "schools", "university", "andrew.cmu.edu", "CMU")
	err = dataStore.PutCell(context.TODO(), CMU.RowKey, CMU.ColumnName, CMU.RefKey, CMU)
	utils.CheckErr(err)

	Yahoo := newBusiness(newUUID(), "companies", "technology", "yahoo-inc.com", "Yahoo!")
	err = dataStore.PutCell(context.TODO(), Yahoo.RowKey, Yahoo.ColumnName, Yahoo.RefKey, Yahoo)
	utils.CheckErr(err)

	cells, _, err := dataStore.GetCellsByFieldLatest(context.TODO(), "schools", "category", "university")
	utils.CheckErr(err)
	for _, cell := range cells {
		fmt.Println(cell.String())
	}
}
