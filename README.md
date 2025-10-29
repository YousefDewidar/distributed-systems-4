
This repository contains a small RPC-based chatroom (server + client) implemented in Go.

Usage (Windows PowerShell):

1. Start the server (default address 127.0.0.1:1234):

```powershell
go run server.go
```

Or build the server:

```powershell
go build -o chatserver server.go
./chatserver
```

2. Run the client (specify a name if you like):

```powershell
go run client.go --name Alice
```

Client commands:
- Type a message and press Enter to send it to the server. The server will return the full chat history after each send.
- Type `history` to fetch the chat history without sending a message.
- Type `exit` to quit the client.

Notes:
- The client reads full lines using a buffered reader, so messages with spaces are supported.
- If the server is down, the client will attempt to reconnect with exponential backoff for a few attempts; subsequent operations also attempt reconnection.
