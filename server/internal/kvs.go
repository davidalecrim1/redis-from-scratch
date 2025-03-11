package internal

import (
	"fmt"
	"sync"
	"time"
)

var ErrKeyDoesntExist = fmt.Errorf("the provided key doesn't exist")

type KeyValueStorage struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewKeyValueStorage() *KeyValueStorage {
	return &KeyValueStorage{
		data: make(map[string][]byte),
	}
}

func (kv *KeyValueStorage) Set(key, value []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[string(key)] = []byte(value)
	return nil
}

func (kv *KeyValueStorage) Get(key []byte) ([]byte, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	data, ok := kv.data[string(key)]
	if !ok {
		return nil, ErrKeyDoesntExist
	}

	return data, nil
}

func (kv *KeyValueStorage) Delete(key []byte) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	delete(kv.data, string(key))
}

func (kv *KeyValueStorage) SetWithExpiration(key []byte, value []byte, miliseconds int) error {
	if err := kv.Set(key, value); err != nil {
		return err
	}

	go kv.expire(key, miliseconds)
	return nil
}

func (kv *KeyValueStorage) expire(key []byte, miliseconds int) {
	time.Sleep(time.Millisecond * time.Duration(miliseconds))
	kv.Delete(key)
}
