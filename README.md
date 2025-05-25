

# GoRedis 🚀

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
<!-- [![Build Status](https://img.shields.io/github/workflow/status/YOUR_USERNAME/GoRedis/CI)](https://github.com/YOUR_USERNAME/GoRedis/actions) -->
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/GoRedis)](https://goreportcard.com/report/github.com/YOUR_USERNAME/GoRedis)
[![Coverage Status](https://coveralls.io/repos/github/YOUR_USERNAME/GoRedis/badge.svg?branch=main)](https://coveralls.io/github/YOUR_USERNAME/GoRedis?branch=main) -->

A high-performance, lightweight Redis clone built in Go. GoRedis implements the Redis protocol and provides a fast, reliable in-memory data structure store that can be used as a database, cache, and message broker.

## ✨ Features

### Current Features
- 🔧 **RESP Protocol Support** - Full Redis Serialization Protocol implementation
- 💾 **Multiple Data Types** - Strings, Lists, Hashes, Sets, and Sorted Sets
- 🔄 **Persistence** - RDB snapshots and AOF (Append Only File) support
- 📡 **Pub/Sub Messaging** - Real-time message broadcasting
- 🔐 **Security** - Password authentication and ACL support
- 📊 **Monitoring** - INFO, CONFIG, and MONITOR commands
- ⚡ **High Performance** - Optimized for speed with concurrent client handling
- 🔗 **Replication** - Master-slave replication support
- 🎯 **Clustering** - Basic cluster mode with consistent hashing

### Roadmap
- [ ] Streams data type
- [ ] Lua scripting support
- [ ] GeoSpatial commands
- [ ] Full Redis Cluster protocol
- [ ] Redis Sentinel support

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Git

### Installation

#### From Source
```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/GoRedis.git
cd GoRedis

# Build the project
make build

# Or using go build directly
go build -o goredis-server cmd/server/main.go
go build -o goredis-client cmd/client/main.go
```

#### Using Go Install
```bash
go install github.com/YOUR_USERNAME/GoRedis/cmd/server@latest
```

### Running the Server

#### Basic Usage
```bash
# Start with default configuration
./goredis-server

# Start with custom configuration
./goredis-server -c /path/to/redis.conf

# Start with specific port
./goredis-server -p 6380
```

#### Using Docker
```bash
# Build Docker image
docker build -t goredis .

# Run container
docker run -d -p 6379:6379 --name goredis-server goredis

# Run with persistent storage
docker run -d -p 6379:6379 -v /path/to/data:/data --name goredis-server goredis
```

### Connecting to the Server

#### Using GoRedis Client
```bash
./goredis-client
127.0.0.1:6379> SET key "Hello, World!"
OK
127.0.0.1:6379> GET key
"Hello, World!"
```

#### Using Official Redis Client
```bash
redis-cli -p 6379
127.0.0.1:6379> PING
PONG
```

#### Using Go Code
```go
package main

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
)

func main() {
    ctx := context.Background()
    
    // Connect to GoRedis
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    // Set a key
    err := rdb.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }

    // Get a key
    val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key:", val)
}
```

## 📚 Supported Commands

### String Commands
- `SET key value [EX seconds] [PX milliseconds] [NX|XX]`
- `GET key`
- `MGET key [key ...]`
- `MSET key value [key value ...]`
- `INCR key`
- `DECR key`
- `APPEND key value`
- `STRLEN key`

### List Commands
- `LPUSH key value [value ...]`
- `RPUSH key value [value ...]`
- `LPOP key`
- `RPOP key`
- `LLEN key`
- `LRANGE key start stop`
- `LINDEX key index`

### Hash Commands
- `HSET key field value`
- `HGET key field`
- `HMSET key field value [field value ...]`
- `HMGET key field [field ...]`
- `HGETALL key`
- `HDEL key field [field ...]`
- `HLEN key`

### Pub/Sub Commands
- `SUBSCRIBE channel [channel ...]`
- `UNSUBSCRIBE [channel [channel ...]]`
- `PUBLISH channel message`

### Server Commands
- `PING [message]`
- `INFO [section]`
- `CONFIG GET parameter`
- `CONFIG SET parameter value`
- `SAVE`
- `BGSAVE`
- `LASTSAVE`
- `FLUSHDB`
- `FLUSHALL`

[View full command reference →](docs/commands.md)

## ⚙️ Configuration

### Configuration File
Create a `redis.conf` file:

```conf
# Network
bind 127.0.0.1
port 6379
timeout 0
tcp-keepalive 300

# General
daemonize no
loglevel notice
logfile ""
databases 16

# Persistence
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir ./

# Replication
replicaof <masterip> <masterport>
masterauth <master-password>

# Security
requirepass yourpassword

# Limits
maxclients 10000
maxmemory 0
maxmemory-policy noeviction

# Append Only Mode
appendonly no
appendfilename "appendonly.aof"
appendfsync everysec
```

### Environment Variables
```bash
GOREDIS_PORT=6379
GOREDIS_BIND=127.0.0.1
GOREDIS_PASSWORD=secret
GOREDIS_MAXCLIENTS=10000
```

## 🏗️ Architecture

GoRedis follows a modular architecture:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│   Server    │────▶│   Storage   │
└─────────────┘     └─────────────┘     └─────────────┘
                           │                     │
                           ▼                     ▼
                    ┌─────────────┐     ┌─────────────┐
                    │   Protocol  │     │ Persistence │
                    │   (RESP)    │     │  (RDB/AOF)  │
                    └─────────────┘     └─────────────┘
```

### Key Components
- **Server**: Handles client connections and request routing
- **Protocol**: RESP protocol encoding/decoding
- **Storage**: In-memory data structure management
- **Commands**: Command handlers and execution
- **Persistence**: RDB snapshots and AOF logging
- **Replication**: Master-slave synchronization

[Read more about the architecture →](docs/architecture.md)

## 📊 Performance

### Benchmarks
Running on MacBook Pro M1, 16GB RAM:

| Command | Operations/sec | Avg Latency |
|---------|---------------|-------------|
| SET     |      |      |
| GET     |       |      |
| LPUSH   |       |      |
| LPOP    |        |      |
| SADD    |       |      |
| HSET    |       |      |

### Run Benchmarks
```bash
# Using redis-benchmark
redis-benchmark -p 6379 -q -n 100000

# Using Go benchmarks
go test -bench=. -benchmem ./...
```

## 🧪 Testing

### Run Tests
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test
go test -v ./internal/storage/...

# Run integration tests
make test-integration
```

### Test Coverage
- Unit Tests: 85%+ coverage
- Integration Tests: Comprehensive client-server scenarios
- Benchmark Tests: Performance regression testing

## 🛠️ Development

### Project Structure
```
GoRedis/
├── cmd/
│   ├── server/        # Server entry point
│   └── client/        # Client entry point
├── internal/
│   ├── server/        # Server implementation
│   ├── protocol/      # RESP protocol
│   ├── storage/       # Data structures
│   ├── commands/      # Command handlers
│   ├── persistence/   # RDB/AOF
│   ├── pubsub/        # Pub/Sub system
│   └── cluster/       # Clustering
├── pkg/               # Public packages
├── tests/             # Integration tests
├── docs/              # Documentation
└── examples/          # Usage examples
```

### Building from Source
```bash
# Install dependencies
go mod download

# Run tests
make test

# Build binaries
make build

# Run linter
make lint

# Generate documentation
make docs
```


## 🙏 Acknowledgments

- [Redis](https://redis.io/) - The amazing in-memory data structure store that inspired this project
- [go-redis](https://github.com/go-redis/redis) - Redis Go client used for testing
- [RESP Protocol Specification](https://redis.io/docs/reference/protocol-spec/) - Protocol documentation

