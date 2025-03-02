package v1

import (
	"log"
	"net/http"
	"sync"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/media/pkg/jwt"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域（生产环境需限制）
	},
}

// 管理所有活跃的 WebSocket 连接
type WebSocketPool struct {
	clients map[uint]*websocket.Conn
	mu      sync.Mutex
}

var pool = WebSocketPool{
	clients: make(map[uint]*websocket.Conn),
}

// 处理 WebSocket 连接
func HandleWebSocket(c *gin.Context) {

	token := c.Query("token") // 从 URL 参数获取 token
	if token == "" {
		token = c.GetHeader("Authorization")
		// 移除 "Bearer " 前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		c.Abort()
		return
	}

	// 验证 token 并获取用户 ID
	claims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token无效"})
		c.Abort()
		return
	}

	userID := claims.UserID

	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 将连接加入池中
	pool.mu.Lock()
	pool.clients[userID] = conn
	pool.mu.Unlock()

	// 监听连接关闭
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			pool.mu.Lock()
			delete(pool.clients, userID)
			pool.mu.Unlock()
			break
		}
	}
	c.Set("userID", claims.UserID)
	c.Next()
}

// 广播消息给指定用户
func SendNotificationToUser(userID uint, message []byte) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if conn, ok := pool.clients[userID]; ok {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("WebSocket send error:", err)
		}
	}
}
