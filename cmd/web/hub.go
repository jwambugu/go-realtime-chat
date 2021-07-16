package main

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type subscription struct {
	connection *connection
	roomID     string
}

type message struct {
	data   []byte
	roomID string
}

// hub maintains the set of active connections and broadcasts messages to the connections
type hub struct {
	// registered connections to a specific room
	rooms map[string]map[*connection]bool

	// inbound messages from the connections.
	broadcast chan message

	// register adds requests from the connections.
	register chan subscription

	// unregister removes requests from the connections.
	unregister chan subscription
}

var h = &hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (s subscription) writePump() {
	conn := s.connection
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()

		if err := conn.ws.Close(); err != nil {
			log.Println(err.Error())
			return
		}
	}()

	for {
		select {
		case msg, ok := <-conn.send:
			if !ok {
				if err := conn.write(websocket.CloseMessage, []byte{}); err != nil {
					log.Println(err.Error())
					return
				}
				return
			}

			if err := conn.write(websocket.TextMessage, msg); err != nil {
				log.Println(err.Error())
				return
			}
		case <-ticker.C:
			if err := conn.write(websocket.PingMessage, []byte{}); err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (s subscription) readPump() {
	conn := s.connection

	defer func() {
		h.unregister <- s

		if err := conn.ws.Close(); err != nil {
			log.Println(err.Error())
			return
		}
	}()

	conn.ws.SetReadLimit(maxMessageSize)

	// Wait for a message during maximum pongWait seconds.
	// If socket is still live, increase the duration by pongWait seconds.
	// So, if connection is no more connected, the for loop will break and defer function will
	// run closing the WebSocket for that connection.
	if err := conn.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err.Error())
		return
	}

	conn.ws.SetPongHandler(func(string) error {
		return conn.ws.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, msg, err := conn.ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		m := message{
			data:   msg,
			roomID: s.roomID,
		}

		h.broadcast <- m
	}
}

func (h *hub) run() {
	for {
		select {
		case sub := <-h.register:
			connections := h.rooms[sub.roomID]

			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[sub.roomID] = connections
			}

			h.rooms[sub.roomID][sub.connection] = true
		case sub := <-h.unregister:
			connections := h.rooms[sub.roomID]

			if connections != nil {
				// Check if the connection  exists, if it does, delete it
				if _, ok := connections[sub.connection]; ok {
					delete(connections, sub.connection)
					close(sub.connection.send)

					// If we have no connections in this room, delete the room
					if len(connections) == 0 {
						delete(h.rooms, sub.roomID)
					}
				}
			}
		case msg := <-h.broadcast:
			connections := h.rooms[msg.roomID]

			for c := range connections {
				select {
				case c.send <- msg.data:
				default:
					close(c.send)

					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, msg.roomID)
					}
				}
			}
		}
	}
}
