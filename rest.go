package main

import (
	"net/http"
)

const (
	apiEndpoint string = "/api"
)
var (
	data map[string]map[int]string = make(map[string]map[int]string)
	autoinc map[string]int = make(map[string]int)
)

type CollectionHandler struct {
	Model string
}

type ItemHandler struct {
	Model string
}

func (h *CollectionHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static-rest")))
	//http.Handle(apiEndpoint)
	data["todo"] = make(map[int]string)
	autoinc["todo"] = 0
	
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
