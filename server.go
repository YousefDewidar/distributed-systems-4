package main

import (
    "flag"
    "fmt"
    "log"
    "net"
    "net/rpc"
    "sync"
)

// MessageArgs represents a chat message sent by a client.
type MessageArgs struct {
    Sender string
    Text   string
}

// HistoryReply is returned to clients containing the full chat history.
type HistoryReply struct {
    Messages []string
}

// ChatServer holds the chat history and exposes RPC methods.
type ChatServer struct {
    mu   sync.Mutex
    msgs []string
}

// Send appends a message to the chat history and returns the full history.
func (c *ChatServer) Send(args MessageArgs, reply *HistoryReply) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    entry := fmt.Sprintf("%s: %s", args.Sender, args.Text)
    c.msgs = append(c.msgs, entry)
    // copy to reply
    reply.Messages = append([]string(nil), c.msgs...)
    return nil
}

// History returns the full chat history without adding a new message.
func (c *ChatServer) History(_ struct{}, reply *HistoryReply) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    reply.Messages = append([]string(nil), c.msgs...)
    return nil
}

func main() {
    addr := flag.String("addr", "127.0.0.1:1234", "server listen address")
    flag.Parse()

    server := &ChatServer{}
    if err := rpc.RegisterName("ChatServer", server); err != nil {
        log.Fatalf("rpc register: %v", err)
    }

    ln, err := net.Listen("tcp", *addr)
    if err != nil {
        log.Fatalf("listen %s: %v", *addr, err)
    }
    defer ln.Close()

    log.Printf("Chat server listening on %s", *addr)
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("accept error: %v", err)
            continue
        }
        go rpc.ServeConn(conn)
    }
}
