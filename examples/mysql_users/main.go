package main

import (
	"context"

	"code.jogchat.internal/go-schemaless"
	"code.jogchat.internal/go-schemaless/utils"
	"code.jogchat.internal/go-schemaless/core"
	st "code.jogchat.internal/go-schemaless/storage/mysql"
	"github.com/satori/go.uuid"
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
	"code.jogchat.internal/go-schemaless/models"
)

func newBackend(user, pass, host, port, schemaName string) *st.Storage {
	m := st.New().WithUser(user).
		WithPass(pass).
		WithHost(host).
		WithPort(port).
		WithDatabase(schemaName)

	err := m.WithZap()
	if err != nil {
		panic(err)
	}

	err = m.Open()
	if err != nil {
		panic(err)
	}

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

func main() {
	jsonFile, err := os.Open("config/config.json")
	utils.CheckErr(err)
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	utils.CheckErr(err)

	var config map[string][]map[string]string
	json.Unmarshal(bytes, &config)

	shards := getShards(config)

	//kv := schemaless.New().WithSource(shards)
	//defer kv.Destroy(context.TODO())

	// We're going to demonstrate jump hash+metro hash with MySQL-backed
	// storage. This example implements multiple shard schemas on a single
	// node.

	// You decide the refKey's purpose. For example, it can
	// be used as a record version number, or for sort-order.

	//for i := 0; i < 1000; i++ {
	//	refKey := int64(i)
	//	kv.PutCell(context.TODO(), newUUID(), "PII", refKey, models.Cell{RefKey: refKey, Body: fakeUserJSON()})
	//}

	dataStore := schemaless.New().WithSource(shards)
	defer dataStore.Destroy(context.TODO())

	rowKey := newUUID().Bytes()
	colKey := "school"
	refKey := time.Now().UnixNano()
	blob, err := json.Marshal(map[string]string {
		"id": newUUID().String(),
		"category": "university",
		"domain": "illinois.edu",
		"name": "UIUC",
	})
	utils.CheckErr(err)

	err = dataStore.PutCell(context.TODO(), rowKey, colKey, refKey,
		models.Cell{RefKey: refKey, Body: blob})
	utils.CheckErr(err)

	cell, _, err := dataStore.GetCellLatest(context.TODO(), rowKey, colKey)
	utils.CheckErr(err)
	fmt.Println(cell.String())
}
