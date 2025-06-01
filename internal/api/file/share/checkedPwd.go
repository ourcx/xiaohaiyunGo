package share

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
)

func Checked(c *gin.Context) {
	var dataJson share.CheckUrlDataJSONString
	var data share.UrlDataJSONString
	err := c.ShouldBind(&dataJson)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err,
		})
		return
	}

	oneId, err := uuid.Parse(dataJson.OneId)

	_, err = app.Engine.Table("url_data").Where("one_id=?", oneId[:]).Get(&data)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err,
			"data": data,
		})
		return
	}

	if data.Password == dataJson.Password {
		c.JSON(200, gin.H{
			"code": 200,
		})
		return
	} else {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "验证码不匹配",
		})
		return
	}
}
