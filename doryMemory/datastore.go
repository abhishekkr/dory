package doryMemory

import "time"

/*
package for Create/Read/Delete actions over Cache2Go Table for items post aes encryption
*/

/*
DataStore is an interface for all datastore backends that can be used.
Mainly so I can write actual "unit" tests.
*/
type DataStore interface {
	Add(string, time.Duration, []byte) error
	Delete(string) error
	Exists(string) bool
	Value(string) ([]byte, error)

	List() []string
	Count() int
	Purge() error
	PurgeOne(string) error
}
