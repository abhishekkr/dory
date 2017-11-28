package doryMemory

import (
	"log"
	"path"
	"strconv"
	"time"

	golcrypt "github.com/abhishekkr/gol/golcrypt"
	"github.com/abhishekkr/gol/golenv"
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
	return dataStore.Exists(auth.Name)
}

/*
Set stores an encrypted value with random/provided Token at Auth-Path in Cache2Go Table.
*/
func (auth *LocalAuth) Set(dataStore DataStore) bool {
	if dataStore == nil {
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
	dataStore.Add(auth.Name, ttl, []byte(auth.Value.Cipher))

	return auth.Exists(dataStore)
}

/*
Get fetchs a value decrypted by Token stored at a Auth-Path in Cache2Go Table.
*/
func (auth *LocalAuth) Get(dataStore DataStore) bool {
	var err error

	if dataStore == nil {
		return false
	}
	if auth.Value.Key == nil {
		return false
	}

	if !auth.Exists(dataStore) {
		return false
	}

	auth.Value.Cipher, err = dataStore.Value(auth.Name)

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
func (auth *LocalAuth) Delete(dataStore DataStore) bool {
	var err error

	if dataStore == nil {
		log.Println("delete triggered for missing auth-store")
		return false
	}
	if auth.Value.Key == nil {
		log.Println("delete triggered for empty key")
		return false
	}
	if !auth.Exists(dataStore) {
		log.Println("delete triggered for missing auth identifier")
		return false
	}

	auth.Value.Cipher, err = dataStore.Value(auth.Name)

	if err = auth.Value.Decrypt(); err != nil {
		log.Println("to delete decrypt shall pass;", err)
		return false
	}
	auth.Value.Cipher = nil
	auth.Value.DataBlob = nil

	err = dataStore.Delete(auth.Name)
	if err != nil {
		log.Println("delete triggered but", err)
		return false
	}

	return true
}
