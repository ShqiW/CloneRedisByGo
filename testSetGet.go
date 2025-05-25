package main

import (
    "fmt"
    "github.com/ShqiW/CloneRedisByGo/internal/storage"   // 改成你的实际模块名
    "github.com/ShqiW/CloneRedisByGo/internal/commands"  // 改成你的实际模块名
)

func main() {
    // 创建内存存储
    store := storage.NewMemoryStorage()
    
    // 创建命令处理器
    handler := commands.NewHandler(store)
    
    fmt.Println("=== GoRedis SET/GET 测试 ===")
    
    // 测试 SET 命令
    fmt.Println("\n测试 SET 命令:")
    result := handler.Execute([]string{"SET", "name", "GoRedis"})
    fmt.Printf("SET name GoRedis => %s", result)
    
    result = handler.Execute([]string{"SET", "version", "1.0"})
    fmt.Printf("SET version 1.0 => %s", result)
    
    // 测试 GET 命令
    fmt.Println("\n测试 GET 命令:")
    result = handler.Execute([]string{"GET", "name"})
    fmt.Printf("GET name => %s", result)
    
    result = handler.Execute([]string{"GET", "version"})
    fmt.Printf("GET version => %s", result)
    
    // 测试不存在的键
    result = handler.Execute([]string{"GET", "notexist"})
    fmt.Printf("GET notexist => %s", result)
    
    // 测试 PING
    fmt.Println("\n测试 PING 命令:")
    result = handler.Execute([]string{"PING"})
    fmt.Printf("PING => %s", result)
}