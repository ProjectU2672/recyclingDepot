package main

import (
    "fmt"
    "net/http"
    "code.google.com/p/go.net/websocket"
    "strings"
//    "encoding/json"
//    "io"
//    "log"
)

type Message struct {
    Message string `json:"message"`
}

func EchoServer(ws *websocket.Conn) {
    var m Message
    for {
        if err := websocket.JSON.Receive(ws, &m); err != nil {
            fmt.Println("Receive error: " + err.Error())
            continue
        }
        m.Message = strings.ToUpper(m.Message)
        if err := websocket.JSON.Send(ws, m); err != nil {
            fmt.Println("Send error: " + err.Error())
        }
    }
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("static")))
    http.Handle("/echo", websocket.Handler(EchoServer))
    err := http.ListenAndServe(":12345", nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}
