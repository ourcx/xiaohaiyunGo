package chat

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
	"xiaohaiyun/internal/utils"
	webSocketLeg "xiaohaiyun/internal/utils/webSocket"
	logs "xiaohaiyun/log"
)

// WebSocket配置
var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// 全局变量（需线程安全）
var (
	clients        = make(map[*websocket.Conn]string) // 已连接的客户端
	clientsMu      sync.Mutex                         // 保护clients的互斥锁
	broadcast      = make(chan Message)               // 广播通道
	status         Response                           //对方的状态
	wg             sync.WaitGroup
	reverseClients = make(map[string]*websocket.Conn)
)

// Message 消息结构体
type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	JWT      string `json:"jwt" xorm:"jwt"`
	Email    string `json:"email"`
	Type     string `json:"type"`
	ToUser   string `json:"toUser"`
	Avatar   string `json:"avatar"`
}

//记录了谁发消息，发了什么，密钥是什么，什么类型的，发给谁的问题

type Response struct {
	ToUser string `json:"to"`
	Status int    `json:"status"`
	//对方的状态，如果是0说明对方在线，1为对方已经接受，2为对方不在线，需要特别处理
}

// HandleWebSocket Gin路由：处理WebSocket连接
func HandleWebSocket(c *gin.Context) {
	// 升级HTTP连接为WebSocket
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	//从gin哪里拿到writer和request这两个东西
	if err != nil {
		fmt.Printf("WebSocket升级失败: %v\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
		//这个地方是对前端发生意外关闭的处理
	}
	log.Printf("websocket心跳开始")
	//每次心跳开始都是用户链接,可以在这里做一个日志记录

	// 注册新客户端（加锁保证线程安全）
	clientsMu.Lock()
	clients[ws] = " "
	clientsMu.Unlock()

	defer func(ws *websocket.Conn) {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				// 判断是否为正常关闭（如客户端主动断开）
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("连接异常关闭: %v", err)
				}

				break // 退出循环，触发 defer 中的 close()
			}
		}
		log.Printf("websocket心跳停止")
		clientsMu.Lock()
		delete(clients, ws)
		clientsMu.Unlock()
		err = ws.Close()
		fmt.Println(clients)
		if err != nil {
		}
	}(ws)

	// 循环读取客户端消息
	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			// 客户端断开时清理资源
			clientsMu.Lock()
			delete(clients, ws)
			clientsMu.Unlock()
			break
		}

		//到时候在这里验证JWT并且解析JWT就不用输入名字了，而且可以检验他的凭证，防止出现过期的问题
		claims, _ := utils.ParseUserJWT(msg.JWT)
		userClaims, _ := utils.ExtractUserClaims(claims)

		clientsMu.Lock()
		clients[ws] = userClaims.Email
		msg.Email = userClaims.Email
		//msg.Username = userClaims.Name
		msg.JWT = " "
		//通过JWT知道发信人的身份，如果JWT缺失和过期，要去退出他的登录状态
		clientsMu.Unlock()
		//并发锁

		if webSocketLeg.Legislation(msg.Text) {
			err = logs.Writer("log/GroupLog.log", msg.Text+"=>存在违禁词，由用户："+msg.Username+"发出\n")
			//go是直接用根目录来定位的,和js的语法不一样,搞混了就一直用了../
			if err != nil {
				fmt.Println("错误：" + err.Error())
				return
			}
			msg.Text = "***违规词***"
		}

		// 将消息推送到广播通道
		go HistoryMessage()
		//处理历史消息（非群组的）

		if msg.Type == "message" {
			if msg.ToUser != "" {
				wg.Add(1)
				go sendMessageToUser(msg, c)
				wg.Wait()
				//私聊的函数
			} else {
				now := time.Now()
				// 生成 "年/月/日 时:分:秒" 格式（月份和日期不带前导零）
				// 使用格式模板：2006/1/2 15:04:05
				formatted := now.Format("2006/1/2 15:04:05")
				var msgHistory = MessageHistory{
					Username: msg.Username,
					ToUser:   msg.ToUser,
					Text:     msg.Text,
					JWT:      msg.JWT,
					Type:     msg.Type,
					Email:    msg.Email,
					Date:     formatted,
					Avatar:   msg.Avatar,
				}
				SetGroupHistoryByWebsocket(msgHistory)
				fmt.Println("进入群组")
				fmt.Println(msg)
				clientsMu.Lock()
				broadcast <- msg
				clientsMu.Unlock()
			}
		} else if msg.Type == "online" {
			clientsMu.Lock()
			broadcast <- Message{
				Type:     "online",
				Username: msg.Username,
				Text:     msg.Username + "用户已上线",
				Email:    msg.Email,
			}
			clientsMu.Unlock()
			//对于不同类型的数据进行对应的处理，达到一个websocket服务私聊和群聊的功能
		} else if msg.Type == "ping" {
			clientsMu.Lock()
			broadcast <- Message{
				Type:     "pong",
				Username: msg.Username,
				Email:    msg.Email,
			}
			clientsMu.Unlock()
		}
	}
}

// HandleMessages 后台任务：监听广播并发送消息
func HandleMessages() {
	for msg := range broadcast {
		// 遍历所有客户端发送消息（加锁）
		clientsMu.Lock()
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				err := client.Close()
				if err != nil {
					return
				}
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}
