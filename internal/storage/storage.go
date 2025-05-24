package storage

// Storage define the key-value storage interface
type Storage interface {
    // basic operation
    Set(key string, value []byte) error
    Get(key string) ([]byte, error)
    Delete(key string) error
    Exists(key string) bool
    
    // batch operation
    Keys(pattern string) ([]string, error)
    Clear() error
}