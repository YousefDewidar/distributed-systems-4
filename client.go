package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net/rpc"
    "os"
    "strings"
    "time"
)

// Types should match the server definitions
type MessageArgs struct {
    Sender string
    Text   string
}

type HistoryReply struct {
    Messages []string
}

func dialWithRetry(addr string) (*rpc.Client, error) {
    var client *rpc.Client
    var err error
    backoff := time.Second
    for i := 0; i < 5; i++ {
        client, err = rpc.Dial("tcp", addr)
        if err == nil {
            return client, nil
        }
        log.Printf("dial error: %v; retrying in %v", err, backoff)
        time.Sleep(backoff)
        backoff *= 2
    }
    return nil, err
}

func printHistory(h HistoryReply) {
    fmt.Println("--- Chat history ---")
    for _, m := range h.Messages {
        fmt.Println(m)
    }
    fmt.Println("--------------------")
}

func main() {
    addr := flag.String("addr", "127.0.0.1:1234", "server address")
    name := flag.String("name", "anon", "your display name")
    flag.Parse()

    client, err := dialWithRetry(*addr)
    if err != nil {
        log.Printf("initial dial failed: %v", err)
        // continue; we'll attempt reconnect on each operation
    }

    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("Connected to %s as %s. Type messages and press Enter. Type 'history' to fetch history, 'exit' to quit.\n", *addr, *name)
    for {
        fmt.Print("> ")
        line, err := reader.ReadString('\n')
        if err != nil {
            log.Printf("read error: %v", err)
            break
        }
        text := strings.TrimSpace(line)
        if text == "exit" {
            fmt.Println("bye")
            break
        }
        // fetch history without sending a message
        if text == "history" {
            // ensure connection
            if client == nil {
                client, err = dialWithRetry(*addr)
                if err != nil {
                    log.Printf("cannot connect to server: %v", err)
                    continue
                }
            }
            var h HistoryReply
            callErr := client.Call("ChatServer.History", struct{}{}, &h)
            if callErr != nil {
                log.Printf("history call error: %v", callErr)
                client.Close()
                client = nil
                continue
            }
            printHistory(h)
            continue
        }

        // send message
        if client == nil {
            client, err = dialWithRetry(*addr)
            if err != nil {
                log.Printf("cannot connect to server: %v", err)
                continue
            }
        }
        args := MessageArgs{Sender: *name, Text: text}
        var reply HistoryReply
        err = client.Call("ChatServer.Send", args, &reply)
        if err != nil {
            log.Printf("send error: %v", err)
            client.Close()
            client = nil
            continue
        }
        printHistory(reply)
    }
    if client != nil {
        client.Close()
    }
}
