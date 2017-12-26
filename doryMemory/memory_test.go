package doryMemory

import (
	"testing"
	"time"

	"github.com/muesli/cache2go"
	"github.com/stretchr/testify/mock"
)

type MockCacheTable struct {
	mock.Mock
}

type MockCacheItem struct {
	mock.Mock
}

func (m *MockCacheTable) Add(key string, ttl time.Duration, data []byte) *MockCacheItem {
	args := m.Called(key, ttl, data)
	return //args.Bool(0), args.Error(1)
}

func (m *MockCacheTable) Exists(key string) bool {
	args := m.Called(key)
	return args.Bool(0) //, args.Error(1)
}

/*
 */

func TestNewLocalAuthStore(t *testing.T) {
	var localAuthStore interface{} = NewLocalAuthStore("test")
	if _, ok := localAuthStore.(*cache2go.CacheTable); !ok {
		t.Error("Cache2Go CacheTable instantiation failed.")
	}
}

func TestSet(t *testing.T) {
	mockLocalAuthStore := new(MockCacheTable)
	localAuth := LocalAuth{Name: "test", TTLSecond: 1}

	mockLocalAuthStore.On("Add", localAuth.Name, localAuth.TTLSecond).Return()

	localAuth.Set(mockLocalAuthStore)
}

func TestExists(t *testing.T) {
	localAuthStore := NewLocalAuthStore("test")
	locaAuth := LocalAuth{Name: "test"}
	if locaAuth.Exists(localAuthStore) {
		t.Error("LocalAuth Exists gives false positive for key presence.")
	}
}

/*
func TestDelete(t *testing.T) {
	localAuth := NewLocalAuthStore("test")

}

func TestGet(t *testing.T) {
	localAuth := NewLocalAuthStore("test")

}
*/
