package chat

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models"
	"xiaohaiyun/internal/utils/chat"
)

// 当用户在的时候搞
func sendMessageToUser(msg Message, c *gin.Context) {
	status.ToUser = msg.ToUser

	for conn, value := range clients {
		reverseClients[value] = conn // 将值作为键，原键作为值
	}

	for {
		if conn, exists := reverseClients[msg.ToUser]; exists {

			fmt.Println("进入发送消息")

			err := conn.WriteJSON(msg)

			if err != nil {

				fmt.Println("进入离线发送模式")
				//用户离线，需要存储信息
				var message = models.Message{
					ToUser:   msg.ToUser,
					Text:     msg.Text,
					Email:    msg.Email,
					Type:     "message",
					Username: msg.Username,
				}
				err := chat.SaveUnsentMessage(app.Engine, message)
				if err != nil {
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  message,
						"err":  "历史消息存储失败",
					})
					return
				}

				clientsMu.Lock()
				delete(clients, conn)
				clientsMu.Unlock()
				//下线就这样删掉
			}

			if conn, exists = reverseClients[msg.Email]; exists {
				err := conn.WriteJSON(msg)

				if err != nil {
					log.Printf("2222error: %v", err)
					clientsMu.Lock()
					delete(clients, conn)
					clientsMu.Unlock()
					//下线就这样删掉
				}
				break
			}

			break

		} else {
			fmt.Println("进入离线发送模式")
			//用户离线，需要存储信息
			var message = models.Message{
				ToUser:   msg.ToUser,
				Text:     msg.Text,
				Email:    msg.Email,
				Type:     "message",
				Username: msg.Username,
			}
			err := chat.SaveUnsentMessage(app.Engine, message)
			if err != nil {
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  message,
					"err":  "历史消息存储失败",
				})
				return
			}

			if conn, exists = reverseClients[msg.Email]; exists {
				err := conn.WriteJSON(msg)

				if err != nil {
					log.Printf("11111error: %v", err)
					clientsMu.Lock()
					delete(clients, conn)
					clientsMu.Unlock()
					//下线就这样删掉
				}
			}
			break
		}
	}

	defer wg.Done()
	//协程完成后告诉WAIT
	return
	//要通过心跳机制来检查是否下线了客户端
}

// HistoryMessage 当websocket启动的时候，如果有对应的用户在线，就去把历史消息发过去
func HistoryMessage() {

	//要把结束的方法在最前面挂起,然后才能保证外面的wg.wait不会出错
	users, err := chat.GetAllToUsers(app.Engine)
	if err != nil {
		return
	}

	for i, email := range users {
		conn, exists := reverseClients[email]

		if !exists {
			//reverseClients或者email为空的情况
			continue
		}

		if conn == nil {
			continue
		}
		//因为对数据库操作，不要挂后台上，容易消耗很多流量，有人上线就自检一下就好了
		if chat.HasMessageByCreatedAt(app.Engine, email) {
			fmt.Println("进入拿取历史消息")
			m, err2 := chat.GetMessagesByToUser(app.Engine, email)
			if err2 != nil {
				return
			}
			for _, m := range m {
				err := conn.WriteJSON(m)
				if err != nil {
					log.Printf("error: %v", err)
					return
				}
			}
			users[i] = ""
		}
	}
}
