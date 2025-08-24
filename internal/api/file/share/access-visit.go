package share

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
)

// 访问量加一
func AccessVisit(c *gin.Context) {
	var oneID share.GetUrlDataJSONString

	if err := c.ShouldBind(&oneID); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "绑定数据错误"})
		return
	}

	fmt.Println("Received oneID:", oneID.OneId)

	// 将十六进制字符串转换为二进制（假设 one_id 是32字符的十六进制字符串）
	binID, err := uuid.Parse(oneID.OneId)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "one_id 格式错误，必须是有效的十六进制"})
		return
	}

	// 使用二进制查询并执行原子递增
	_, err = app.Engine.Exec(
		"UPDATE url_data SET visit_count = visit_count + 1 WHERE one_id = ?",
		binID[:],
	)

	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "更新访问量失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "访问量更新成功"})
}
