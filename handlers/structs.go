package handlers

import (
	"messenger/database"
	"messenger/ws"
)

type DB struct { // database for all logic
	s *database.Storage // storage
	H *ws.Hub           //websocket hub
}

func NewDB(s *database.Storage, h *ws.Hub) *DB {
	return &DB{s: s, H: h}
}
