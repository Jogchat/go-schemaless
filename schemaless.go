package schemaless

import (
	"code.jogchat.internal/go-schemaless/storage/mysql"
	"code.jogchat.internal/go-schemaless/utils"
	"code.jogchat.internal/go-schemaless/core"
	"os"
	"io/ioutil"
	"encoding/json"
)

func newBackend(user, pass, host, port, schemaName string) *mysql.Storage {
	m := mysql.New().WithUser(user).
		WithPass(pass).
		WithHost(host).
		WithPort(port).
		WithDatabase(schemaName)

	m.WithZap()
	m.Open()

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

func InitDataStore() *core.KVStore {
	jsonFile, err := os.Open("config/config.json")
	utils.CheckErr(err)
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	utils.CheckErr(err)

	var config map[string][]map[string]string
	json.Unmarshal(bytes, &config)

	shards := getShards(config)

	return core.New(shards)
}
