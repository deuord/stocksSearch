package controllers

import (
	"log"
	"net/http"

	"stocksSearch/backend/services"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSController WebSocket控制器
type WSController struct {
	hub *services.Hub
}

// NewWSController 创建WebSocket控制器
func NewWSController(hub *services.Hub) *WSController {
	return &WSController{hub: hub}
}

// HandleWS 处理前端WebSocket连接
func (wc *WSController) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ WebSocket升级失败: %v", err)
		return
	}

	client := &services.Client{}
	client.Conn = conn
	client.Hub = wc.hub
	client.Send = make(chan []byte, 256)

	wc.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
