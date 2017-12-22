package doryMemory

import (
	"fmt"
	"path"
	"strconv"
	"time"

	golcrypt "github.com/abhishekkr/gol/golcrypt"
	golenv "github.com/abhishekkr/gol/golenv"
	gollog "github.com/abhishekkr/gol/gollog"
	golrandom "github.com/abhishekkr/gol/golrandom"

	"github.com/muesli/cache2go"
	"github.com/peterbourgon/diskv"
)

var (
	DoryDiskvBaseDir = golenv.OverrideIfEnv("DORY_DISKV_BASE_DIR", "/tmp")
	DoryDiskvCacheMB = golenv.OverrideIfEnv("DORY_DISKV_CACHE_MB", "128")
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
func NewLocalAuthStore(cacheName string) DataStore {
	return &Cache2Go{CacheTable: cache2go.Cache(cacheName)}
}

/*
NewDiskv instantiates and return a GolevelDB DBEngine.
*/
func NewDiskv(cacheName string) DataStore {
	diskvCacheMB, err := strconv.Atoi(DoryDiskvCacheMB)
	if err != nil {
		diskvCacheMB = 128
	}
	kv := diskv.New(diskv.Options{
		BasePath:     path.Join(DoryDiskvBaseDir, cacheName),
		CacheSizeMax: 1024 * 1024 * uint64(diskvCacheMB), // 128MB
	})
	return &Diskv{KV: kv}
}

/*
Exists checks if a Auth-Path exists in Cache2Go Table.
*/
func (auth *LocalAuth) Exists(dataStore DataStore) bool {
	gollog.Debug(fmt.Sprintf("key '%s' exists: %t", auth.Name, dataStore.Exists(auth.Name)))
	return dataStore.Exists(auth.Name)
}

/*
Set stores an encrypted value with random/provided Token at Auth-Path in Cache2Go Table.
*/
func (auth *LocalAuth) Set(dataStore DataStore) bool {
	if dataStore == nil {
		gollog.Err(fmt.Sprintf("key '%s' sent to corrupted datastore", auth.Name))
		return false
	}
	if auth.Value.Key == nil {
		auth.Value.Key = []byte(golrandom.Token(32)) //size 16/24/32 allowed
	}

	if err := auth.Value.Encrypt(); err != nil {
		gollog.Err(err.Error())
		return false
	}
	auth.Value.DataBlob = nil

	if auth.TTLSecond == 0 {
		auth.TTLSecond = 300 //default 5minute
	}

	ttl := time.Duration(auth.TTLSecond) * time.Second
	dataStore.Add(auth.Name, ttl, []byte(auth.Value.Cipher))

	gollog.Debug(fmt.Sprintf("SET - '%s' created with '%s'", auth.Name, auth.Value.Key))
	return auth.Exists(dataStore)
}

/*
Get fetchs a value decrypted by Token stored at a Auth-Path in Cache2Go Table.
*/
func (auth *LocalAuth) Get(dataStore DataStore) bool {
	var err error

	if dataStore == nil {
		gollog.Err(fmt.Sprintf("key '%s' asked from corrupted datastore", auth.Name))
		return false
	}
	if auth.Value.Key == nil {
		gollog.Err(fmt.Sprintf("key '%s' asked with missing token", auth.Name))
		return false
	}
	if !auth.Exists(dataStore) {
		gollog.Err(fmt.Sprintf("key '%s doesn't exist", auth.Name))
		return false
	}

	auth.Value.Cipher, err = dataStore.Value(auth.Name)
	if err != nil {
		gollog.Err(fmt.Sprintf("key '%s' asked with wrong token", auth.Name))
		return false
	}

	if err = auth.Value.Decrypt(); err != nil {
		gollog.Err(fmt.Sprintf("failed to decrypt - %q", err.Error()))
		return false
	}

	gollog.Debug(fmt.Sprintf("key '%s' fetched with %s", auth.Name, auth.Value.Key))
	return true
}

/*
Delete purges a Auth-Path in Cache2Go Table, if it's value is decipherable by given Token.
*/
func (auth *LocalAuth) Delete(dataStore DataStore) bool {
	var err error

	if !auth.Get(dataStore) {
		return false
	}

	auth.Value.Cipher = nil
	auth.Value.DataBlob = nil

	err = dataStore.Delete(auth.Name)
	if err != nil {
		gollog.Err(fmt.Sprintf("delete failed for %s because %s", auth.Name, err))
		return false
	}

	gollog.Debug(fmt.Sprintf("key '%s' deleted", auth.Name))
	return true
}
