package storage

import (
    "errors"
    "sync"
)

var ErrKeyNotFound = errors.New("key not found")

// MemoryStorage implement the Storage interface
type MemoryStorage struct {
    mu   sync.RWMutex
    data map[string][]byte
}

// NewMemoryStorage create a new memory storage
func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        data: make(map[string][]byte),
    }
}

// Set set the key-value
func (m *MemoryStorage) Set(key string, value []byte) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.data[key] = value
    return nil
}

// Get get the value
func (m *MemoryStorage) Get(key string) ([]byte, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    value, exists := m.data[key]
    if !exists {
        return nil, ErrKeyNotFound
    }
    return value, nil
}

// Delete delete the key
func (m *MemoryStorage) Delete(key string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.data, key)
    return nil
}

// Exists check if the key exists
func (m *MemoryStorage) Exists(key string) bool {
    m.mu.RLock()
    defer m.mu.RUnlock()
    _, exists := m.data[key]
    return exists
}

// Keys get all keys (simple implementation, no support for pattern)
func (m *MemoryStorage) Keys(pattern string) ([]string, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    keys := make([]string, 0, len(m.data))
    for k := range m.data {
        keys = append(keys, k)
    }
    return keys, nil
}

// Clear clear all data
func (m *MemoryStorage) Clear() error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.data = make(map[string][]byte)
    return nil
}