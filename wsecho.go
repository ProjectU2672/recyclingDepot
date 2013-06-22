package main

import (
    "fmt"
    "io"
    "strings"
    "net/http"
    "code.google.com/p/go.net/websocket"
)

type Message struct {
    Message string `json:"message"`
}
type Action struct {
    Action string `json:"action"`
}

var queue chan Message = make(chan Message)
var conns map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)

func EchoServer(ws *websocket.Conn) {
    var m Message
    // Greeting:
    if err := websocket.JSON.Send(ws, Message{Message: `<span style="font-size:2em;">♲ connected!</span>`}); err != nil {
        fmt.Println("Greeting failed: " + err.Error())
    }
    conns[ws] = true
    defer delete(conns, ws)
    for {
        // Get message from client:
        if err := websocket.JSON.Receive(ws, &m); err == io.EOF {
            break
        } else if err != nil {
            fmt.Println("Receive error: " + err.Error())
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
        queue <- m
    }
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("static")))
    http.Handle("/echo", websocket.Handler(EchoServer))

    go func() {
        for m := range queue {
            for ws := range conns {
                if err := websocket.JSON.Send(ws, m); err != nil {
                    fmt.Println("Send error: " + err.Error())
                }
            }
        }
    }();

    err := http.ListenAndServe(":12345", nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}
