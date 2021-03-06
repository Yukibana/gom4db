package cache

import (
	"github.com/tecbot/gorocksdb"
	"time"
)

// We can call it a extension of db
type kvcache struct {
	*DB
	batchModChan chan batchConfig
	setCommands  chan setCommand
}
type batchConfig struct {
	delay       time.Duration
	maxBatchNum int
	stop        bool
}
type setCommand struct {
	key, value []byte
}

// user can choose to use async write or not
//func (c *kvcache) Set(key string, value []byte) error {
//	fmt.Println("Already set",key,"with",Bytes2str(value))
//	keyBytes := Str2bytes(key)
//	return c.RawSet(keyBytes, value)
//}
//func (c *kvcache)Get(){
//
//}
func (c *kvcache) AsyncSet(key string, value []byte) {
	keyBytes := Str2bytes(key)
	c.setCommands <- setCommand{
		key:   keyBytes,
		value: value,
	}
}
func (c *kvcache) ModWriteConfig(newDelay time.Duration, newMaxBatchNum int) {
	newConfig := batchConfig{
		delay:       newDelay,
		maxBatchNum: newMaxBatchNum,
	}
	c.batchModChan <- newConfig
	return
}
func (c *kvcache) Close() {
	stopConfig := batchConfig{
		stop: true,
	}
	c.batchModChan <- stopConfig
	c.DB.Close()
}

func (c *kvcache) BatchDaemon(initConfig batchConfig) {
	batch := gorocksdb.NewWriteBatch()
	defer batch.Destroy()

	localConfig := initConfig
	for {
		select {
		case <-time.After(localConfig.delay):
			err := c.WriteBatch(batch)
			if err != nil {
				// TODO How to solve error
				return
			}
			batch.Clear()
		case newConfig := <-c.batchModChan:
			err := c.WriteBatch(batch)
			if err != nil {
				return
			}
			batch.Clear()
			localConfig = newConfig
		case newSetCommand := <-c.setCommands:
			batch.Put(newSetCommand.key, newSetCommand.value)
		}
	}
}

func NewKeyValueCache() KeyValueCache {
	rocksDb := NewRocksDB("./rocks_data")
	db := NewDb(rocksDb)
	cache := &kvcache{
		DB:           db,
		batchModChan: make(chan batchConfig, 1),
		setCommands:  make(chan setCommand, 1000),
	}

	initBatchConfig := batchConfig{
		delay:       time.Second * 2,
		maxBatchNum: 1000,
	}
	go cache.BatchDaemon(initBatchConfig)
	return cache
}
