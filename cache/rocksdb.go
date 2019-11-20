package cache

import (
	"sync"

	"github.com/golang/groupcache/lru"
	"github.com/tecbot/gorocksdb"
)

type DB struct {
	rdb    *gorocksdb.DB
	wo     *gorocksdb.WriteOptions
	ro     *gorocksdb.ReadOptions
	mu     sync.Mutex
	caches *lru.Cache
}

func NewRocksDB(dir string) *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	rdb, err := gorocksdb.OpenDb(opts, dir)
	if err != nil {
		panic(err)
	}
	return rdb
}
func NewDb(rdb *gorocksdb.DB) *DB {
	db := &DB{rdb: rdb}
	db.wo = gorocksdb.NewDefaultWriteOptions()
	db.ro = gorocksdb.NewDefaultReadOptions()
	db.caches = lru.New(1000)
	return db
}
func (d *DB) WriteBatch(batch *gorocksdb.WriteBatch) error {
	return d.rdb.Write(d.wo, batch)
}

func (d *DB) RawGet(key []byte) ([]byte, error) {
	return d.rdb.GetBytes(d.ro, key)
}

func (d *DB) RawSet(key, value []byte) error {
	return d.rdb.Put(d.wo, key, value)
}

func (d *DB) RawDelete(key []byte) error {
	return d.rdb.Delete(d.wo, key)
}

func (d *DB) Close() {
	d.wo.Destroy()
	d.ro.Destroy()
	d.rdb.Close()
}
