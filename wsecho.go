package main

import (
    "net/http"
    "io"
    "code.google.com/p/go.net/websocket"
)

func EchoServer(ws *websocket.Conn) {
    io.Copy(ws, ws)
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("static")))
    http.Handle("/echo", websocket.Handler(EchoServer))
    err := http.ListenAndServe(":12345", nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}
