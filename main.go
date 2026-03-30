package main

import (
    "log"
    "lab2-chat/server"
)

func main() {
    chatServer := server.NewChatServer(":8080")
    log.Println(" Starting Chat Server Lab 2")
    if err := chatServer.Start(); err != nil {
        log.Fatalf(" Server error: %v", err)
    }
}
