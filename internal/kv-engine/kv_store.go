package kvengine

import (
	"errors"
	"sync"
)

var (
	ErrEmptyKey    = errors.New("Key can't be empty")
	ErrKeyNotFound = errors.New("Key not found")
)

type Storage interface {
	Put(key string, value string, expires uint64) error
	Get(key string) (string, error)
}

type StorageEntry struct {
	value   string
	expires uint64
}

type MemoryStorage struct {
	values     *sync.Map
	timeSource TimeSource
}

func (m MemoryStorage) Put(key string, value string, expires uint64) error {
	if key == "" {
		return ErrEmptyKey
	}
	if expires == 0 {
		expires = ^uint64(0) // maximum value
	}
	if expires <= 2_502_000 {
		expires += m.timeSource.Now()
	}
	m.values.Store(key, StorageEntry{value, expires})
	return nil
}

func (m MemoryStorage) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	for {
		value, ok := m.values.Load(key)
		if !ok {
			return "", ErrKeyNotFound
		}

		s := value.(StorageEntry)
		if m.timeSource.Now() > s.expires {
			if m.values.CompareAndDelete(key, s) {
				return "", ErrKeyNotFound
			}
		} else {
			return s.value, nil
		}
	}
}

func CreateMemoryStorage() Storage {
	return MemoryStorage{&sync.Map{}, SystemTimeSource{}}
}
