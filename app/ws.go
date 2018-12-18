package app

import (
	"log"

	"github.com/gin-gonic/gin"
)

var Huber = newHub()

// Register register stats service
func WSRegister(router *gin.RouterGroup) {
	router.GET("/ws", serveWs)
}

// serveWs handles websocket requests from the peer.
func serveWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: Huber, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	// go client.readPump()
}
