package doryMemory

import (
	"time"

	"github.com/peterbourgon/diskv"
)

/*
Diskv is abstracted here so it can be mocked.
*/
type Diskv struct {
	KV *diskv.Diskv
}

func (kvstore *Diskv) Add(key string, ttl time.Duration, dataBlob []byte) error {
	// ttl of no use here
	kvstore.KV.Write(key, dataBlob)
	return nil
}

func (kvstore *Diskv) Exists(key string) bool {
	return kvstore.KV.Has(key)
}

func (kvstore *Diskv) Delete(key string) error {
	err := kvstore.KV.Erase(key)
	return err
}

func (kvstore *Diskv) Value(key string) ([]byte, error) {
	return kvstore.KV.Read(key)
}

func (kvstore *Diskv) List() []string {
	keyIndex := 0
	keyCount := kvstore.Count()
	keyList := make([]string, keyCount)
	for key := range kvstore.KV.Keys(nil) {
		keyList[keyIndex] = key
		keyIndex += 1
	}
	return keyList
}

func (kvstore *Diskv) Count() int {
	var keyCount int
	for range kvstore.KV.Keys(nil) {
		keyCount += 1
	}
	return keyCount
}

func (kvstore *Diskv) Purge() error {
	var err error
	for key := range kvstore.KV.Keys(nil) {
		err = kvstore.KV.Erase(key)
		if err != nil {
			break
		}
	}
	return err
}
