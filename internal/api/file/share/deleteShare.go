package share

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func DeleteShare(c *gin.Context) {
	var OneID share.GetUrlDataJSONString
	var client = cosFile.Client()
	err := c.ShouldBind(&OneID)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err,
		})
		return
	}
	fmt.Println(OneID)

	session := app.Engine.NewSession()

	oneUUID, err := uuid.Parse(OneID.OneId)
	dir := "share/" + oneUUID.String() + "/"
	//在cos处删除分享的文件
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:  dir,
		MaxKeys: 1000,
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := client.Bucket.Get(context.Background(), opt)
		if err != nil {
			// Error
			break
		}
		for _, content := range v.Contents {
			_, err = client.Object.Delete(context.Background(), content.Key)
			if err != nil {
				// Error
			}
		}
		isTruncated = v.IsTruncated
		marker = v.NextMarker
	}
	//在本地数据库删除相关的记录
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
