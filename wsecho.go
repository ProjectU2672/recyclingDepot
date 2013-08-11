package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Message struct {
	Action  string `json:"action"`
	Message string `json:"message,omitempty"`
	Count   uint   `json:"count,omitempty"`
}

var (
	queue      chan Message             = make(chan Message)
	register   chan *websocket.Conn     = make(chan *websocket.Conn)
	unregister chan *websocket.Conn     = make(chan *websocket.Conn)
	conns      map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)
	status     Message                  = Message{Action: "status"}
)

func EchoServer(ws *websocket.Conn) {
	m := Message{Action: "message"}
	// Greeting:
	m.Message = `<span style="font-size:2em;">♲ connected!</span>`
	if err := websocket.JSON.Send(ws, m); err != nil {
		fmt.Println("Greeting failed: " + err.Error())
	}

	register <- ws

	for {
		// Get message from client:
		if err := websocket.JSON.Receive(ws, &m); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Receive error: " + err.Error())
		}

		// rape it:
		runes := []rune(m.Message)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		m.Message = strings.ToUpper(string(runes))

		// Send it back:
		queue <- m
	}
	unregister <- ws
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/echo", websocket.Handler(EchoServer))

	go func() {
		for {
			select {
			case m := <-queue:
				for ws := range conns {
					if err := websocket.JSON.Send(ws, m); err != nil {
						fmt.Println("Send error: " + err.Error())
					}
				}
			case c := <-register:
				conns[c] = true
				status.Count++
				go func() {
					queue <- status
				}()
			case c := <-unregister:
				delete(conns, c)
				status.Count--
				go func() {
					c.Close()
					queue <- status
				}()
			}
		}

	}()

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
