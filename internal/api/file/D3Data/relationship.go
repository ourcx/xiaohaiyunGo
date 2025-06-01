package D3

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/RelationShip"
	"xiaohaiyun/internal/repositories"
	cosFile "xiaohaiyun/internal/utils/cos"
)

type simperUser struct {
	Name  string `json:"name"`
	Email string `json:"mail"`
	ID    int    `json:"id"`
}

type Connection struct {
	Source int `json:"source"`
	Target int `json:"target"`
}

func SearchRelationshipByID(c *gin.Context) {
	userID := cosFile.GetID(c).ID

	var re []RelationShip.Friend

	err := app.Engine.Table("friends").Where("user_id=?", userID).Find(&re)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	var wg sync.WaitGroup
	s := repositories.NewUserRepository(app.Engine)
	ch := make(chan simperUser, len(re))

	wg.Add(1)
	fmt.Print(re)
	go func() {
		defer wg.Done()
		// 并发处理逻辑可以写在这里
		for i, r := range re {
			processedData, err := s.GetUserByID(r.FriendId)
			if err != nil {
				continue
			}
			ch <- simperUser{
				Email: processedData.Email,
				Name:  processedData.Name,
				ID:    i + 1,
			}
		}

	}()
	wg.Wait()

	// 通道一定要关闭
	close(ch)

	var processedData []simperUser
	for data := range ch {
		processedData = append(processedData, data)
	}

	// 生成连接关系
	connections := make([]Connection, len(processedData))
	for i, user := range processedData {
		connections[i] = Connection{
			Source: 0,       // 当前用户
			Target: user.ID, // 连接目标用户ID
		}
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"nodes":   processedData,
		"links":   connections,
	})
}
