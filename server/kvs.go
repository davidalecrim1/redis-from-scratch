package main

import (
	"fmt"
	"sync"
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
