package share

import (
	"time"
)

type UrlShare struct {
	ID        []byte `xorm:"BINARY(16) pk not null 'id'"`      // 主键字段
	Url       string `xorm:"VARCHAR(2048) 'url'"`              // URL地址
	Username  string `xorm:"VARCHAR(100) not null 'username'"` // 用户名
	Signature string `xorm:"TEXT 'signature'"`                 // 用户签名
	Email     string `xorm:"VARCHAR(255) 'email'"`             // 邮箱
	UserReqID int    `xorm:"INT 'user_req_id'"`                // 关联字段
	Avatar    string `xorm:"VARCHAR(2048) 'avatar'"`
}

type UrlShareString struct {
	ID        string `xorm:"BINARY(16) pk not null 'id'"`      // 主键字段
	Url       string `xorm:"VARCHAR(2048) 'url'"`              // URL地址
	Username  string `xorm:"VARCHAR(100) not null 'username'"` // 用户名
	Signature string `xorm:"TEXT 'signature'"`                 // 用户签名
	Email     string `xorm:"VARCHAR(255) 'email'"`             // 邮箱
	UserReqID int    `xorm:"INT 'user_req_id'"`                // 关联字段
	Avatar    string `xorm:"VARCHAR(2048) 'avatar'"`
}

type UrlShareJSON struct {
	Url       string `json:"url"`
	UpdatedAt string `json:"updated"`
}

type UrlData struct {
	ID         int64     `xorm:"BIGINT pk autoincr 'id'"`
	ShareID    []byte    `xorm:"BINARY(16) 'share_id'"` // BINARY(16) 存储 UUID
	Files      []string  `xorm:"JSON NOT NULL 'files'"` // JSON 数组
	Password   string    `xorm:"VARCHAR(255) NOT NULL 'password'"`
	ExpiresAt  time.Time `xorm:"TIMESTAMP NOT NULL 'expires_at'"` // 明确时区
	CreatedAt  time.Time `xorm:"created_at created"`              // 自动填充创建时间
	VisitCount int       `xorm:"INT NOT NULL DEFAULT 0 'visit_count'"`
	OneId      []byte    `xorm:"BINARY(16) 'one_id'"` // 关键字段定义
}
type UrlDataJSONString struct {
	Files     []string  `xorm:"JSON notnull" json:"files"`
	Password  string    `xorm:"VARCHAR(255) 'password'" json:"password"`
	ExpiresAt time.Time `xorm:"expires_at notnull" json:"expiresAt"`
}

type GetUrlDataJSONString struct {
	OneId string `xorm:"VARCHAR(16) 'one_id'" json:"one_id"`
}

type CheckUrlDataJSONString struct {
	OneId    string `xorm:"VARCHAR(16) 'one_id'" json:"one_id"`
	Password string `xorm:"VARCHAR(255) 'password'" json:"password"`
}

func (f *UrlData) TableName() string {
	return "url_data"
}
