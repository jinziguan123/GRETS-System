package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Message WebSocket消息结构
type Message struct {
	Type     string      `json:"type"`
	RoomUUID string      `json:"roomUUID"`
	Data     interface{} `json:"data"`
}

// Client WebSocket客户端
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	RoomUUID string
	UserID   string
	UserOrg  string
}

// Hub WebSocket连接集线器
type Hub struct {
	// 已注册的客户端
	clients map[*Client]bool

	// 按房间分组的客户端
	rooms map[string]map[*Client]bool

	// 客户端注册请求
	register chan *Client

	// 客户端注销请求
	unregister chan *Client

	// 来自客户端的消息
	broadcast chan []byte

	// 房间消息广播
	roomBroadcast chan *RoomMessage

	// 客户端连接锁
	mu sync.RWMutex
}

// RoomMessage 房间消息
type RoomMessage struct {
	RoomUUID string
	Message  []byte
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		rooms:         make(map[string]map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan []byte),
		roomBroadcast: make(chan *RoomMessage),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true

			// 添加到房间
			if h.rooms[client.RoomUUID] == nil {
				h.rooms[client.RoomUUID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomUUID][client] = true
			h.mu.Unlock()

			log.Printf("用户 %s 连接到房间 %s", client.UserID, client.RoomUUID)

			// 发送连接成功消息
			joinMsg := Message{
				Type:     "join",
				RoomUUID: client.RoomUUID,
				Data:     map[string]string{"message": "成功加入聊天室"},
			}
			if data, err := json.Marshal(joinMsg); err == nil {
				select {
				case client.Send <- data:
				default:
					close(client.Send)
					delete(h.clients, client)
					if h.rooms[client.RoomUUID] != nil {
						delete(h.rooms[client.RoomUUID], client)
					}
				}
			}

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if h.rooms[client.RoomUUID] != nil {
					delete(h.rooms[client.RoomUUID], client)
					if len(h.rooms[client.RoomUUID]) == 0 {
						delete(h.rooms, client.RoomUUID)
					}
				}
				close(client.Send)
				log.Printf("用户 %s 断开房间 %s 连接", client.UserID, client.RoomUUID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
					if h.rooms[client.RoomUUID] != nil {
						delete(h.rooms[client.RoomUUID], client)
					}
				}
			}
			h.mu.RUnlock()

		case roomMsg := <-h.roomBroadcast:
			h.mu.RLock()
			if roomClients, exists := h.rooms[roomMsg.RoomUUID]; exists {
				for client := range roomClients {
					select {
					case client.Send <- roomMsg.Message:
					default:
						close(client.Send)
						delete(h.clients, client)
						delete(roomClients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToRoom 向指定房间广播消息
func (h *Hub) BroadcastToRoom(roomUUID string, message []byte) {
	h.roomBroadcast <- &RoomMessage{
		RoomUUID: roomUUID,
		Message:  message,
	}
}

// GetRoomClients 获取房间内的客户端数量
func (h *Hub) GetRoomClients(roomUUID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if roomClients, exists := h.rooms[roomUUID]; exists {
		return len(roomClients)
	}
	return 0
}
