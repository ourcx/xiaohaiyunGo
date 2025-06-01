package models

import "time"

type Message struct {
	Username string `json:"username" xorm:"varchar(255) not null 'username'"`
	Text     string `json:"text"     xorm:"text not null 'text'"`
	Email    string `json:"email"    xorm:"varchar(255) 'email'"`
	Type     string `json:"type"     xorm:"varchar(50) not null 'type'"`
	ToUser   string `json:"toUser"   xorm:"varchar(255) not null 'to_user'"`
}

type PendingMessage struct {
	ID        int              `json:"id" xorm:"pk autoincr 'id'"`
	Message   `xorm:"extends"` // 嵌入 Message 结构
	CreatedAt time.Time        `json:"createdAt"  xorm:"created_at datetime not null"`
	Status    string           `json:"status"     xorm:"varchar(20) default 'pending'"`
}

// TableName 声明表名（可选）
func (PendingMessage) TableName() string {
	return "pending_messages"
}
