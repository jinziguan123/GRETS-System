package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写入等待时间
	writeWait = 10 * time.Second

	// Pong等待时间
	pongWait = 60 * time.Second

	// Ping发送间隔
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 32 * 1024 // 增加到32KB
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许跨域连接
		return true
	},
}

// readPump 从WebSocket连接读取消息并发送到Hub
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		// 解析消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("消息解析错误: %v", err)
			continue
		}

		// 处理不同类型的消息
		switch msg.Type {
		case "ping":
			// 心跳包回应
			pongMsg := Message{
				Type:     "pong",
				RoomUUID: c.RoomUUID,
				Data:     map[string]string{"message": "pong"},
			}
			if data, err := json.Marshal(pongMsg); err == nil {
				select {
				case c.Send <- data:
				default:
					return
				}
			}
		case "message":
			// 聊天消息，这里可以加入业务逻辑处理
			log.Printf("收到来自房间 %s 用户 %s 的消息: %v", c.RoomUUID, c.UserID, msg.Data)

			// 广播消息到房间内的其他客户端
			c.Hub.BroadcastToRoom(c.RoomUUID, message)
		}
	}
}

// writePump 从Hub接收消息并写入WebSocket连接
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将队列中的聊天消息添加到当前WebSocket消息中
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWS 处理WebSocket连接请求
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, roomUUID, userID, userOrg string) {
	log.Printf("尝试建立WebSocket连接: roomUUID=%s, userID=%s, userOrg=%s", roomUUID, userID, userOrg)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	log.Printf("WebSocket连接升级成功: roomUUID=%s, userID=%s", roomUUID, userID)

	client := &Client{
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		RoomUUID: roomUUID,
		UserID:   userID,
		UserOrg:  userOrg,
	}

	client.Hub.register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}
