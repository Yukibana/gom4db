package cache

type KeyValueCache interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
	Del(key string) error
	Close()
	GetStat() KeyValueCacheStat
}
type KeyValueCacheStat struct {
}

type HashMapCache interface {
}
type ListCache interface {
}
type SortedSetCache interface {
}
type Cache interface {
	KeyValueCache
	HashMapCache
	ListCache
	SortedSetCache
}
