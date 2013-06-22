package main

// FIXME race conditions!
import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Message struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Count   uint   `json:"count"`
}

var (
	queue  chan Message             = make(chan Message)
	conns  map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)
	status Message                  = Message{Action: "status"}
)

func EchoServer(ws *websocket.Conn) {
	m := Message{Action: "message"}
	// Greeting:
	m.Message = `<span style="font-size:2em;">♲ connected!</span>`
	if err := websocket.JSON.Send(ws, m); err != nil {
		fmt.Println("Greeting failed: " + err.Error())
	}

	conns[ws] = true
	status.Count++
	queue <- status
	defer func() {
		delete(conns, ws)
		ws.Close()
		status.Count--
		queue <- status
	}()

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
	}()

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
