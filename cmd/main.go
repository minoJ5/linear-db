// use bytes fields to seperate a string
// https://pkg.go.dev/bytes#FieldsFunc

package main

import (
	hs "linear-db/pkg/httpserver"
	ws "linear-db/pkg/websocket"

	"sync"
)


func main() {
	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(2)
	go ws.WsConnect()
	go hs.HttpServer()
	wg.Wait()
}
