package models

import "time"

type User struct {
	ID       int64  `xorm:"pk autoincr 'id'"`
	Name     int64  `xorm:"varchar(10) not null 'name'"`
	Password string `xorm:"int not null 'password'"`
	Email    string `xorm:"varchar(40) not null 'email'"`
	//File         int64  `xorm:"'file'"`
	//Relationship string `xorm:"char(1) 'relationship'"`
	//Data         string `xorm:"varchar(500) 'data'"`
}
type UserReq struct {
	ID       int    `xorm:"'id' pk autoincr"` // 主键自增
	Name     string `xorm:"'name' varchar(10) not null"`
	Password string `xorm:"'password' text not null"` // 存储哈希后的密码
	Email    string `xorm:"'email' varchar(40) not null"`
}

type UserReqByEmail struct {
	ID       int    `xorm:"'id' pk autoincr"` // 主键自增
	Name     string `xorm:"'name' varchar(10) not null"`
	Password string `xorm:"'password' text not null"` // 存储哈希后的密码
	Email    string `xorm:"'email' varchar(40) not null"`
	Code     string `xorm:"'code' varchar(6) not null"`
}

// SafeUserReq 用于返回的脱敏结构体（避免暴露密码哈希）
type SafeUserReq struct {
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UserloginPost struct {
	Password string `xorm:"'password' text not null"` // 存储哈希后的密码
	Email    string `xorm:"'email' varchar(40) not null"`
}

type UserLogin struct {
	Id        int       `xorm:"pk autoincr 'id'"`
	Email     string    `xorm:"pk varchar(255)"`                          // 主键
	LoginTime time.Time `xorm:"timestamp not null default 'now()' index"` // 自动时间戳+索引
	LoginIp   string    `xorm:"varchar(45) not null"`                     // 兼容 IPv6
}

// TableName 方法用于返回表名
func (u User) TableName() string {
	return "user"
}
