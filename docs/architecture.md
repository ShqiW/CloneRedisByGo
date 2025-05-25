# GoRedis Architecture Documentation

## Table of Contents

1. [Overview](#overview)
2. [System Architecture](#system-architecture)
3. [Core Components](#core-components)
4. [Data Structures](#data-structures)
5. [Network Layer](#network-layer)
6. [Persistence](#persistence)
7. [Replication](#replication)
8. [Clustering](#clustering)
9. [Performance Optimizations](#performance-optimizations)
10. [Security Architecture](#security-architecture)
11. [Monitoring and Observability](#monitoring-and-observability)

## Overview

GoRedis is a high-performance, Redis-compatible in-memory data structure store built in Go. It implements the Redis protocol while leveraging Go's concurrency model and type safety to provide a robust and scalable solution.

### Design Goals

1. **Redis Protocol Compatibility**: Full support for RESP (Redis Serialization Protocol)
2. **High Performance**: Optimized for low latency and high throughput
3. **Concurrency**: Efficient handling of thousands of concurrent connections
4. **Reliability**: Data persistence and replication support
5. **Scalability**: Horizontal scaling through clustering
6. **Maintainability**: Clean, modular architecture

### Key Design Decisions

- **Language Choice**: Go provides excellent concurrency primitives, garbage collection, and performance
- **Memory Model**: All data structures are kept in memory for maximum performance
- **Threading Model**: One goroutine per client connection with shared data structures
- **Persistence Strategy**: Dual approach with RDB snapshots and AOF logging
- **Networking**: Epoll/Kqueue based event loop for efficient I/O handling

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                          GoRedis Server                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │
│  │   Client    │  │   Client    │  │   Client    │           │
│  │ Connection  │  │ Connection  │  │ Connection  │  ...       │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘           │
│         │                 │                 │                   │
│  ┌──────┴─────────────────┴─────────────────┴──────┐          │
│  │            Connection Manager (Epoll)             │          │
│  └──────────────────────┬───────────────────────────┘          │
│                         │                                       │
│  ┌──────────────────────┴───────────────────────────┐          │
│  │              Command Dispatcher                   │          │
│  └──────┬───────────────────────────────────┬───────┘          │
│         │                                   │                   │
│  ┌──────┴──────┐                    ┌──────┴──────┐           │
│  │   Command   │                    │   Command   │           │
│  │  Handlers   │                    │   Queue     │           │
│  └──────┬──────┘                    └──────┬──────┘           │
│         │                                   │                   │
│  ┌──────┴───────────────────────────────────┴──────┐          │
│  │              Storage Engine                      │          │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐  │          │
│  │  │ String │ │  List  │ │  Hash  │ │  Set   │  │          │
│  │  └────────┘ └────────┘ └────────┘ └────────┘  │          │
│  └──────────────────┬───────────────────────────────┘          │
│                     │                                           │
│  ┌─────────────────┴────────────────────────────┐             │
│  │          Persistence Layer                     │             │
│  │  ┌─────────────┐    ┌─────────────┐         │             │
│  │  │     RDB     │    │     AOF     │         │             │
│  │  └─────────────┘    └─────────────┘         │             │
│  └───────────────────────────────────────────────┘             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Component Interaction Flow

```
Client Request → TCP Connection → Connection Handler → RESP Parser
    ↓
Command Dispatcher → Command Validator → Command Handler
    ↓
Storage Engine → Data Structure Operations
    ↓
Response Builder → RESP Encoder → TCP Response → Client
```

## Core Components

### 1. Server Component (`internal/server/server.go`)

The main server component manages the lifecycle of the GoRedis instance.

```go
type Server struct {
    // Configuration
    config     *Config
    
    // Networking
    listener   net.Listener
    clients    map[int64]*Client
    clientsMu  sync.RWMutex
    
    // Storage
    storage    storage.Storage
    
    // Command handling
    commands   map[string]CommandFunc
    
    // Persistence
    rdb        *persistence.RDB
    aof        *persistence.AOF
    
    // Replication
    role       string // "master" or "slave"
    slaves     []*SlaveConn
    master     *MasterConn
    
    // Stats
    stats      *ServerStats
    
    // Lifecycle
    shutdown   chan bool
    wg         sync.WaitGroup
}
```

**Responsibilities:**
- Initialize and manage all subsystems
- Accept client connections
- Coordinate shutdown procedures
- Maintain server statistics

### 2. Connection Manager (`internal/server/connection.go`)

Handles client connections efficiently using epoll/kqueue.

```go
type ConnectionManager struct {
    epoll      *Epoll
    conns      map[int]*ClientConn
    bufferPool *sync.Pool
}

type ClientConn struct {
    id         int64
    conn       net.Conn
    reader     *bufio.Reader
    writer     *bufio.Writer
    
    // Authentication
    authenticated bool
    user          *User
    
    // Transaction state
    multi      bool
    queue      []Command
    watching   map[string]uint64
    
    // Pub/Sub state
    subscriber bool
    channels   map[string]bool
    patterns   map[string]bool
}
```

**Key Features:**
- Non-blocking I/O with epoll/kqueue
- Connection pooling and reuse
- Buffer management for efficient memory usage
- Per-connection state management

### 3. RESP Protocol (`internal/protocol/resp.go`)

Implements the Redis Serialization Protocol for client-server communication.

```go
type Value struct {
    Type  ValueType
    Str   string
    Num   int
    Bulk  string
    Array []Value
    Error error
}

type Parser struct {
    reader *bufio.Reader
}

type Encoder struct {
    writer *bufio.Writer
}
```

**Protocol Types:**
- Simple Strings: `+OK\r\n`
- Errors: `-Error message\r\n`
- Integers: `:1000\r\n`
- Bulk Strings: `$6\r\nfoobar\r\n`
- Arrays: `*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n`

### 4. Command Dispatcher (`internal/commands/handler.go`)

Routes commands to appropriate handlers.

```go
type CommandHandler struct {
    storage    storage.Storage
    pubsub     *pubsub.PubSub
    
    // Command registry
    commands   map[string]CommandFunc
    
    // Command metadata
    metadata   map[string]CommandMeta
}

type CommandFunc func(conn *ClientConn, args []Value) Value

type CommandMeta struct {
    Name       string
    Arity      int  // -n means >= n args
    Flags      []string
    FirstKey   int
    LastKey    int
    KeyStep    int
}
```

**Command Processing Pipeline:**
1. Parse RESP command
2. Validate command existence and arity
3. Check authentication and ACL permissions
4. Execute pre-command hooks (e.g., key watching)
5. Call command handler
6. Execute post-command hooks (e.g., replication, AOF)
7. Send response

### 5. Storage Engine (`internal/storage/storage.go`)

The core data storage abstraction.

```go
type Storage interface {
    // String operations
    Set(key string, value string, opts ...SetOption) error
    Get(key string) (string, error)
    
    // List operations
    LPush(key string, values ...string) (int, error)
    LPop(key string) (string, error)
    
    // Hash operations
    HSet(key, field string, value string) error
    HGet(key, field string) (string, error)
    
    // Set operations
    SAdd(key string, members ...string) (int, error)
    SMembers(key string) ([]string, error)
    
    // Sorted set operations
    ZAdd(key string, members ...ZMember) (int, error)
    ZRange(key string, start, stop int) ([]string, error)
    
    // Key operations
    Del(keys ...string) int
    Exists(keys ...string) int
    Expire(key string, seconds int) bool
    TTL(key string) int
    
    // Database operations
    FlushDB() error
    DBSize() int
}
```

## Data Structures

### Memory Layout

GoRedis uses optimized data structures for each Redis type:

#### String Storage
```go
type StringObject struct {
    value []byte
    encoding int // RAW, INT, EMBSTR
}
```
- Small strings (<= 44 bytes) use embedded string optimization
- Integer strings are stored as native integers
- Large strings use dynamic allocation

#### List Storage
```go
type ListObject struct {
    encoding int // LINKEDLIST or ZIPLIST
    list     interface{}
}

type LinkedList struct {
    head   *ListNode
    tail   *ListNode
    length int
}

type ZipList struct {
    data []byte
    size int
}
```
- Small lists use ziplist (compressed representation)
- Large lists use doubly-linked lists
- Automatic conversion based on size thresholds

#### Hash Storage
```go
type HashObject struct {
    encoding int // ZIPLIST or HT
    hash     interface{}
}

type HashTable struct {
    buckets  []*HashEntry
    size     int
    used     int
    sizemask int
}
```
- Small hashes use ziplist
- Large hashes use hash tables with chain collision resolution
- Progressive rehashing for smooth growth

#### Set Storage
```go
type SetObject struct {
    encoding int // INTSET or HT
    set      interface{}
}

type IntSet struct {
    encoding int // INT16, INT32, or INT64
    data     []byte
    length   int
}
```
- Integer-only sets use intset (sorted array)
- Mixed sets use hash tables
- Automatic encoding upgrade when needed

#### Sorted Set Storage
```go
type ZSetObject struct {
    dict     *HashTable  // member -> score mapping
    skiplist *SkipList   // score-ordered structure
}

type SkipList struct {
    header *SkipListNode
    tail   *SkipListNode
    length int
    level  int
}
```
- Dual representation for O(1) score lookups and O(log N) range queries
- Skip list provides efficient ordered traversal
- Dictionary provides fast member lookups

### Memory Optimization Strategies

1. **Object Sharing**: Common small integers (-128 to 127) are shared
2. **Lazy Deletion**: Large objects deleted in background
3. **Memory Pools**: Reusable buffers for network I/O
4. **Reference Counting**: Efficient memory management
5. **Compression**: Automatic compression for suitable data structures

## Network Layer

### Event-Driven Architecture

```go
type EventLoop struct {
    epoll       *Epoll
    timers      *TimerQueue
    fileEvents  []FileEvent
    timeEvents  []TimeEvent
    stop        bool
}

func (el *EventLoop) Run() {
    for !el.stop {
        // Wait for events with timeout
        events := el.epoll.Wait(el.timers.NextTimeout())
        
        // Process file events
        for _, event := range events {
            el.processFileEvent(event)
        }
        
        // Process timer events
        el.processTimerEvents()
    }
}
```

### Connection Handling

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Accept    │────▶│   Read      │────▶│   Parse     │
│   Socket    │     │   Request   │     │   RESP      │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                    ┌─────────────┐            ▼
                    │   Write     │     ┌─────────────┐
                    │   Response  │◀────│   Execute   │
                    └─────────────┘     │   Command   │
                                        └─────────────┘
```

### Buffer Management

```go
type BufferPool struct {
    small  sync.Pool // 4KB buffers
    medium sync.Pool // 16KB buffers
    large  sync.Pool // 64KB buffers
}

func (bp *BufferPool) Get(size int) []byte {
    switch {
    case size <= 4096:
        return bp.small.Get().([]byte)
    case size <= 16384:
        return bp.medium.Get().([]byte)
    default:
        return bp.large.Get().([]byte)
    }
}
```

## Persistence

### RDB (Redis Database) Snapshots

RDB provides point-in-time snapshots of the dataset.

```go
type RDB struct {
    filename   string
    tempFile   string
    
    // Snapshot scheduling
    saveParams []SaveParam
    
    // Progress tracking
    lastSave   time.Time
    dirtyCount int
}

type SaveParam struct {
    seconds int
    changes int
}
```

**RDB File Format:**
```
+-------+-------------+----------+---------------+-----+-------+
| REDIS | RDB-VERSION | AUX-DATA | DB-SELECTOR   | KEY | VALUE |
+-------+-------------+----------+---------------+-----+-------+
| 5bytes| 4 bytes     | Variable | 1-5 bytes     | Var | Var   |
+-------+-------------+----------+---------------+-----+-------+
```

**Saving Process:**
1. Fork child process (or use goroutine)
2. Child writes dataset to temporary file
3. Use copy-on-write for efficiency
4. Atomic rename when complete

### AOF (Append Only File)

AOF logs every write operation for durability.

```go
type AOF struct {
    file       *os.File
    filename   string
    
    // Sync policy
    syncPolicy int // ALWAYS, EVERYSEC, NO
    
    // Buffer management
    buffer     []byte
    
    // Rewrite state
    rewriting  bool
    rewriteBuf []byte
}
```

**AOF Rewrite Process:**
1. Fork child process
2. Child writes compacted dataset
3. Parent buffers new writes
4. Merge parent buffer with child output
5. Atomic file replacement

### Persistence Configuration

```yaml
# RDB Configuration
save 900 1      # Save after 900s if 1+ keys changed
save 300 10     # Save after 300s if 10+ keys changed
save 60 10000   # Save after 60s if 10000+ keys changed

rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb

# AOF Configuration
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
```

## Replication

### Master-Slave Architecture

```
┌──────────────┐
│    Master    │
│  (Primary)   │
└──────┬───────┘
       │
   ┌───┴───┐
   │       │
┌──┴───┐ ┌─┴────┐
│Slave │ │Slave │
│  1   │ │  2   │
└──────┘ └──────┘
```

### Replication Protocol

1. **Initial Sync:**
   ```
   SLAVE: PSYNC ? -1
   MASTER: +FULLRESYNC <runid> <offset>
   MASTER: $<rdb-size>
   MASTER: <rdb-data>
   ```

2. **Command Propagation:**
   ```go
   type ReplicationBuffer struct {
       buffer   []byte
       offset   int64
       capacity int
   }
   
   func (m *Master) PropagateCommand(cmd Command) {
       encoded := cmd.Encode()
       
       for _, slave := range m.slaves {
           slave.Send(encoded)
       }
       
       m.replBuffer.Append(encoded)
   }
   ```

3. **Partial Resynchronization:**
   - Slaves track replication offset
   - Master maintains replication backlog
   - Reconnecting slaves can resume from offset

### Replication Features

- **Cascading Replication**: Slaves can have sub-slaves
- **Read-Only Slaves**: Slaves reject write commands by default
- **Replication Timeout**: Detect and handle network partitions
- **Min-Slaves Configuration**: Refuse writes without enough slaves

## Clustering

### Cluster Topology

```
┌─────────────────────────────────────────────┐
│              16384 Hash Slots               │
├─────────────┬─────────────┬────────────────┤
│  0 - 5460   │ 5461 - 10922│ 10923 - 16383  │
│   Node A    │    Node B   │    Node C      │
└─────────────┴─────────────┴────────────────┘
```

### Cluster Implementation

```go
type ClusterNode struct {
    name     string
    ip       string
    port     int
    flags    uint16
    
    slots    []SlotRange
    
    // Failure detection
    pingTime time.Time
    pongTime time.Time
    
    // Connections
    link     *ClusterLink
}

type Cluster struct {
    myself    *ClusterNode
    nodes     map[string]*ClusterNode
    
    // Slot mapping
    slots     [16384]*ClusterNode
    
    // Cluster state
    state     int // OK or FAIL
    
    // Epoch for configuration
    epoch     uint64
}
```

### Key Distribution

```go
func KeyHashSlot(key string) uint16 {
    // Extract hash tag if present
    start := strings.IndexByte(key, '{')
    if start != -1 {
        end := strings.IndexByte(key[start+1:], '}')
        if end != -1 {
            key = key[start+1 : start+1+end]
        }
    }
    
    // CRC16 with specific polynomial
    return crc16(key) & 0x3FFF
}
```

### Cluster Communication

**Cluster Bus Protocol:**
- Gossip-based failure detection
- Configuration propagation
- Heartbeat messages (PING/PONG)
- Failure reports (FAIL)
- Configuration updates (UPDATE)

## Performance Optimizations

### 1. Lock-Free Data Structures

```go
type LockFreeList struct {
    head atomic.Pointer[Node]
}

func (l *LockFreeList) Push(value interface{}) {
    newNode := &Node{value: value}
    for {
        head := l.head.Load()
        newNode.next = head
        if l.head.CompareAndSwap(head, newNode) {
            break
        }
    }
}
```

### 2. Sharded Locks

```go
type ShardedMutex struct {
    shards [256]sync.RWMutex
}

func (sm *ShardedMutex) Lock(key string) {
    shard := hash(key) & 0xFF
    sm.shards[shard].Lock()
}
```

### 3. Memory Pooling

```go
var pools = map[int]*sync.Pool{
    64:   {New: func() interface{} { return make([]byte, 64) }},
    512:  {New: func() interface{} { return make([]byte, 512) }},
    4096: {New: func() interface{} { return make([]byte, 4096) }},
}
```

### 4. Zero-Copy Operations

```go
// Use unsafe to avoid copying
func StringToBytes(s string) []byte {
    return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
    return unsafe.String(&b[0], len(b))
}
```

### 5. Batch Operations

```go
type BatchProcessor struct {
    commands []Command
    results  []Result
    mu       sync.Mutex
}

func (bp *BatchProcessor) Process() {
    bp.mu.Lock()
    defer bp.mu.Unlock()
    
    // Process all commands in single lock acquisition
    for i, cmd := range bp.commands {
        bp.results[i] = cmd.Execute()
    }
}
```

### 6. CPU Cache Optimization

```go
// Align structures to cache lines
type CacheLinePadded struct {
    value uint64
    _     [7]uint64 // 64 bytes total
}
```

## Security Architecture

### Authentication

```go
type AuthManager struct {
    requirePass string
    users       map[string]*User
}

type User struct {
    username     string
    passwords    []string
    flags        UserFlags
    allowedCmds  map[string]bool
    deniedCmds   map[string]bool
    keyPatterns  []string
}
```

### ACL (Access Control List)

```go
type ACL struct {
    users    map[string]*User
    defaultUser *User
    
    // ACL log
    log      []ACLLogEntry
    logSize  int
}

type ACLLogEntry struct {
    timestamp time.Time
    username  string
    command   string
    reason    string
}
```

### TLS Support

```go
type TLSConfig struct {
    CertFile       string
    KeyFile        string
    CAFile         string
    
    ClientAuth     bool
    MinVersion     uint16
    CipherSuites   []uint16
}

func (s *Server) ListenAndServeTLS(config *TLSConfig) error {
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientAuth:   tls.RequireAndVerifyClientCert,
        MinVersion:   tls.VersionTLS12,
    }
    
    listener := tls.NewListener(s.listener, tlsConfig)
    return s.Serve(listener)
}
```

## Monitoring and Observability

### Metrics Collection

```go
type Metrics struct {
    // Commands
    commandsProcessed   atomic.Uint64
    commandsPerSec      atomic.Uint64
    
    // Connections
    connectionsReceived atomic.Uint64
    connectionsActive   atomic.Int64
    
    // Memory
    memoryUsed         atomic.Uint64
    memoryPeak         atomic.Uint64
    
    // Keys
    keysTotal          atomic.Uint64
    keysExpired        atomic.Uint64
    keysEvicted        atomic.Uint64
    
    // Network
    bytesReceived      atomic.Uint64
    bytesSent          atomic.Uint64
    
    // Persistence
    lastSaveTime       atomic.Int64
    rdbSaveInProgress  atomic.Bool
    aofRewriteInProgress atomic.Bool
}
```

### INFO Command Implementation

```go
func (s *Server) GetInfo(section string) map[string]string {
    info := make(map[string]string)
    
    switch section {
    case "server":
        info["redis_version"] = VERSION
        info["process_id"] = strconv.Itoa(os.Getpid())
        info["tcp_port"] = strconv.Itoa(s.config.Port)
        info["uptime_in_seconds"] = strconv.FormatInt(time.Since(s.startTime).Seconds(), 10)
        
    case "clients":
        info["connected_clients"] = strconv.Itoa(len(s.clients))
        info["blocked_clients"] = strconv.Itoa(s.blockedClients)
        
    case "memory":
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        info["used_memory"] = strconv.FormatUint(m.Alloc, 10)
        info["used_memory_human"] = humanizeBytes(m.Alloc)
        info["used_memory_peak"] = strconv.FormatUint(s.memoryPeak, 10)
        
    case "stats":
        info["total_connections_received"] = strconv.FormatUint(s.metrics.connectionsReceived.Load(), 10)
        info["total_commands_processed"] = strconv.FormatUint(s.metrics.commandsProcessed.Load(), 10)
        info["instantaneous_ops_per_sec"] = strconv.FormatUint(s.metrics.commandsPerSec.Load(), 10)
    }
    
    return info
}
```

### Slow Query Log

```go
type SlowLog struct {
    entries    []SlowLogEntry
    maxLen     int
    threshold  time.Duration
    mu         sync.RWMutex
}

type SlowLogEntry struct {
    id         uint64
    timestamp  time.Time
    duration   time.Duration
    command    []string
    clientAddr string
}

func (s *SlowLog) Add(cmd []string, duration time.Duration, client string) {
    if duration < s.threshold {
        return
    }
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    entry := SlowLogEntry{
        id:         atomic.AddUint64(&s.nextID, 1),
        timestamp:  time.Now(),
        duration:   duration,
        command:    cmd,
        clientAddr: client,
    }
    
    s.entries = append(s.entries, entry)
    if len(s.entries) > s.maxLen {
        s.entries = s.entries[1:]
    }
}
```

### Event Notification

```go
type EventNotifier struct {
    subscribers map[string][]chan Event
    mu          sync.RWMutex
}

type Event struct {
    Type      string // "expired", "evicted", "del"
    Key       string
    Database  int
    Timestamp time.Time
}

func (en *EventNotifier) Notify(event Event) {
    en.mu.RLock()
    defer en.mu.RUnlock()
    
    pattern := fmt.Sprintf("__key%s@%d__:%s", event.Type, event.Database, event.Key)
    
    for _, ch := range en.subscribers[pattern] {
        select {
        case ch <- event:
        default: // Don't block
        }
    }
}
```

## Advanced Features

### Lua Scripting Engine

```go
type ScriptEngine struct {
    vm       *lua.LState
    scripts  map[string]*Script
    sha1sum  func([]byte) string
}

type Script struct {
    SHA1     string
    Source   string
    Compiled *lua.LFunction
}

func (se *ScriptEngine) Eval(script string, keys []string, args []string) (interface{}, error) {
    L := se.vm
    
    // Set up KEYS and ARGV tables
    keysTable := L.NewTable()
    for i, key := range keys {
        L.RawSetInt(keysTable, i+1, lua.LString(key))
    }
    L.SetGlobal("KEYS", keysTable)
    
    argsTable := L.NewTable()
    for i, arg := range args {
        L.RawSetInt(argsTable, i+1, lua.LString(arg))
    }
    L.SetGlobal("ARGV", argsTable)
    
    // Execute script
    if err := L.DoString(script); err != nil {
        return nil, err
    }
    
    // Get return value
    ret := L.Get(-1)
    L.Pop(1)
    
    return luaValueToGo(ret), nil
}
```

### Transactions

```go
type Transaction struct {
    commands []Command
    watched  map[string]uint64
}

func (t *Transaction) Execute(storage Storage) []Result {
    // Check watched keys
    for key, version := range t.watched {
        if storage.GetVersion(key) != version {
            return []Result{{Err: ErrWatchedKeyModified}}
        }
    }
    
    // Execute all commands atomically
    storage.Lock()
    defer storage.Unlock()
    
    results := make([]Result, len(t.commands))
    for i, cmd := range t.commands {
        results[i] = cmd.Execute()
    }
    
    return results
}
```

### Pipeline Support

```go
type Pipeline struct {
    conn     *ClientConn
    commands []Command
    results  []Result
}

func (p *Pipeline) Execute() error {
    // Send all commands without waiting for responses
    for _, cmd := range p.commands {
        if err := p.conn.WriteCommand(cmd); err != nil {
            return err
        }
    }
    
    // Flush the write buffer
    if err := p.conn.Flush(); err != nil {
        return err
    }
    
    // Read all responses
    p.results = make([]Result, len(p.commands))
    for i := range p.commands {
        resp, err := p.conn.ReadResponse()
        if err != nil {
            return err
        }
        p.results[i] = resp
    }
    
    return nil
}
```

### Memory Management

```go
type MemoryManager struct {
    maxMemory    uint64
    policy       EvictionPolicy
    
    // LRU/LFU tracking
    accessTime   map[string]time.Time
    accessCount  map[string]uint64
    
    // Memory pressure callback
    onEviction   func(key string)
}

type EvictionPolicy int

const (
    NoEviction EvictionPolicy = iota
    AllKeysLRU
    VolatileLRU
    AllKeysLFU
    VolatileLFU
    AllKeysRandom
    VolatileRandom
    VolatileTTL
)

func (mm *MemoryManager) EvictKeys(bytesNeeded uint64) {
    switch mm.policy {
    case AllKeysLRU:
        mm.evictLRU(bytesNeeded, false)
    case VolatileLRU:
        mm.evictLRU(bytesNeeded, true)
    case AllKeysLFU:
        mm.evictLFU(bytesNeeded, false)
    // ... other policies
    }
}
```

## Testing Architecture

### Unit Testing Strategy

```go
// Test individual components in isolation
func TestRESPParser(t *testing.T) {
    tests := []struct {
        input    string
        expected Value
    }{
        {"+OK\r\n", Value{Type: String, Str: "OK"}},
        {"$6\r\nfoobar\r\n", Value{Type: Bulk, Bulk: "foobar"}},
        {"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", Value{Type: Array, Array: []Value{...}}},
    }
    
    for _, tt := range tests {
        parser := NewParser(strings.NewReader(tt.input))
        got, err := parser.Parse()
        assert.NoError(t, err)
        assert.Equal(t, tt.expected, got)
    }
}
```

### Integration Testing

```go
func TestServerIntegration(t *testing.T) {
    // Start test server
    server := NewTestServer()
    defer server.Close()
    
    // Connect multiple clients
    clients := make([]*redis.Client, 10)
    for i := range clients {
        clients[i] = redis.NewClient(&redis.Options{
            Addr: server.Addr,
        })
    }
    
    // Run concurrent operations
    var wg sync.WaitGroup
    for i, client := range clients {
        wg.Add(1)
        go func(id int, c *redis.Client) {
            defer wg.Done()
            
            // Perform operations
            key := fmt.Sprintf("key%d", id)
            err := c.Set(ctx, key, "value", 0).Err()
            assert.NoError(t, err)
            
            val, err := c.Get(ctx, key).Result()
            assert.NoError(t, err)
            assert.Equal(t, "value", val)
        }(i, client)
    }
    
    wg.Wait()
}
```

### Benchmarking

```go
func BenchmarkSET(b *testing.B) {
    server := NewTestServer()
    defer server.Close()
    
    client := redis.NewClient(&redis.Options{
        Addr: server.Addr,
    })
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            key := fmt.Sprintf("key%d", i)
            if err := client.Set(ctx, key, "value", 0).Err(); err != nil {
                b.Fatal(err)
            }
            i++
        }
    })
    
    b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "ops/s")
}
```

### Chaos Testing

```go
type ChaosTest struct {
    server    *Server
    clients   []*redis.Client
    
    // Chaos parameters
    networkDelay   time.Duration
    packetLoss     float64
    nodeFailures   int
}

func (ct *ChaosTest) Run() {
    // Start background load
    go ct.generateLoad()
    
    // Introduce failures
    go ct.introduceNetworkChaos()
    go ct.simulateNodeFailures()
    
    // Verify consistency
    ct.verifyDataConsistency()
    ct.verifyReplication()
    ct.verifyClusterState()
}
```

## Deployment Considerations

### Docker Support

```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o goredis cmd/server/main.go

# Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/goredis .
COPY configs/redis.conf .

EXPOSE 6379
CMD ["./goredis", "-c", "redis.conf"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: goredis
spec:
  serviceName: goredis
  replicas: 3
  selector:
    matchLabels:
      app: goredis
  template:
    metadata:
      labels:
        app: goredis
    spec:
      containers:
      - name: goredis
        image: goredis:latest
        ports:
        - containerPort: 6379
        - containerPort: 16379  # Cluster bus
        volumeMounts:
        - name: data
          mountPath: /data
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
```

### Monitoring Stack

```yaml
# Prometheus configuration
scrape_configs:
  - job_name: 'goredis'
    static_configs:
    - targets: ['goredis:9121']
    metrics_path: '/metrics'
```

## Performance Tuning Guide

### System Tuning

```bash
# Kernel parameters (sysctl.conf)
net.core.somaxconn = 65535
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.ip_local_port_range = 1024 65535
net.core.netdev_max_backlog = 65535
vm.overcommit_memory = 1

# File descriptors
ulimit -n 65535
```

### Configuration Optimization

```yaml
# Performance-oriented configuration
tcp-backlog 511
timeout 0
tcp-keepalive 300

# Memory optimization
maxmemory 8gb
maxmemory-policy allkeys-lru
maxmemory-samples 5

# Persistence tuning
save ""  # Disable RDB for pure cache mode
appendonly no  # Disable AOF for pure cache mode

# Client output buffer limits
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit replica 256mb 64mb 60
client-output-buffer-limit pubsub 32mb 8mb 60
```

## Future Enhancements

### Planned Features

1. **Redis Modules API**
   - Plugin system for custom commands
   - Dynamic library loading
   - Module isolation and sandboxing

2. **Active-Active Replication**
   - Multi-master support
   - Conflict-free replicated data types (CRDTs)
   - Cross-datacenter replication

3. **Advanced Data Types**
   - Streams (Redis Streams compatible)
   - Bloom filters
   - HyperLogLog
   - Geospatial indexes

4. **Performance Improvements**
   - io_uring support for Linux
   - SIMD optimizations
   - GPU acceleration for certain operations

5. **Enterprise Features**
   - LDAP/AD integration
   - Audit logging
   - Data encryption at rest
   - Multi-tenancy support

## Conclusion

GoRedis demonstrates how Go's concurrency primitives and type safety can be leveraged to build a high-performance, Redis-compatible data store. The architecture emphasizes:

- **Modularity**: Clear separation of concerns
- **Performance**: Optimized data structures and algorithms
- **Reliability**: Comprehensive persistence and replication
- **Scalability**: Horizontal scaling through clustering
- **Maintainability**: Clean code and extensive testing

The project serves as both a production-ready Redis alternative and an educational resource for understanding the internals of in-memory data stores.