package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c connection) write(messageType int, payload []byte) error {
	fmt.Printf("Writing message %s of type %d\n", string(payload), messageType)

	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}

	return c.ws.WriteJSON(map[string]interface{}{
		"messageType": messageType,
		"body":        string(payload),
	})
}

func serveWebsockets(w http.ResponseWriter, r *http.Request, roomID string) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	conn := &connection{
		ws:   ws,
		send: make(chan []byte, 256),
	}

	sub := subscription{
		connection: conn,
		roomID:     roomID,
	}

	h.register <- sub

	go sub.writePump()
	go sub.readPump()

}
