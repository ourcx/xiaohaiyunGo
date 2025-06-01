package share

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
)

func DeleteShare(c *gin.Context) {
	var OneID share.GetUrlDataJSONString

	err := c.ShouldBind(&OneID)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err,
		})
		return
	}

	session := app.Engine.NewSession()

	oneUUID, err := uuid.Parse(OneID.OneId)
	i, err := session.Table("url_data").Where("one_id=?", oneUUID[:]).Delete(share.UrlData{})
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "msg": err, "data": "删除失败",
		})
		return
	} else if i == 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不存在分享链接",
		})
		return
	}

	if err = session.Commit(); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "事务提交失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "链接已失效",
	})

}
