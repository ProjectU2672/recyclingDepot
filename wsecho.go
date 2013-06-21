package main

import (
    "fmt"
    "strings"
    "net/http"
    "code.google.com/p/go.net/websocket"
)

type Message struct {
    Message string `json:"message"`
}

func EchoServer(ws *websocket.Conn) {
    var m Message
    for {
        // Get message from client:
        if err := websocket.JSON.Receive(ws, &m); err != nil {
            fmt.Println("Receive error: " + err.Error())
            continue
        }

        // rape it:
        old := m.Message
        m.Message = ""
        n := len(old)
        for i := 0; i < n; i++ {
            m.Message += old[(n-i-1):(n-i)]
        }
        m.Message = strings.ToUpper(m.Message)

        // Send it back:
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
