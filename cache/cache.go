package cache

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
}

func NewCache()Cache{
	rdb := NewRocksDB("./rocks_data")
	db :=  NewDb(rdb,100)
	return db
}
