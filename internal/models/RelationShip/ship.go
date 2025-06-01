package RelationShip

import "time"

//字段说明
//status
//- pending：好友请求待确认
//- accepted：已确认好友
//- blocked：已拉黑

//要实现的功能
//好友请求
//好友列表的查看
//好友聊天沟通
//群组的功能
//好友群组文件分享
//删除好友

type Friend struct {
	FriendshipId int       `xorm:"'friendship_id' pk autoincr"`
	UserId       int       `xorm:"'user_id' notnull index(idx_user_friend)"`
	FriendId     int       `xorm:"'friend_id' notnull index(idx_friend_user)"`
	Status       string    `xorm:"'status' enum('pending','accepted','blocked') default('pending')"`
	CreatedAt    time.Time `xorm:"'created_at' created"`
	UpdatedAt    time.Time `xorm:"'updated_at' updated"`
}

// TableName 设置表名
func (Friend) TableName() string {
	return "friends"
}

type FriendEmail struct {
	Email string `json:"'email'"`
}
