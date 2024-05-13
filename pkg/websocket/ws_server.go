package websocket

import (
	"fmt"
	//hs "linear-db/pkg/httpserver"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type cType struct {
	UserType string
	UserKey string
}

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

func upgradeProtocol(w http.ResponseWriter, r *http.Request)(*websocket.Conn, error){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}


func handler(w http.ResponseWriter, r *http.Request) {
	var h cType
	h.UserType = r.Header.Get("User-Type")
	h.UserKey = r.Header.Get("User-Key")
	if h.UserType != "admin" {
		fmt.Printf("[%s] tried to connect to WS\n", r.RemoteAddr)
		http.Error(w, "Not allowed", http.StatusForbidden)
		return
	}
	conn, err := upgradeProtocol(w, r)
	if err != nil {
		log.Printf("Could not start the websocket: %s", err)
	}
	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()
	// go func() {
	// 	for update := range hs.Comm {
	// 		conn.WriteMessage(websocket.TextMessage, []byte(update))
	// 	}
	// }()
	// go func() {
	// 	for {
	// 		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Connection is alive %v", time.Now())))
	// 		time.Sleep(5000 * time.Millisecond)
	// 	}
	// }()
	go func() {
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
	}()
}

func WsConnect() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
