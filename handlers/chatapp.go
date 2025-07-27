package handlers

import (
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ // смена протокола на вебсокет
	CheckOrigin: func(r *http.Request) bool { return true }, // Разрешаем CORS для теста, потом жестче
}

func (s *DB) ServeWs(c *gin.Context) {
	log.Println("▶▶▶ ServeWs ENTER") // ← ПЕРВЫЙ ЛОГ
	usernameIf, _ := c.Get("username")
	username := usernameIf.(string)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ServeWs Upgrade error:", err)
		return
	}
	log.Printf("▶▶▶ ServeWs: upgraded for %s\n", username)

	client := NewClient(conn, username, make(chan []byte))
	s.H.Register <- client
	log.Println("▶▶▶ ServeWs: client registered in hub")

	go s.H.WritePump(client) // write without database
	go s.ReadPump(client)    // read with database
	log.Println("▶▶▶ ServeWs: launched pumps")
}
