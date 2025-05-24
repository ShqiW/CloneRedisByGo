package main

import (
    "fmt"
    "github.com/ShqiW/CloneRedisByGo/internal/storage"
)


func HandleSet(store storage.Storage, key, value string) string {
    err := store.Set(key, value)
    if err != nil {
        return fmt.Sprintf("-ERR %s\r\n", err.Error())
    }
    return "+OK\r\n"
}

func HandleGet(store storage.Storage, key string) string {
    value, exists := store.Get(key)
    if !exists {
        return "$-1\r\n"  // nil
    }
    return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
}

func main() {
    // create memory storage
    store := storage.NewMemoryStorage()
    
    // test SET
    fmt.Println("测试 SET 命令:")
    result := HandleSet(store, "name", "GoRedis")
    fmt.Printf("SET name GoRedis => %s", result)
    
    // test GET
    fmt.Println("\n测试 GET 命令:")
    result = HandleGet(store, "name")
    fmt.Printf("GET name => %s", result)
    
    // test non-existent key
    result = HandleGet(store, "notexist")
    fmt.Printf("GET notexist => %s", result)
}