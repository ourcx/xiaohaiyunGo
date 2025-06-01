package share

import (
	"github.com/gin-gonic/gin"
	"xiaohaiyun/internal/models/share"
	share2 "xiaohaiyun/internal/utils/share"
)

func GetUrl(c *gin.Context) {

	var dataJson share.UrlDataJSONString
	var OneID share.GetUrlDataJSONString
	if err := c.ShouldBindJSON(&OneID); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}
	share2.XData(c, dataJson, nil, OneID.OneId)
}
