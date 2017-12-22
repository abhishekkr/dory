package doryMemory

import (
	"fmt"
	"time"

	golcrypt "github.com/abhishekkr/gol/golcrypt"

	"github.com/muesli/cache2go"
)

/*
Cache2Go is abstracted here so it can be mocked.
*/
type Cache2Go struct {
	CacheTable *cache2go.CacheTable
}

func (cache *Cache2Go) Add(key string, ttl time.Duration, dataBlob []byte) error {
	cache.CacheTable.Add(key, ttl, dataBlob)
	return nil
}

func (cache *Cache2Go) Exists(key string) bool {
	return cache.CacheTable.Exists(key)
}

func (cache *Cache2Go) Delete(key string) error {
	_, err := cache.CacheTable.Delete(key)
	return err
}

func (cache *Cache2Go) Value(key string) ([]byte, error) {
	var cipherData golcrypt.Cipher
	cipherCacheItem, err := cache.CacheTable.Value(key)

	cipherData = cipherCacheItem.Data().([]byte)
	return cipherData, err
}

func (cache *Cache2Go) List() []string {
	keyIndex := 0
	keyList := make([]string, cache.CacheTable.Count())
	cache.CacheTable.Foreach(func(key interface{}, item *cache2go.CacheItem) {
		keyList[keyIndex] = fmt.Sprintf("%q", key)
		keyIndex += 1
	})
	return keyList
}

func (cache *Cache2Go) Count() int {
	return cache.CacheTable.Count()
}

func (cache *Cache2Go) Purge() error {
	var err error
	cache.CacheTable.Flush()
	return err
}
