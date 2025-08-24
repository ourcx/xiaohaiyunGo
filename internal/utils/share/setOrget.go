package share

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"sync"
	"time"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func XData(c *gin.Context, dataJson share.UrlDataJSONString, binaryID []byte, oneID string, userId string) {
	// 开启事务
	// 获取用户ID
	var userID int
	if userId != "" {
		num, err := strconv.Atoi(userId)
		if err != nil {
			// 处理错误
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "字符串转化失败: " + err.Error(),
				"data": nil,
			})
			return
		}
		userID = num
	} else {
		userID = cosFile.GetID(c).ID
	}

	session := app.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "事务启动失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 查询关联的分享记录
	var urlShare share.UrlShare
	has, err := session.Where("user_req_id = ?", userID).Get(&urlShare)
	if err != nil {
		_ = session.Rollback()
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询分享记录失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	if !has {
		_ = session.Rollback()
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "未找到对应的分享记录",
			"data": nil,
		})
		return
	}

	// 转换ShareID（假设urlShare.ID是[]byte类型）
	shareUUID, err := uuid.FromBytes(urlShare.ID)

	if err != nil {
		_ = session.Rollback()
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "ID格式错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 准备插入数据
	urlData := share.UrlData{
		ShareID:    urlShare.ID, // 直接使用二进制格式
		Files:      dataJson.Files,
		Password:   dataJson.Password, // 实际应存储哈希值
		ExpiresAt:  dataJson.ExpiresAt,
		VisitCount: 0,
		OneId:      binaryID,
	}

	// 执行插入
	if binaryID == nil {
		oneUUID, err := uuid.Parse(oneID)
		if err != nil {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  "UUID 格式错误: " + err.Error(),
				"data": nil,
			})
			return
		}
		fmt.Println(oneUUID)
		var getUrlData share.UrlData
		_, err2 := session.Table("url_data").Where("one_id=?", oneUUID[:]).Get(&getUrlData)
		fmt.Println(getUrlData)
		if err2 != nil {
			return
		}
		// 提交事务
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
			"msg":  "获取成功",
			"data": gin.H{
				"share_id":    shareUUID.String(),
				"expires_at":  getUrlData.ExpiresAt.Format(time.RFC3339),
				"created_at":  time.Now(),
				"visit_count": getUrlData.VisitCount,
				"files":       getUrlData.Files,
				"username":    urlShare.Username,
				"id":          getUrlData.ShareID,
				"signature":   urlShare.Signature,
				"email":       urlShare.Email,
				"avatar":      urlShare.Avatar,
				"ond_id":      getUrlData.OneId,
			},
			//这是一些总的信息，还要加上一点其他的生成短链才行
		})
		return

	} else {
		var wg sync.WaitGroup
		uuidObj, err := uuid.FromBytes(binaryID)
		for _, item := range dataJson.Files {
			wg.Add(1)

			// 显式传递参数解决循环变量问题
			go func(it string, bid string) {
				defer wg.Done()
				if err := CopyOptions(it, bid); err != nil {
					fmt.Printf("复制失败 %s: %v\n", it, err)
				}
				return
			}(item, uuidObj.String()) // 传参给匿名函数
		}
		wg.Wait() // 等待所有任务完成

		//还是文件复制一份到共享桶，共享桶会把外链链接发回来给你的
		if _, err = session.Insert(&urlData); err != nil {
			_ = session.Rollback()
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "创建分享数据失败: " + err.Error(),
				"data": nil,
			})
			return
		}
	}

	// 提交事务
	if err = session.Commit(); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "事务提交失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 清除敏感信息
	urlData.Password = ""
	uuidObj, err := uuid.FromBytes(binaryID)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": gin.H{
			"share_id":    shareUUID.String(),
			"expires_at":  urlData.ExpiresAt.Format(time.RFC3339),
			"created_at":  time.Now(),
			"visit_count": urlData.VisitCount,
			"files":       urlData.Files,
			"username":    urlShare.Username,
			"id":          urlData.ShareID,
			"signature":   urlShare.Signature,
			"email":       urlShare.Email,
			"avatar":      urlShare.Avatar,
			"ond_id":      uuidObj.String(),
		},
		//这是一些总的信息，还要加上一点其他的生成短链才行
	})
}
