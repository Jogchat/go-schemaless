package main

import (
	"context"

	"github.com/dgryski/go-metro"
	"github.com/icrowley/fake"
	"code.jogchat.internal/go-schemaless"
	"code.jogchat.internal/go-schemaless/core"
	"code.jogchat.internal/go-schemaless/models"
	st "code.jogchat.internal/go-schemaless/storage/mysql"
	"github.com/satori/go.uuid"
	"os"
	//"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
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

func hash64(b []byte) uint64 { return metro.Hash64(b, 0) }

func newUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func fakeUserJSON() string {
	name := fake.FirstName() + " " + fake.LastName()
	return "{\"name" + "\": \"" + name + "\"}"
}

func main() {
	jsonFile, err := os.Open("config/config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var config map[string][]map[string]string
	json.Unmarshal(bytes, &config)

	shards := getShards(config)
	kv := schemaless.New().WithSource(shards)
	defer kv.Destroy(context.TODO())

	// We're going to demonstrate jump hash+metro hash with MySQL-backed
	// storage. This example implements multiple shard schemas on a single
	// node.

	// You decide the refKey's purpose. For example, it can
	// be used as a record version number, or for sort-order.
	for i := 0; i < 1000; i++ {
		refKey := int64(i)
		kv.PutCell(context.TODO(), newUUID(), "PII", refKey, models.Cell{RefKey: refKey, Body: fakeUserJSON()})
	}
}
