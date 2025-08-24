package share

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"xiaohaiyun/internal/models/share"
	share2 "xiaohaiyun/internal/utils/share"
)

func CreateUrl(c *gin.Context) {
	// 绑定请求数据
	id := uuid.New()                  // 生成UUIDv4
	binaryID, _ := id.MarshalBinary() // 转换为16字节
	var dataJson share.UrlDataJSONString
	if err := c.ShouldBindJSON(&dataJson); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}
	share2.XData(c, dataJson, binaryID, "", "")
}
