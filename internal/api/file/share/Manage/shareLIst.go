package Manage

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"time"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// 在您的 handler 包或一个公共的 response DTO 包中
type UrlDataResponseItem struct {
	FileName   []string  `json:"fileName"`
	FileVisit  int       `json:"fileVisit"`
	CreateTime time.Time `json:"createTime"`
	OneId      string    `json:"oneId"` // 这里是格式化后的 UUID 字符串
	Password   string    `json:"password"`
	// ... 其他对应字段
}

// ShareList 拿到用户的分享链接列表
func ShareList(c *gin.Context) {
	userID := cosFile.GetID(c).ID // 假设此函数能正确获取用户ID

	session := app.Engine.NewSession()
	defer session.Close() // 良好的习惯，确保会话关闭

	var urlShare share.UrlShare   // 父分享记录
	var dbUrlData []share.UrlData // 从数据库获取的原始 urlData 列表

	// 1. 获取 urlShare 记录 (父分享记录)
	found, err := session.Table("url_share").Where("user_req_id=?", userID).Get(&urlShare)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "查询分享主记录失败: " + err.Error()})
		return
	}
	if !found {
		c.JSON(200, gin.H{ // 或者 404，根据业务需求
			"code":          200, // 或者 404
			"msg":           "未找到用户的分享主记录",
			"data":          []UrlDataResponseItem{}, // 返回空列表
			"parentShareId": "",
		})
		return
	}

	// 2. 获取关联的 urlData 列表
	// 假设 urlShare.ID 是 []byte 类型，并且是 url_data 表中 share_id 的外键
	err = session.Table("url_data").Where("share_id=?", urlShare.ID).Find(&dbUrlData)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "查询分享数据列表失败: " + err.Error()})
		return
	}

	// (关于 session.Commit() 的说明：如果前面的操作都是只读，并且没有显式开启事务 session.Begin()，
	// 那么这里的 Commit() 通常是不必要的。如果确实需要事务来保证读取一致性，
	// 则应使用 session.Begin(), 并在出错时 session.Rollback(), 成功时 session.Commit()。)
	// 如果没有其他写操作或特殊事务需求，此处可以考虑移除 session.Commit()。
	// 为了安全，如果之前的代码包含了它且您不确定，可以暂时保留，但建议审视其必要性。
	// if err = session.Commit(); err != nil {
	//    c.JSON(500, gin.H{"code": 500, "msg": "事务提交失败: " + err.Error()})
	//    return
	// }

	// 3. 格式化 urlShare.ID (父分享ID)
	parentShareIdStr := ""
	if urlShare.ID != nil && len(urlShare.ID) == 16 {
		parentUUID, err := uuid.FromBytes(urlShare.ID)
		if err == nil {
			parentShareIdStr = parentUUID.String()
		} else {
			log.Printf("警告: 转换父分享ID (urlShare.ID) 失败: %v", err)
			// 可以选择在这里返回错误，或者让 parentShareIdStr 为空
		}
	} else if urlShare.ID != nil {
		log.Printf("警告: 父分享ID (urlShare.ID) 长度不正确，期望16字节，得到 %d字节", len(urlShare.ID))
	}

	// 4. 转换 dbUrlData 列表为其响应格式 (处理每个item的OneId)
	responseDataList := make([]UrlDataResponseItem, 0, len(dbUrlData)) // 初始化以避免nil

	for i, item := range dbUrlData {
		oneIdStr := "" // 默认值
		if item.OneId != nil && len(item.OneId) == 16 {
			id, err := uuid.FromBytes(item.OneId)
			if err == nil {
				oneIdStr = id.String()
			} else {
				// 记录特定item的OneId转换错误，但可能继续处理其他item
				log.Printf("处理urlData列表项 %d 的OneId转换失败: %v (原始OneId: %x)", i, err, item.OneId)
				// oneIdStr 保持为空字符串或设置为一个错误指示符，如 "invalid_id_format"
			}
		} else if item.OneId != nil { // OneId存在但长度不正确
			log.Printf("urlData列表项 %d 的OneId长度不正确: 期望16字节, 得到 %d字节 (原始OneId: %x)", i, len(item.OneId), item.OneId)
		}
		// else: item.OneId 为 nil, oneIdStr 保持为空字符串

		// 构建响应对象，你需要根据你的 UrlData 和 UrlDataResponseItem 结构体来映射字段
		responseItem := UrlDataResponseItem{
			OneId: oneIdStr,
			// 从 item (share.UrlData 类型) 复制其他字段到 responseItem
			// 例如:
			FileName:   item.Files,
			FileVisit:  item.VisitCount,
			CreateTime: item.CreatedAt,
			Password:   item.Password,
			// ... 其他字段
		}
		responseDataList = append(responseDataList, responseItem)
	}

	// 5. 成功返回
	c.JSON(200, gin.H{
		"code":          200,
		"data":          responseDataList, // 这是处理过的列表
		"parentShareId": parentShareIdStr, // 这是 urlShare.ID 转换后的字符串
	})
}
