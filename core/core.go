package core

import (
	jh "code.jogchat.internal/dgryski-go-shardedkv/choosers/jump"
	"code.jogchat.internal/go-schemaless/storage/mysql"
	"context"
	"code.jogchat.internal/go-schemaless/models"
	"sync"
	"code.jogchat.internal/dgryski-go-metro"
	"code.jogchat.internal/go-schemaless/utils"
)

// KVStore is a sharded key-value store
type KVStore struct {
	continuum Chooser
	storages  map[string]*mysql.Storage

	migration Chooser
	mstorages map[string]*mysql.Storage

	// we avoid holding the lock during a call to a storage engine, which may block
	mu	sync.Mutex
}

// Chooser maps keys to shards
type Chooser interface {
	// SetBuckets sets the list of known buckets from which the chooser should select
	SetBuckets([]string) error
	// Choose returns a bucket for a given key
	Choose(key string) string
	// Buckets returns the list of known buckets
	Buckets() []string
}

// Shard is a named storage backend
type Shard struct {
	Name    string
	Backend *mysql.Storage
}

func hash64(b []byte) uint64 { return metro.Hash64(b, 0) }


// New returns a KVStore that uses chooser to shard the keys across the provided shards
func New(shards []Shard) *KVStore {
	chooser := jh.New(hash64)

	var buckets []string
	kv := &KVStore{
		continuum: chooser,
		storages:  make(map[string]*mysql.Storage),
		// what about migration?
	}
	for _, shard := range shards {
		buckets = append(buckets, shard.Name)
		kv.AddShard(shard.Name, shard.Backend)
	}
	chooser.SetBuckets(buckets)
	return kv
}

func (kv *KVStore) GetCell(ctx context.Context, rowKey []byte, columnKey string, refKey int64) (cell models.Cell, found bool, err error) {
	var storage *mysql.Storage
	var migStorage *mysql.Storage

	kv.mu.Lock()
	defer kv.mu.Unlock()

	if kv.migration != nil {
		shard := kv.migration.Choose(string(rowKey))
		migStorage = kv.mstorages[shard]
	}
	shard := kv.continuum.Choose(string(rowKey))
	storage = kv.storages[shard]

	if migStorage != nil {
		val, ok, err := (*migStorage).GetCell(ctx, rowKey, columnKey, refKey)
		if ok {
			return val, ok, err
		}
	}

	return (*storage).GetCell(ctx, rowKey, columnKey, refKey)
}

func (kv *KVStore) GetCellLatest(ctx context.Context, rowKey []byte, columnKey string) (cell models.Cell, found bool, err error) {
	var storage *mysql.Storage
	var migStorage *mysql.Storage

	kv.mu.Lock()
	defer kv.mu.Unlock()

	if kv.migration != nil {
		shard := kv.migration.Choose(string(rowKey))
		migStorage = kv.mstorages[shard]
	}

	if migStorage != nil {
		val, ok, err := (*migStorage).GetCellLatest(ctx, rowKey, columnKey)
		if err != nil {
			return val, ok, err
		}
		if ok {
			return val, ok, nil
		}
	}

	shard := kv.continuum.Choose(string(rowKey))
	storage = kv.storages[shard]

	return (*storage).GetCellLatest(ctx, rowKey, columnKey)
}

func (kv *KVStore) GetCellsByFieldLatest(ctx context.Context, columnKey string, field string, value interface{}) (cells []models.Cell, found bool, err error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	for _, storage := range kv.storages {
		cells_, found, err := (*storage).GetCellsByFieldLatest(ctx, columnKey, field, value)
		if !found {
			continue
		}
		utils.CheckErr(err)
		cells = append(cells, cells_...)
	}

	found = true
	if len(cells) == 0 {
		found = false
	}
	return cells, found, nil
}

// PutCell
func (kv *KVStore) PutCell(ctx context.Context, rowKey []byte, columnKey string, refKey int64, cell models.Cell) error {
	var storage *mysql.Storage

	kv.mu.Lock()
	defer kv.mu.Unlock()

	if kv.migration != nil {
		shard := kv.migration.Choose(string(rowKey))
		storage = kv.mstorages[shard]

		return (*storage).PutCell(ctx, rowKey, columnKey, refKey, cell)
	}

	shard := kv.continuum.Choose(string(rowKey))
	storage = kv.storages[shard]

	return (*storage).PutCell(ctx, rowKey, columnKey, refKey, cell)
}

func (kv *KVStore) PartitionRead(ctx context.Context, partitionNumber int, location string, value interface{}, limit int) (cells []models.Cell, found bool, err error) {

	kv.mu.Lock()
	defer kv.mu.Unlock()

	if kv.migration != nil {
		buckets := kv.migration.Buckets()
		shard := buckets[partitionNumber]
		migStorage := kv.mstorages[shard]

		if migStorage != nil {
			return (*migStorage).PartitionRead(ctx, partitionNumber, location, value, limit)
		}
	}

	buckets := kv.continuum.Buckets()
	shard := buckets[partitionNumber]
	storage := kv.storages[shard]

	return (*storage).PartitionRead(ctx, partitionNumber, location, value, limit)
}

// ResetConnection implements Storage.ResetConnection()
func (kv *KVStore) ResetConnection(ctx context.Context, key string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if kv.migration != nil {
		shard := kv.migration.Choose(key)
		migStorage := kv.mstorages[shard]

		if migStorage != nil {
			err := (*migStorage).ResetConnection(ctx, key)
			if err != nil {
				return err
			}
		}
	}
	shard := kv.continuum.Choose(key)
	storage := kv.storages[shard]

	return (*storage).ResetConnection(ctx, key)
}

// Destroy implements Storage.Destroy()
func (kv *KVStore) Destroy(ctx context.Context) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if kv.migration != nil {
		for _, migStorage := range kv.mstorages {
			err := (*migStorage).Destroy(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for _, store := range kv.storages {
		err := (*store).Destroy(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddShard adds a shard from the list of known shards
func (kv *KVStore) AddShard(shard string, storage *mysql.Storage) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.storages[shard] = storage
}

// DeleteShard removes a shard from the list of known shards
func (kv *KVStore) DeleteShard(shard string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	delete(kv.storages, shard)
}

// BeginMigration begins a continuum migration.  All the shards in the new
// continuum must already be known to the KVStore via AddShard().
func (kv *KVStore) BeginMigration(continuum Chooser) {

	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.migration = continuum
	kv.mstorages = kv.storages
}

// BeginMigrationWithShards begins a continuum migration using the new set of shards.
func (kv *KVStore) BeginMigrationWithShards(continuum Chooser, shards []Shard) {

	kv.mu.Lock()
	defer kv.mu.Unlock()

	var buckets []string
	mstorages := make(map[string]*mysql.Storage)
	for _, shard := range shards {
		buckets = append(buckets, shard.Name)
		mstorages[shard.Name] = shard.Backend
	}

	continuum.SetBuckets(buckets)

	kv.migration = continuum
	kv.mstorages = mstorages
}

// EndMigration ends a continuum migration and marks the migration continuum
// as the new primary
func (kv *KVStore) EndMigration() {

	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.continuum = kv.migration
	kv.migration = nil

	kv.storages = kv.mstorages
	kv.mstorages = nil
}
