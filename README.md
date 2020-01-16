# chat-room-backend

## description

Go语言学习中，用Go语言写的一个简单聊天室后端。
基于websocket建立服务端与客户端全双工通信。
主要利用了Go语言的`channel`及`goroutine`等特性，以及第三方websocket包 `gorilla/websocket`。

## build and run

```bash
go build gochatroom
./gochatroom --addr=:8081
```

to be improved
