package doryMemory

import (
	"log"
	"time"

	golcrypt "github.com/abhishekkr/gol/golcrypt"
	golrandom "github.com/abhishekkr/gol/golrandom"

	"github.com/muesli/cache2go"
)

type LocalAuth struct {
	Name      string
	Value     golcrypt.AESBlock
	TTLSecond uint64
}

func NewLocalAuthStore(cacheName string) *cache2go.CacheTable {
	return cache2go.Cache(cacheName)
}

func (auth *LocalAuth) Set(localAuthStore *cache2go.CacheTable) bool {
	if localAuthStore == nil {
		return false
	}
	if auth.Value.Key == nil {
		auth.Value.Key = []byte(golrandom.Token(32)) //size 16/24/32 allowed
	}

	if err := auth.Value.Encrypt(); err != nil {
		log.Println(err)
		return false
	}
	auth.Value.DataBlob = nil

	if auth.TTLSecond == 0 {
		auth.TTLSecond = 300 //default 5minute
	}

	ttl := time.Duration(auth.TTLSecond) * time.Second
	localAuthStore.Add(auth.Name, ttl, &auth.Value.Cipher)

	return localAuthStore.Exists(auth.Name)
}

func (auth *LocalAuth) Get(localAuthStore *cache2go.CacheTable) bool {
	if localAuthStore == nil {
		return false
	}
	if auth.Value.Key == nil {
		return false
	}

	cipherCacheItem, err := localAuthStore.Value(auth.Name)
	cipherData := cipherCacheItem.Data().(*golcrypt.Cipher)
	auth.Value.Cipher = *cipherData

	if err != nil {
		return false
	}

	if err = auth.Value.Decrypt(); err != nil {
		log.Println("failed to decrypt;", err)
		return false
	}

	return true
}

func (auth *LocalAuth) Delete(localAuthStore *cache2go.CacheTable) bool {
	var err error

	if localAuthStore == nil {
		return true
	}

	if auth.Value.Key == nil || !localAuthStore.Exists(auth.Name) {
		return true
	}

	_, err = localAuthStore.Delete(auth.Name)
	if err != nil {
		return false
	}

	return true
}
