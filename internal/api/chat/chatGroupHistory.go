package chat

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"xiaohaiyun/internal/app"
)

type MessageHistory struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	JWT      string `json:"jwt" xorm:"jwt"`
	Email    string `json:"email"`
	Type     string `json:"type"`
	ToUser   string `json:"toUser"`
	Date     string `json:"date" xorm:"date"`
	Avatar   string `json:"avatar"`
}

// SetGroupHistoryByWebsocket SetGroupHistory 群组聊天的数据库缓存系统
func SetGroupHistoryByWebsocket(msg MessageHistory) {

	affected, err := app.Engine.Table("messagesGroupHistory").Insert(&msg)
	if err != nil {
		log.Printf("插入失败: %v", err)
		return
	}
	log.Printf("插入成功，影响行数: %d", affected)
}

// SetGroupHistory 群组聊天的数据库缓存系统
func SetGroupHistory(c *gin.Context) {
	var msg Message
	err := c.ShouldBind(&msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	affected, err := app.Engine.Table("messagesGroupHistory").Insert(&msg)
	if err != nil {
		log.Printf("插入失败: %v", err)
		return
	}
	log.Printf("插入成功，影响行数: %d", affected)
}

func GetGroupHistory(c *gin.Context) {
	var messages []MessageHistory

	// 添加错误处理
	err := app.Engine.
		Table("messagesGroupHistory").
		Find(&messages)

	if err != nil {
		log.Printf("查询消息历史失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"messages": err})
	}

	if len(messages) == 0 {
		log.Println("没有历史消息")
		c.JSON(http.StatusNoContent, gin.H{"message": "no"})
	}
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
