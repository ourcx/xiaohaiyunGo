package share

import (
	"github.com/gin-gonic/gin"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
)

// 访问量加一
func AccessVisit(c *gin.Context) {
	var oneID share.GetUrlDataJSONString

	err := c.ShouldBind(&oneID)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400, "msg": "绑定数据错误",
		})
		return
	}

	_, err = app.Engine.Table("url_data").Where("one_id=?", oneID.OneId).Update("visit_num", "visit_num+1")
	if err != nil {
		return
	}

	c.JSON(200, nil)
}
