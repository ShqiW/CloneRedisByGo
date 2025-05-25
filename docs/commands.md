# GoRedis Command Reference

This document provides a comprehensive reference for all commands supported by GoRedis, organized by data type and functionality.

## Table of Contents

- [Connection Commands](#connection-commands)
- [String Commands](#string-commands)
- [List Commands](#list-commands)
- [Hash Commands](#hash-commands)
- [Set Commands](#set-commands)
- [Sorted Set Commands](#sorted-set-commands)
- [Pub/Sub Commands](#pubsub-commands)
- [Transaction Commands](#transaction-commands)
- [Server Commands](#server-commands)
- [Persistence Commands](#persistence-commands)
- [Replication Commands](#replication-commands)
- [Cluster Commands](#cluster-commands)

## Connection Commands

### PING
Test the connection to the server.
```
PING [message]
```
- **Time Complexity:** O(1)
- **Returns:** PONG or the message if provided
- **Example:**
  ```
  > PING
  PONG
  > PING "Hello"
  "Hello"
  ```

### AUTH
Authenticate to the server.
```
AUTH password
```
- **Time Complexity:** O(1)
- **Returns:** OK or error
- **Example:**
  ```
  > AUTH mypassword
  OK
  ```

### SELECT
Select the Redis database.
```
SELECT index
```
- **Time Complexity:** O(1)
- **Returns:** OK
- **Example:**
  ```
  > SELECT 1
  OK
  ```

### QUIT
Close the connection.
```
QUIT
```
- **Time Complexity:** O(1)
- **Returns:** OK
- **Example:**
  ```
  > QUIT
  OK
  ```

## String Commands

### SET
Set key to hold the string value.
```
SET key value [EX seconds] [PX milliseconds] [NX|XX] [KEEPTTL]
```
- **Time Complexity:** O(1)
- **Options:**
  - `EX seconds`: Set expire time in seconds
  - `PX milliseconds`: Set expire time in milliseconds
  - `NX`: Only set if key doesn't exist
  - `XX`: Only set if key exists
  - `KEEPTTL`: Retain the TTL associated with the key
- **Returns:** OK or nil
- **Example:**
  ```
  > SET mykey "Hello"
  OK
  > SET mykey "World" NX
  (nil)
  > SET mykey "World" XX
  OK
  > SET tempkey "temp" EX 10
  OK
  ```

### GET
Get the value of key.
```
GET key
```
- **Time Complexity:** O(1)
- **Returns:** String value or nil
- **Example:**
  ```
  > GET mykey
  "World"
  > GET nonexisting
  (nil)
  ```

### MGET
Get values of multiple keys.
```
MGET key [key ...]
```
- **Time Complexity:** O(N) where N is number of keys
- **Returns:** Array of values
- **Example:**
  ```
  > MGET key1 key2 key3
  1) "value1"
  2) "value2"
  3) (nil)
  ```

### MSET
Set multiple keys to multiple values.
```
MSET key value [key value ...]
```
- **Time Complexity:** O(N) where N is number of keys
- **Returns:** OK
- **Example:**
  ```
  > MSET key1 "value1" key2 "value2" key3 "value3"
  OK
  ```

### INCR
Increment the integer value of a key by one.
```
INCR key
```
- **Time Complexity:** O(1)
- **Returns:** Integer value after increment
- **Example:**
  ```
  > SET counter 10
  OK
  > INCR counter
  11
  > INCR newcounter
  1
  ```

### INCRBY
Increment the integer value of a key by given amount.
```
INCRBY key increment
```
- **Time Complexity:** O(1)
- **Returns:** Integer value after increment
- **Example:**
  ```
  > INCRBY counter 5
  16
  > INCRBY counter -3
  13
  ```

### DECR
Decrement the integer value of a key by one.
```
DECR key
```
- **Time Complexity:** O(1)
- **Returns:** Integer value after decrement
- **Example:**
  ```
  > DECR counter
  12
  ```

### DECRBY
Decrement the integer value of a key by given amount.
```
DECRBY key decrement
```
- **Time Complexity:** O(1)
- **Returns:** Integer value after decrement
- **Example:**
  ```
  > DECRBY counter 3
  9
  ```

### APPEND
Append a value to a key.
```
APPEND key value
```
- **Time Complexity:** O(1)
- **Returns:** Length of string after append
- **Example:**
  ```
  > SET greeting "Hello"
  OK
  > APPEND greeting " World"
  11
  > GET greeting
  "Hello World"
  ```

### STRLEN
Get the length of the value stored in a key.
```
STRLEN key
```
- **Time Complexity:** O(1)
- **Returns:** Length of string
- **Example:**
  ```
  > STRLEN greeting
  11
  ```

### SETEX
Set key with expiration in seconds.
```
SETEX key seconds value
```
- **Time Complexity:** O(1)
- **Returns:** OK
- **Example:**
  ```
  > SETEX tempkey 60 "temporary"
  OK
  ```

### SETNX
Set key only if it doesn't exist.
```
SETNX key value
```
- **Time Complexity:** O(1)
- **Returns:** 1 if set, 0 if not
- **Example:**
  ```
  > SETNX newkey "value"
  1
  > SETNX newkey "othervalue"
  0
  ```

### GETSET
Set key and return old value atomically.
```
GETSET key value
```
- **Time Complexity:** O(1)
- **Returns:** Old value or nil
- **Example:**
  ```
  > SET mykey "oldvalue"
  OK
  > GETSET mykey "newvalue"
  "oldvalue"
  ```

### GETRANGE
Get substring of the string stored at key.
```
GETRANGE key start end
```
- **Time Complexity:** O(N) where N is length of returned string
- **Returns:** Substring
- **Example:**
  ```
  > SET mykey "Hello World"
  OK
  > GETRANGE mykey 0 4
  "Hello"
  > GETRANGE mykey -5 -1
  "World"
  ```

### SETRANGE
Overwrite part of string at key starting at offset.
```
SETRANGE key offset value
```
- **Time Complexity:** O(1)
- **Returns:** Length of string after modification
- **Example:**
  ```
  > SET mykey "Hello World"
  OK
  > SETRANGE mykey 6 "Redis"
  11
  > GET mykey
  "Hello Redis"
  ```

## List Commands

### LPUSH
Insert values at the head of the list.
```
LPUSH key value [value ...]
```
- **Time Complexity:** O(1) for each value
- **Returns:** Length of list after push
- **Example:**
  ```
  > LPUSH mylist "world"
  1
  > LPUSH mylist "hello"
  2
  ```

### RPUSH
Insert values at the tail of the list.
```
RPUSH key value [value ...]
```
- **Time Complexity:** O(1) for each value
- **Returns:** Length of list after push
- **Example:**
  ```
  > RPUSH mylist "first"
  3
  > RPUSH mylist "second"
  4
  ```

### LPOP
Remove and return element from head of list.
```
LPOP key [count]
```
- **Time Complexity:** O(N) where N is count
- **Returns:** Popped element(s) or nil
- **Example:**
  ```
  > LPOP mylist
  "hello"
  > LPOP mylist 2
  1) "world"
  2) "first"
  ```

### RPOP
Remove and return element from tail of list.
```
RPOP key [count]
```
- **Time Complexity:** O(N) where N is count
- **Returns:** Popped element(s) or nil
- **Example:**
  ```
  > RPOP mylist
  "second"
  ```

### LLEN
Get the length of a list.
```
LLEN key
```
- **Time Complexity:** O(1)
- **Returns:** Length of list
- **Example:**
  ```
  > LLEN mylist
  0
  ```

### LRANGE
Get a range of elements from list.
```
LRANGE key start stop
```
- **Time Complexity:** O(S+N) where S is start offset and N is number of elements
- **Returns:** List of elements
- **Example:**
  ```
  > RPUSH mylist a b c d e
  5
  > LRANGE mylist 0 2
  1) "a"
  2) "b"
  3) "c"
  > LRANGE mylist -3 -1
  1) "c"
  2) "d"
  3) "e"
  ```

### LINDEX
Get element at index in list.
```
LINDEX key index
```
- **Time Complexity:** O(N) where N is index
- **Returns:** Element at index or nil
- **Example:**
  ```
  > LINDEX mylist 0
  "a"
  > LINDEX mylist -1
  "e"
  ```

### LSET
Set element at index in list.
```
LSET key index value
```
- **Time Complexity:** O(N) where N is index
- **Returns:** OK or error
- **Example:**
  ```
  > LSET mylist 0 "A"
  OK
  ```

### LREM
Remove elements from list.
```
LREM key count value
```
- **Time Complexity:** O(N+M) where N is list length and M is removed elements
- **Parameters:**
  - `count > 0`: Remove from head to tail
  - `count < 0`: Remove from tail to head
  - `count = 0`: Remove all occurrences
- **Returns:** Number of removed elements
- **Example:**
  ```
  > RPUSH mylist a b a c a
  5
  > LREM mylist 2 a
  2
  ```

### LTRIM
Trim list to specified range.
```
LTRIM key start stop
```
- **Time Complexity:** O(N) where N is number of removed elements
- **Returns:** OK
- **Example:**
  ```
  > LTRIM mylist 1 3
  OK
  ```

### BLPOP
Blocking left pop from lists.
```
BLPOP key [key ...] timeout
```
- **Time Complexity:** O(1)
- **Returns:** Popped element or nil after timeout
- **Example:**
  ```
  > BLPOP list1 list2 10
  1) "list1"
  2) "value"
  ```

### BRPOP
Blocking right pop from lists.
```
BRPOP key [key ...] timeout
```
- **Time Complexity:** O(1)
- **Returns:** Popped element or nil after timeout
- **Example:**
  ```
  > BRPOP list1 list2 10
  1) "list2"
  2) "value"
  ```

## Hash Commands

### HSET
Set field in hash.
```
HSET key field value [field value ...]
```
- **Time Complexity:** O(1) for each field
- **Returns:** Number of fields added
- **Example:**
  ```
  > HSET user:1 name "John" age 30
  2
  ```

### HGET
Get value of field in hash.
```
HGET key field
```
- **Time Complexity:** O(1)
- **Returns:** Value or nil
- **Example:**
  ```
  > HGET user:1 name
  "John"
  ```

### HMSET
Set multiple fields in hash (deprecated, use HSET).
```
HMSET key field value [field value ...]
```
- **Time Complexity:** O(N) where N is number of fields
- **Returns:** OK
- **Example:**
  ```
  > HMSET user:2 name "Jane" age 25
  OK
  ```

### HMGET
Get values of multiple fields in hash.
```
HMGET key field [field ...]
```
- **Time Complexity:** O(N) where N is number of fields
- **Returns:** List of values
- **Example:**
  ```
  > HMGET user:1 name age city
  1) "John"
  2) "30"
  3) (nil)
  ```

### HGETALL
Get all fields and values in hash.
```
HGETALL key
```
- **Time Complexity:** O(N) where N is hash size
- **Returns:** List of field-value pairs
- **Example:**
  ```
  > HGETALL user:1
  1) "name"
  2) "John"
  3) "age"
  4) "30"
  ```

### HDEL
Delete fields from hash.
```
HDEL key field [field ...]
```
- **Time Complexity:** O(N) where N is number of fields
- **Returns:** Number of deleted fields
- **Example:**
  ```
  > HDEL user:1 age
  1
  ```

### HLEN
Get number of fields in hash.
```
HLEN key
```
- **Time Complexity:** O(1)
- **Returns:** Number of fields
- **Example:**
  ```
  > HLEN user:1
  1
  ```

### HEXISTS
Check if field exists in hash.
```
HEXISTS key field
```
- **Time Complexity:** O(1)
- **Returns:** 1 if exists, 0 if not
- **Example:**
  ```
  > HEXISTS user:1 name
  1
  > HEXISTS user:1 age
  0
  ```

### HKEYS
Get all field names in hash.
```
HKEYS key
```
- **Time Complexity:** O(N) where N is hash size
- **Returns:** List of field names
- **Example:**
  ```
  > HKEYS user:1
  1) "name"
  ```

### HVALS
Get all values in hash.
```
HVALS key
```
- **Time Complexity:** O(N) where N is hash size
- **Returns:** List of values
- **Example:**
  ```
  > HVALS user:1
  1) "John"
  ```

### HINCRBY
Increment integer value of hash field.
```
HINCRBY key field increment
```
- **Time Complexity:** O(1)
- **Returns:** Value after increment
- **Example:**
  ```
  > HSET user:1 points 100
  1
  > HINCRBY user:1 points 20
  120
  ```

### HINCRBYFLOAT
Increment float value of hash field.
```
HINCRBYFLOAT key field increment
```
- **Time Complexity:** O(1)
- **Returns:** Value after increment
- **Example:**
  ```
  > HSET product:1 price 9.99
  1
  > HINCRBYFLOAT product:1 price 0.01
  "10"
  ```

## Set Commands

### SADD
Add members to set.
```
SADD key member [member ...]
```
- **Time Complexity:** O(1) for each member
- **Returns:** Number of added members
- **Example:**
  ```
  > SADD myset "hello" "world"
  2
  > SADD myset "hello"
  0
  ```

### SREM
Remove members from set.
```
SREM key member [member ...]
```
- **Time Complexity:** O(N) where N is number of members
- **Returns:** Number of removed members
- **Example:**
  ```
  > SREM myset "hello"
  1
  ```

### SMEMBERS
Get all members of set.
```
SMEMBERS key
```
- **Time Complexity:** O(N) where N is set cardinality
- **Returns:** List of members
- **Example:**
  ```
  > SMEMBERS myset
  1) "world"
  ```

### SISMEMBER
Check if value is member of set.
```
SISMEMBER key member
```
- **Time Complexity:** O(1)
- **Returns:** 1 if member, 0 if not
- **Example:**
  ```
  > SISMEMBER myset "world"
  1
  ```

### SCARD
Get cardinality (size) of set.
```
SCARD key
```
- **Time Complexity:** O(1)
- **Returns:** Set cardinality
- **Example:**
  ```
  > SCARD myset
  1
  ```

### SPOP
Remove and return random members from set.
```
SPOP key [count]
```
- **Time Complexity:** O(1)
- **Returns:** Popped member(s)
- **Example:**
  ```
  > SPOP myset
  "world"
  ```

### SRANDMEMBER
Get random members from set.
```
SRANDMEMBER key [count]
```
- **Time Complexity:** O(1) or O(N)
- **Returns:** Random member(s)
- **Example:**
  ```
  > SRANDMEMBER myset 2
  1) "member1"
  2) "member2"
  ```

### SUNION
Get union of sets.
```
SUNION key [key ...]
```
- **Time Complexity:** O(N) where N is total elements
- **Returns:** List of union members
- **Example:**
  ```
  > SUNION set1 set2
  1) "a"
  2) "b"
  3) "c"
  ```

### SINTER
Get intersection of sets.
```
SINTER key [key ...]
```
- **Time Complexity:** O(N*M) worst case
- **Returns:** List of intersection members
- **Example:**
  ```
  > SINTER set1 set2
  1) "b"
  ```

### SDIFF
Get difference of sets.
```
SDIFF key [key ...]
```
- **Time Complexity:** O(N) where N is total elements
- **Returns:** List of difference members
- **Example:**
  ```
  > SDIFF set1 set2
  1) "a"
  ```

## Sorted Set Commands

### ZADD
Add members to sorted set.
```
ZADD key [NX|XX] [CH] [INCR] score member [score member ...]
```
- **Time Complexity:** O(log(N)) for each member
- **Options:**
  - `NX`: Only add new members
  - `XX`: Only update existing members
  - `CH`: Return number of changed elements
  - `INCR`: Increment score
- **Returns:** Number of added members
- **Example:**
  ```
  > ZADD leaderboard 100 "player1" 200 "player2"
  2
  ```

### ZREM
Remove members from sorted set.
```
ZREM key member [member ...]
```
- **Time Complexity:** O(M*log(N))
- **Returns:** Number of removed members
- **Example:**
  ```
  > ZREM leaderboard "player1"
  1
  ```

### ZSCORE
Get score of member in sorted set.
```
ZSCORE key member
```
- **Time Complexity:** O(1)
- **Returns:** Score or nil
- **Example:**
  ```
  > ZSCORE leaderboard "player2"
  "200"
  ```

### ZRANK
Get rank of member in sorted set (low to high).
```
ZRANK key member
```
- **Time Complexity:** O(log(N))
- **Returns:** Rank (0-based) or nil
- **Example:**
  ```
  > ZRANK leaderboard "player2"
  0
  ```

### ZREVRANK
Get reverse rank of member in sorted set (high to low).
```
ZREVRANK key member
```
- **Time Complexity:** O(log(N))
- **Returns:** Rank (0-based) or nil
- **Example:**
  ```
  > ZREVRANK leaderboard "player2"
  0
  ```

### ZRANGE
Get range of members by rank.
```
ZRANGE key start stop [WITHSCORES]
```
- **Time Complexity:** O(log(N)+M)
- **Returns:** List of members
- **Example:**
  ```
  > ZRANGE leaderboard 0 -1 WITHSCORES
  1) "player2"
  2) "200"
  ```

### ZREVRANGE
Get range of members by rank (reverse).
```
ZREVRANGE key start stop [WITHSCORES]
```
- **Time Complexity:** O(log(N)+M)
- **Returns:** List of members
- **Example:**
  ```
  > ZREVRANGE leaderboard 0 -1
  1) "player2"
  ```

### ZRANGEBYSCORE
Get range of members by score.
```
ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]
```
- **Time Complexity:** O(log(N)+M)
- **Returns:** List of members
- **Example:**
  ```
  > ZRANGEBYSCORE leaderboard 100 300
  1) "player2"
  ```

### ZCARD
Get cardinality of sorted set.
```
ZCARD key
```
- **Time Complexity:** O(1)
- **Returns:** Number of members
- **Example:**
  ```
  > ZCARD leaderboard
  1
  ```

### ZCOUNT
Count members in score range.
```
ZCOUNT key min max
```
- **Time Complexity:** O(log(N))
- **Returns:** Number of members
- **Example:**
  ```
  > ZCOUNT leaderboard 100 300
  1
  ```

## Pub/Sub Commands

### SUBSCRIBE
Subscribe to channels.
```
SUBSCRIBE channel [channel ...]
```
- **Time Complexity:** O(N) where N is number of channels
- **Returns:** Subscription confirmation
- **Example:**
  ```
  > SUBSCRIBE news sports
  1) "subscribe"
  2) "news"
  3) (integer) 1
  ```

### UNSUBSCRIBE
Unsubscribe from channels.
```
UNSUBSCRIBE [channel [channel ...]]
```
- **Time Complexity:** O(N) where N is number of channels
- **Returns:** Unsubscription confirmation
- **Example:**
  ```
  > UNSUBSCRIBE news
  1) "unsubscribe"
  2) "news"
  3) (integer) 0
  ```

### PUBLISH
Publish message to channel.
```
PUBLISH channel message
```
- **Time Complexity:** O(N+M) where N is clients and M is patterns
- **Returns:** Number of clients that received message
- **Example:**
  ```
  > PUBLISH news "Breaking news!"
  2
  ```

### PSUBSCRIBE
Subscribe to channels matching patterns.
```
PSUBSCRIBE pattern [pattern ...]
```
- **Time Complexity:** O(N) where N is number of patterns
- **Returns:** Subscription confirmation
- **Example:**
  ```
  > PSUBSCRIBE news.* sports.*
  1) "psubscribe"
  2) "news.*"
  3) (integer) 1
  ```

### PUNSUBSCRIBE
Unsubscribe from patterns.
```
PUNSUBSCRIBE [pattern [pattern ...]]
```
- **Time Complexity:** O(N+M)
- **Returns:** Unsubscription confirmation
- **Example:**
  ```
  > PUNSUBSCRIBE news.*
  1) "punsubscribe"
  2) "news.*"
  3) (integer) 0
  ```

### PUBSUB
Inspect pub/sub subsystem.
```
PUBSUB subcommand [argument [argument ...]]
```
- **Subcommands:**
  - `CHANNELS [pattern]`: List active channels
  - `NUMSUB [channel ...]`: Get subscriber count
  - `NUMPAT`: Get pattern subscriber count
- **Example:**
  ```
  > PUBSUB CHANNELS
  1) "news"
  2) "sports"
  ```

## Transaction Commands

### MULTI
Start a transaction.
```
MULTI
```
- **Time Complexity:** O(1)
- **Returns:** OK
- **Example:**
  ```
  > MULTI
  OK
  ```

### EXEC
Execute transaction.
```
EXEC
```
- **Time Complexity:** O(N) where N is number of commands
- **Returns:** Array of results
- **Example:**
  ```
  > MULTI
  OK
  > SET key1 "value1"
  QUEUED
  > SET key2 "value2"
  QUEUED
  > EXEC
  1) OK
  2) OK
  ```

### DISCARD
Discard transaction.
```
DISCARD
```
- **Time Complexity:** O(N) where N is number of queued commands
- **Returns:** OK
- **Example:**
  ```
  > MULTI
  OK
  > SET key1 "value1"
  QUEUED
  > DISCARD
  OK
  ```

### WATCH
Watch keys for modifications.
```
WATCH key [key ...]
```
- **Time Complexity:** O(1) for each key
- **Returns:** OK
- **Example:**
  ```
  > WATCH mykey
  OK
  > MULTI
  OK
  > SET mykey "newvalue"
  QUEUED
  > EXEC
  (nil)  # If mykey was modified
  ```

### UNWATCH
Unwatch all keys.
```
UNWATCH
```
- **Time Complexity:** O(1)
- **Returns:** OK
- **Example:**
  ```
  > UNWATCH
  OK
  ```

## Server Commands

### INFO
Get server information and statistics.
```
INFO [section]
```
- **Time Complexity:** O(1)
- **Sections:**
  - `server`: General server info
  - `clients`: Client connections
  - `memory`: Memory usage
  - `persistence`: RDB/AOF info
  - `stats`: General statistics
  - `replication`: Master/slave info
  - `cpu`: CPU usage
  - `keyspace`: Database statistics
  - `all`: All sections (default)
- **Example:**
  ```
  > INFO server
  # Server
  redis_version:7.0.0
  redis_mode:standalone
  process_id:1234
  tcp_port:6379
  uptime_in_seconds:3600
  ```

### CONFIG GET
Get configuration parameters.
```
CONFIG GET parameter
```
- **Time Complexity:** O(N)
- **Returns:** Parameter and value
- **Example:**
  ```
  > CONFIG GET maxmemory
  1) "maxmemory"
  2) "0"
  ```

### CONFIG SET
Set configuration parameters.
```
CONFIG SET parameter value
```
- **Time Complexity:** O(1)
- **Returns:** OK or error
- **Example:**
  ```
  > CONFIG SET maxmemory 1gb
  OK
  ```

### CONFIG RESETSTAT
Reset statistics.
```
CONFIG RESETSTAT
```
- **Time Complexity:** O(1)
- **Returns:** OK
- **Example:**
  ```
  > CONFIG RESETSTAT
  OK
  ```

### CONFIG REWRITE
Rewrite configuration file.
```
CONFIG REWRITE
```
- **Time Complexity:** O(N)
- **Returns:** OK or error
- **Example:**
  ```
  > CONFIG REWRITE
  OK
  ```

### DBSIZE
Get number of keys in current database.
```
DBSIZE
```
- **Time Complexity:** O(1)
- **Returns:** Number of keys
- **Example:**
  ```
  > DBSIZE
  42
  ```

### FLUSHDB
Remove all keys from current database.
```
FLUSHDB [ASYNC]
```
- **Time Complexity:** O(N) where N is number of keys
- **Returns:** OK
- **Example:**
  ```
  > FLUSHDB
  OK
  ```

### FLUSHALL
Remove all keys from all databases.
```
FLUSHALL [ASYNC]
```
- **Time Complexity:** O(N) where N is total keys
- **Returns:** OK
- **Example:**
  ```
  > FLUSHALL
  OK
  ```

### TIME
Get current server time.
```
TIME
```
- **Time Complexity:** O(1)
- **Returns:** Unix timestamp and microseconds
- **Example:**
  ```
  > TIME
  1) "1609459200"
  2) "123456"
  ```

### MONITOR
Monitor all commands in real time.
```
MONITOR
```
- **Time Complexity:** O(N) where N is number of commands
- **Returns:** Stream of commands
- **Example:**
  ```
  > MONITOR
  OK
  1609459200.123456 [0 127.0.0.1:12345] "SET" "key" "value"
  ```

## Persistence Commands

### SAVE
Synchronously save dataset to disk.
```
SAVE
```
- **Time Complexity:** O(N) where N is number of keys
- **Returns:** OK
- **Example:**
  ```
  > SAVE
  OK
  ```

### BGSAVE
Asynchronously save dataset to disk.
```
BGSAVE [SCHEDULE]
```
- **Time Complexity:** O(1) to start, O(N) in background
- **Returns:** OK or scheduled message
- **Example:**
  ```
  > BGSAVE
  Background saving started
  ```

### LASTSAVE
Get timestamp of last successful save.
```
LASTSAVE
```
- **Time Complexity:** O(1)
- **Returns:** Unix timestamp
- **Example:**
  ```
  > LASTSAVE
  1609459200
  ```

### BGREWRITEAOF
Asynchronously rewrite AOF file.
```
BGREWRITEAOF
```
- **Time Complexity:** O(1) to start, O(N) in background
- **Returns:** OK
- **Example:**
  ```
  > BGREWRITEAOF
  Background append only file rewriting started
  ```

## Replication Commands

### SLAVEOF / REPLICAOF
Make server a replica of another instance.
```
SLAVEOF host port
REPLICAOF host port
```
- **Time Complexity:** O(1)
- **Returns:** OK
- **Example:**
  ```
  > SLAVEOF 192.168.1.100 6379
  OK
  > REPLICAOF NO ONE
  OK
  ```

### ROLE
Get replication role.
```
ROLE
```
- **Time Complexity:** O(1)
- **Returns:** Role information
- **Example:**
  ```
  > ROLE
  1) "master"
  2) (integer) 3129659
  3) 1) 1) "192.168.1.101"
        2) "6379"
        3) "3129242"
  ```

## Cluster Commands

### CLUSTER INFO
Get cluster state.
```
CLUSTER INFO
```
- **Time Complexity:** O(1)
- **Returns:** Cluster information
- **Example:**
  ```
  > CLUSTER INFO
  cluster_state:ok
  cluster_slots_assigned:16384
  cluster_slots_ok:16384
  cluster_known_nodes:6
  cluster_size:3
  ```

### CLUSTER NODES
Get cluster nodes information.
```
CLUSTER NODES
```
- **Time Complexity:** O(N) where N is number of nodes
- **Returns:** Node information
- **Example:**
  ```
  > CLUSTER NODES
  07c4...b1d master - 0 1609459200 1 connected 0-5460
  e7d...35e master - 0 1609459201 2 connected 5461-10922
  ```

### CLUSTER MEET
Add node to cluster.
```
CLUSTER MEET ip port
```
- **Time Complexity:** O(1)
- **Returns:** OK
- **Example:**
  ```
  > CLUSTER MEET 192.168.1.102 6379
  OK
  ```

### CLUSTER ADDSLOTS
Assign slots to node.
```
CLUSTER ADDSLOTS slot [slot ...]
```
- **Time Complexity:** O(N) where N is number of slots
- **Returns:** OK
- **Example:**
  ```
  > CLUSTER ADDSLOTS 0 1 2 3 4
  OK
  ```

### CLUSTER FAILOVER
Force failover.
```
CLUSTER FAILOVER [FORCE|TAKEOVER]
```
- **Time Complexity:** O(1)
- **Returns:** OK
- **Example:**
  ```
  > CLUSTER FAILOVER
  OK
  ```

## Key Management Commands

### DEL
Delete keys.
```
DEL key [key ...]
```
- **Time Complexity:** O(N) where N is number of keys
- **Returns:** Number of deleted keys
- **Example:**
  ```
  > DEL key1 key2 key3
  2
  ```

### EXISTS
Check if keys exist.
```
EXISTS key [key ...]
```
- **Time Complexity:** O(1)
- **Returns:** Number of existing keys
- **Example:**
  ```
  > EXISTS key1 key2
  1
  ```

### EXPIRE
Set key expiration in seconds.
```
EXPIRE key seconds
```
- **Time Complexity:** O(1)
- **Returns:** 1 if set, 0 if key doesn't exist
- **Example:**
  ```
  > EXPIRE mykey 60
  1
  ```

### EXPIREAT
Set key expiration as Unix timestamp.
```
EXPIREAT key timestamp
```
- **Time Complexity:** O(1)
- **Returns:** 1 if set, 0 if key doesn't exist
- **Example:**
  ```
  > EXPIREAT mykey 1609459200
  1
  ```

### TTL
Get time to live in seconds.
```
TTL key
```
- **Time Complexity:** O(1)
- **Returns:** TTL in seconds, -1 if no expire, -2 if doesn't exist
- **Example:**
  ```
  > TTL mykey
  58
  ```

### PTTL
Get time to live in milliseconds.
```
PTTL key
```
- **Time Complexity:** O(1)
- **Returns:** TTL in milliseconds
- **Example:**
  ```
  > PTTL mykey
  57842
  ```

### PERSIST
Remove expiration from key.
```
PERSIST key
```
- **Time Complexity:** O(1)
- **Returns:** 1 if removed, 0 if key doesn't exist or has no expire
- **Example:**
  ```
  > PERSIST mykey
  1
  ```

### KEYS
Find all keys matching pattern.
```
KEYS pattern
```
- **Time Complexity:** O(N) where N is number of keys
- **Returns:** List of matching keys
- **Warning:** This command should be used carefully in production
- **Example:**
  ```
  > KEYS user:*
  1) "user:1"
  2) "user:2"
  ```

### SCAN
Incrementally iterate over keys.
```
SCAN cursor [MATCH pattern] [COUNT count] [TYPE type]
```
- **Time Complexity:** O(1) for each call
- **Returns:** Next cursor and keys
- **Example:**
  ```
  > SCAN 0 MATCH user:* COUNT 10
  1) "17"
  2) 1) "user:1"
     2) "user:2"
  ```

### TYPE
Get type of key.
```
TYPE key
```
- **Time Complexity:** O(1)
- **Returns:** Type string: string, list, set, zset, hash, none
- **Example:**
  ```
  > TYPE mykey
  string
  ```

### RENAME
Rename a key.
```
RENAME key newkey
```
- **Time Complexity:** O(1)
- **Returns:** OK or error
- **Example:**
  ```
  > RENAME mykey mynewkey
  OK
  ```

### RENAMENX
Rename key only if new key doesn't exist.
```
RENAMENX key newkey
```
- **Time Complexity:** O(1)
- **Returns:** 1 if renamed, 0 if newkey exists
- **Example:**
  ```
  > RENAMENX mykey mynewkey
  1
  ```

### RANDOMKEY
Get random key.
```
RANDOMKEY
```
- **Time Complexity:** O(1)
- **Returns:** Random key or nil
- **Example:**
  ```
  > RANDOMKEY
  "user:42"
  ```

### MOVE
Move key to another database.
```
MOVE key db
```
- **Time Complexity:** O(1)
- **Returns:** 1 if moved, 0 if not
- **Example:**
  ```
  > MOVE mykey 1
  1
  ```

## Scripting Commands

### EVAL
Execute Lua script.
```
EVAL script numkeys key [key ...] arg [arg ...]
```
- **Time Complexity:** Depends on script
- **Returns:** Script result
- **Example:**
  ```
  > EVAL "return redis.call('GET', KEYS[1])" 1 mykey
  "value"
  ```

### EVALSHA
Execute Lua script by SHA1 hash.
```
EVALSHA sha1 numkeys key [key ...] arg [arg ...]
```
- **Time Complexity:** Depends on script
- **Returns:** Script result
- **Example:**
  ```
  > EVALSHA e0e1f9fabfc9d4800c877a703b823ac0578ff8db 1 mykey
  "value"
  ```

### SCRIPT LOAD
Load script into cache.
```
SCRIPT LOAD script
```
- **Time Complexity:** O(N) where N is script length
- **Returns:** SHA1 hash
- **Example:**
  ```
  > SCRIPT LOAD "return redis.call('GET', KEYS[1])"
  "e0e1f9fabfc9d4800c877a703b823ac0578ff8db"
  ```

### SCRIPT EXISTS
Check if scripts exist in cache.
```
SCRIPT EXISTS sha1 [sha1 ...]
```
- **Time Complexity:** O(N) where N is number of scripts
- **Returns:** Array of 0/1 values
- **Example:**
  ```
  > SCRIPT EXISTS e0e1f9fabfc9d4800c877a703b823ac0578ff8db
  1) (integer) 1
  ```

### SCRIPT FLUSH
Remove all scripts from cache.
```
SCRIPT FLUSH
```
- **Time Complexity:** O(N) where N is number of cached scripts
- **Returns:** OK
- **Example:**
  ```
  > SCRIPT FLUSH
  OK
  ```

### SCRIPT KILL
Kill running script.
```
SCRIPT KILL
```
- **Time Complexity:** O(1)
- **Returns:** OK or error
- **Example:**
  ```
  > SCRIPT KILL
  OK
  ```

## Notes

1. **Time Complexity**: Understanding time complexity helps in optimizing queries and avoiding performance issues.

2. **Atomic Operations**: All Redis commands are atomic, meaning they complete fully or not at all.

3. **Error Handling**: Commands return specific error messages when operations fail (e.g., wrong type, out of range).

4. **Type Safety**: Redis enforces type safety - string commands work only on strings, list commands only on lists, etc.

5. **Memory Usage**: Be aware of memory implications, especially with commands that load entire data structures.

6. **Production Considerations**:
   - Avoid `KEYS` in production, use `SCAN` instead
   - Be careful with blocking operations
   - Monitor slow commands with `SLOWLOG`
   - Use `MONITOR` sparingly as it impacts performance

For more detailed information about specific commands, visit the [Redis Documentation](https://redis.io/commands).