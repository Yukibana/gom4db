package cache

func (d *DB) Get(key string) ([]byte, error) {
	keyBytes := Str2bytes(key)
	return d.RawGet(rawKey(keyBytes, STRING))
}
func (d *DB) Set(key string, value []byte) error {
	keyBytes := Str2bytes(key)
	return d.RawSet(rawKey(keyBytes, STRING), value)
}
func (d *DB) Del(key string) error {
	keyBytes := Str2bytes(key)
	return d.RawDelete(keyBytes)
}
func (d DB) GetStat() KeyValueCacheStat {
	return KeyValueCacheStat{}
}
func NewKeyValueCache() KeyValueCache {
	rocksDb := NewRocksDB("./rocks_data")
	db := NewDb(rocksDb)
	return db
}
