package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/ShqiW/CloneRedisByGo/internal/server"
)

// set command
// # Basic usage
// SET name "John"
// # Returns: OK

// # Set with 10 seconds expiration
// SET session:123 "user_data" EX 10
// # Returns: OK

// # Set with 5000 milliseconds (5 seconds) expiration
// SET temp:key "temporary_data" PX 5000
// # Returns: OK

// # Set only if key doesn't exist (useful for distributed locks)
// SET lock:order:123 "1" NX EX 30
// # Returns: OK (success) or (nil) (failure)

// # Set only if key exists (update existing)
// SET user:name "Jane" XX
// # Returns: OK (if key exists) or (nil) (if key doesn't exist)

// # Combined: Set only if not exists with 60 seconds expiration
// SET distributed:lock "process-123" NX EX 60

// set command
func (h *Handler) set(args []protocol.Value) protocol.Value {
    if len(args) < 2 {
        return protocol.ErrorValue("ERR wrong number of arguments for 'set' command")
    }
    
    key := args[0].String()
    value := args[1].Bytes()
    
    var ttl time.Duration
    
    // Parse additional options
    for i := 2; i < len(args); i++ {
        arg := strings.ToUpper(args[i].String())
        switch arg {
        case "EX": // Expire in seconds
            if i+1 >= len(args) {
                return protocol.ErrorValue("ERR syntax error")
            }
            seconds, err := strconv.Atoi(args[i+1].String())
            if err != nil {
                return protocol.ErrorValue("ERR value is not an integer")
            }
            ttl = time.Duration(seconds) * time.Second
            i++
        case "PX": // Expire in milliseconds
            if i+1 >= len(args) {
                return protocol.ErrorValue("ERR syntax error")
            }
            millis, err := strconv.Atoi(args[i+1].String())
            if err != nil {
                return protocol.ErrorValue("ERR value is not an integer")
            }
            ttl = time.Duration(millis) * time.Millisecond
            i++
        case "NX": // Set only if key does not exist
            if h.storage.Exists(key) > 0 {
                return protocol.NullBulkValue()
            }
        case "XX": // Set only if key exists
            if h.storage.Exists(key) == 0 {
                return protocol.NullBulkValue()
            }
        }
    }
    
    if err := h.storage.Set(key, value, ttl); err != nil {
        return protocol.ErrorValue("ERR " + err.Error())
    }
    
    return protocol.StringValue("OK")
}

// get command
// # Get existing key
// GET name
// # Returns: "John"

// # Get non-existent key
// GET nonexistent
// # Returns: (nil)

// # Get expired key
// GET session:123  # Assuming it's expired
// # Returns: (nil)
func (h *Handler) get(args []protocol.Value) protocol.Value {
    if len(args) != 1 {
        return protocol.ErrorValue("ERR wrong number of arguments for 'get' command")
    }
    
    key := args[0].String()
    value, err := h.storage.Get(key)
    if err != nil {
        return protocol.NullBulkValue()
    }
    
    return protocol.BulkValue(value)
}