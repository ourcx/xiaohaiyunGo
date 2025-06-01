package relationship

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models"
	"xiaohaiyun/internal/models/RelationShip"
	cosFile "xiaohaiyun/internal/utils/cos"
	"xiaohaiyun/internal/utils/rBook"
)

//返回联系人列表
//数据格式
// {
//        src: 'https://s2.loli.net/2025/02/02/ELbK6urJqYvgBPj.jpg',从user_req拿到
//        name: '联系人2',从user_req拿到
//        date: '刚刚',从time拿到
//        id: "3277975910@136.com"从user_req拿到
//        // 加上一个后台生成的群id码
//      },

func GetRBookList(c *gin.Context) {

	userTable := cosFile.GetID(c)
	var ship []RelationShip.Friend // 注意是切片
	err := app.Engine.Table("friends").
		Where("user_id = ?", userTable.ID).
		Find(&ship) // ← 传递切片指针
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//根据jwt查询用户列表

	var friendIDs = make([]int, len(ship))
	for i, fr := range ship {
		friendIDs[i] = fr.FriendId
	}

	results, err := rBook.ReqProfiles(friendIDs, c)
	if err != nil {
		return
	}
	//拿到raw数据，方便后续处理

	contacts := make([]RelationShip.Contact, len(results))
	for i, result := range results {
		userReq := &models.UserReq{}
		_, _ = app.Engine.Where("ID = ?", result.UserProfile.UserId).Get(userReq)
		contacts[i] = RelationShip.Contact{
			Src:  result.UserProfile.AvatarUrl,
			Name: userReq.Name,
			ID:   result.UserReqEmail,
			Date: "刚刚", // 时间格式化函数
		}
	}
	//处理raw数据

	c.JSON(http.StatusOK, gin.H{"data": contacts, "code": 200})
	//返回raw数据

}
