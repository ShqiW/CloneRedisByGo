package commands

import (
    "fmt"
    "strings"
    
    "github.com/ShqiW/CloneRedisByGo/internal/storage"  
)

// Handler command handler
type Handler struct {
    storage storage.Storage
}

// NewHandler create a new command handler
func NewHandler(storage storage.Storage) *Handler {
    return &Handler{
        storage: storage,
    }
}

// Execute execute the command and return the RESP format response
func (h *Handler) Execute(args []string) string {
    if len(args) == 0 {
        return "-ERR no command\r\n"
    }
    
    cmd := strings.ToUpper(args[0])
    
    switch cmd {
    case "SET":
        if len(args) != 3 {
            return "-ERR wrong number of arguments for 'set' command\r\n"
        }
        return h.handleSet(args[1], args[2])
        
    case "GET":
        if len(args) != 2 {
            return "-ERR wrong number of arguments for 'get' command\r\n"
        }
        return h.handleGet(args[1])
        
    case "PING":
        return "+PONG\r\n"
        
    default:
        return fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd)
    }
}

// handleSet handle the SET command
func (h *Handler) handleSet(key, value string) string {
    err := h.storage.Set(key, []byte(value))
    if err != nil {
        return fmt.Sprintf("-ERR %s\r\n", err.Error())
    }
    return "+OK\r\n"
}

// handleGet handle the GET command
func (h *Handler) handleGet(key string) string {
    value, err := h.storage.Get(key)
    if err != nil {
        if err == storage.ErrKeyNotFound {
            return "$-1\r\n"  // nil
        }
        return fmt.Sprintf("-ERR %s\r\n", err.Error())
    }
    return fmt.Sprintf("$%d\r\n%s\r\n", len(value), string(value))
}