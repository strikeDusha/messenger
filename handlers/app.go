package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	connection map[*websocket.Conn]bool
}

func NewWSS() *Server {
	return &Server{
		make(map[*websocket.Conn]bool),
	}
}

// later  - USE MUTEX!!!
// maps in go is not concurrent
func (s *Server) HandleWS(ws *websocket.Conn) {

}

func ChatApp(c *gin.Context) {
	un, _ := c.Get("username")
	c.JSON(200, gin.H{"hello": un})
}
