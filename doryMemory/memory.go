package doryMemory

/*
package for Create/Read/Delete actions over Cache2Go Table for items post aes encryption
*/

import (
	"log"
	"time"

	golcrypt "github.com/abhishekkr/gol/golcrypt"
	golrandom "github.com/abhishekkr/gol/golrandom"

	"github.com/muesli/cache2go"
)

/*
LocalAuth is a struct for Name as Auth-Path, Value with Gol Library Struct for Data/Cipher/Key nodes, TTLSecond for expiry timeout.
*/
type LocalAuth struct {
	Name      string
	Value     golcrypt.AESBlock
	TTLSecond uint64
}

/*
NewLocalAuthStore instantiates and return a Cache2Go Table store.
*/
func NewLocalAuthStore(cacheName string) *cache2go.CacheTable {
	return cache2go.Cache(cacheName)
}

/*
Exists checks if a Auth-Path exists in Cache2Go Table.
*/
func (auth *LocalAuth) Exists(localAuthStore *cache2go.CacheTable) bool {
	return localAuthStore.Exists(auth.Name)
}

/*
Set stores an encrypted value with random/provided Token at Auth-Path in Cache2Go Table.
*/
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

	return auth.Exists(localAuthStore)
}

/*
Get fetchs a value decrypted by Token stored at a Auth-Path in Cache2Go Table.
*/
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

/*
Delete purges a Auth-Path in Cache2Go Table, if it's value is decipherable by given Token.
*/
func (auth *LocalAuth) Delete(localAuthStore *cache2go.CacheTable) bool {
	var err error

	if localAuthStore == nil {
		log.Println("delete triggered for missing auth-store")
		return false
	}
	if auth.Value.Key == nil {
		log.Println("delete triggered for empty key")
		return false
	}
	if !auth.Exists(localAuthStore) {
		log.Println("delete triggered for missing auth identifier")
		return false
	}

	cipherCacheItem, err := localAuthStore.Value(auth.Name)
	cipherData := cipherCacheItem.Data().(*golcrypt.Cipher)
	auth.Value.Cipher = *cipherData

	if err = auth.Value.Decrypt(); err != nil {
		log.Println("to delete decrypt shall pass;", err)
		return false
	}
	auth.Value.Cipher = nil
	auth.Value.DataBlob = nil

	_, err = localAuthStore.Delete(auth.Name)
	if err != nil {
		log.Println("delete triggered but", err)
		return false
	}

	return true
}
