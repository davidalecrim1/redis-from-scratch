package main

import (
	"fmt"
	"sync"
)

var ErrInvalidKey = fmt.Errorf("the provided key is invalid")

type KeyValueStorage struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewKeyValueStorage() *KeyValueStorage {
	return &KeyValueStorage{
		data: make(map[string][]byte),
	}
}

func (kv *KeyValueStorage) Set(key, value string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[key] = []byte(value)
	return nil
}

func (kv *KeyValueStorage) Get(key string) ([]byte, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	data, ok := kv.data[key]
	if !ok {
		return nil, ErrInvalidKey
	}

	return data, nil
}
