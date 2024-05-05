package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func CheckServer() {
	fmt.Println("server!")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not start the websocket: %s", err)
	}
	defer conn.Close()

	for {
		t, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		log.Printf("\nType\tMessage\n%d\t%v\t%s", t, message, message)
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error writing message: %s", err)
		}
	}
}

func WsConnect() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}