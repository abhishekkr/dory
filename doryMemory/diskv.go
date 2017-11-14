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
