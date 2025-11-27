package testutils

import (
	"ASMO-site-backend/internal/cache"
	"encoding/json"
	"sync"
	"time"
)

// RedisMock реализует cache.Cache интерфейс
type RedisMock struct {
	data  map[string][]byte
	mutex sync.RWMutex
}

var _ cache.Cache = (*RedisMock)(nil)

func NewRedisMock() *RedisMock {
	return &RedisMock{
		data: make(map[string][]byte),
	}
}

func (r *RedisMock) Set(key string, value interface{}, expiration time.Duration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	r.data[key] = jsonData
	return nil
}

func (r *RedisMock) Get(key string, dest interface{}) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	jsonData, exists := r.data[key]
	if !exists {
		return ErrNotFound
	}

	return json.Unmarshal(jsonData, dest)
}

func (r *RedisMock) Delete(key string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.data, key)
	return nil
}

func (r *RedisMock) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.data = make(map[string][]byte)
	return nil
}

func (r *RedisMock) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.data = make(map[string][]byte)
}