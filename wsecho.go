package main

import (
    "fmt"
    "net/http"
    "code.google.com/p/go.net/websocket"
    "strings"
    "encoding/json"
    "io"
    "log"
)

type Message struct {
    Message string `json:"message"`
}

func EchoServer(ws *websocket.Conn) {
    var m Message
    dec := json.NewDecoder(ws)
    for {
        if err := dec.Decode(&m); err == io.EOF {
            fmt.Println("decoder got EOF")
        } else if err != nil {
            log.Fatal(err)
            continue
        }
        go func() {
            m.Message = strings.ToUpper(m.Message)
            b, err := json.Marshal(m)
            if err != nil {
                fmt.Println("Marshal error: " + err.Error())
            }
            ws.Write(b)
        }()
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
