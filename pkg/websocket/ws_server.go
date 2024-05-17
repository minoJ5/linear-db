// Parts of this code are from the example chat in the gorilla websocket package.
// Link: https://github.com/gorilla/websocket/blob/main/examples/chat/client.go
// Licence is in the licenses/ws_license.txt file.

package websocket

import (
	"bytes"
	"fmt"
	hs "linear-db/pkg/httpserver"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type cType struct {
	UserType  string
	UserKey   string
	GetUpdate string
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var H *Hub

func init() {
	H = &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}
func upgradeProtocol(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
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
	h.GetUpdate = r.Header.Get("Get-Update")
	if h.UserType != "admin" {
		fmt.Printf("[%s] tried to connect to WS\n", r.RemoteAddr)
		http.Error(w, "Not allowed", http.StatusForbidden)
		return
	}
	conn, err := upgradeProtocol(w, r)
	if err != nil {
		log.Printf("Could not start the websocket: %s", err)
	}
	client := &Client{hub: H, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(3)
	defer wg.Wait()
	go client.readPump()
	go client.writePump()
	go func() {
		for update := range hs.Comm {
			if h.GetUpdate == "true" {
				message := bytes.TrimSpace(bytes.Replace([]byte(update), newline, space, -1))
				client.hub.broadcast <- message
			}
		}
	}()
	// go func() {
	// 	for {
	// 		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Connection is alive %v", time.Now())))
	// 		time.Sleep(5000 * time.Millisecond)
	// 	}
	// }()
	// go func() {
	// 	for {
	// 		t, message, err := conn.ReadMessage()
	// 		if err != nil {
	// 			log.Printf("Error reading message: %v", err)
	// 			return
	// 		}
	// 		log.Printf("\nType\tMessage\n%d\t%v\t%s", t, message, message)
	// 		err = conn.WriteMessage(websocket.TextMessage, message)
	// 		if err != nil {
	// 			log.Printf("Error writing message: %s", err)
	// 		}
	// 	}
	// }()
}

func WsConnect() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("WS server: Error loading .env file")
	}
	go H.run()
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(os.Getenv("URL_WS"), nil)
}
