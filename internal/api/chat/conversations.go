package chat

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/repositories"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// ContactWithInfo xorm
type ContactWithInfo struct {
	ContactID   string `xorm:"contact_id"`
	DisplayName string `xorm:"display_name"`
	Avatar      string `xorm:"avatar"`
	LastMsg     string `xorm:"last_msg"`
	LastDate    string `xorm:"last_date"`
}

type ContactFromUser struct {
	LastMsg      string `json:"last_msg"`
	LastDate     string `json:"last_date"`
	ContactEmail string `json:"contact_email"`
}

// Contact 联系人表（维护用户-联系人关系）
type Contact struct {
	ID        int64  `xorm:"bigint pk autoincr 'id'"`         // 联系人关系ID
	UserID    string `xorm:"varchar(255) index 'user_id'"`    // 所属用户ID
	ContactID string `xorm:"varchar(255) index 'contact_id'"` // 联系人用户ID
	LastMsg   string `xorm:"text 'last_msg'"`                 // 最后一条消息内容
	LastDate  string `xorm:"varchar(50) 'last_date'"`         // 最后消息时间（如"刚刚"）
	Nickname  string `xorm:"varchar(100) 'nickname'"`         // 备注名（可选）
}

// GetConversations Conversations 这里是获得聊天列表\好友的路由\可以拿到初始消息的
func GetConversations(c *gin.Context) {
	var contacts []ContactWithInfo
	userID := cosFile.GetID(c).ID
	err := app.Engine.Table("contact").Alias("c").
		Join("INNER", "user_req u", "c.contact_id = u.id").
		Where("c.user_id = ?", userID).
		Select(`c.contact_id, 
           COALESCE(NULLIF(c.nickname, ''), u.name) AS display_name,
           u.avatar,
           c.last_msg,
           c.last_date`).
		Find(&contacts)
	fmt.Println(contacts)

	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(), "data": contacts})
		return
	}
	c.JSON(200, gin.H{
		"code": 200, "message": "success", "data": contacts,
	})
}

func SetConversation(c *gin.Context) {
	//设置你的联系人的最新消息,和你的联系人列表
	ID := cosFile.GetID(c).ID
	var content ContactFromUser

	err := c.ShouldBind(&content)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	s := repositories.NewUserRepository(app.Engine)
	table, err := s.GetUserByEmail(content.ContactEmail)

	var contacts Contact
	var userID = strconv.Itoa(ID)
	var contactID = strconv.Itoa(table.ID)

	contacts = Contact{
		UserID:    userID,
		ContactID: contactID,
		LastMsg:   content.LastMsg,
		LastDate:  content.LastDate,
		Nickname:  table.Name,
	}
	fmt.Println(contacts)

	// 批量覆盖写入

	exist, err := app.Engine.Exist(&Contact{UserID: userID, ContactID: contactID})
	if err != nil {
		return
	}
	if !exist {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"data":    contacts,
			"m":       "执行了插入操作",
		})
		_, err = app.Engine.
			Insert(&contacts)
		return
	}

	_, err = app.Engine.
		Where("user_id = ? AND contact_id = ?", userID, contactID).
		Update(&contacts)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    contacts,
		"m":       "执行了更新操作",
	})

	//设置数据的时候判断是插入还是更新
}
