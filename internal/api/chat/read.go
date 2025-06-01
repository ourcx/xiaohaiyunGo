package chat

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Read 这个文件是来标记消息是否已经被读到了
// 或者说是当新消息出现的时候的提示
// 消息的红点显示
func Read(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
