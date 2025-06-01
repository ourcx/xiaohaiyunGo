package userData

type UserProfile struct {
	ProfileId int    `xorm:"'profile_id' pk autoincr"`                                 // 主键自增，映射 profile_id 字段
	UserId    int    `xorm:"'user_id' unique notnull"`                                 // 唯一且非空，映射 user_id 字段
	Signature string `xorm:"'signature' varchar(255) default('')"`                     // 签名，默认空字符串
	AvatarUrl string `xorm:"'avatar_url' varchar(512) default('/default-avatar.png')"` // 头像地址

	// 以下字段用于维护关联关系（非表字段）
	// UserReq   *UserReq `xorm:"-"` // 假设关联的 UserReq 结构体，- 表示不映射到数据库
}

type PostProfile struct {
	Signature string `json:"signature" validate:"max=255"`
	AvatarUrl string `json:"avatar_url" validate:"omitempty,url,max=512"`
}

type GetProfile struct {
	Email     string `json:"'email'"`
	Signature string `json:"signature" validate:"max=255"`
	AvatarUrl string `json:"avatar_url" validate:"omitempty,url,max=512"`
	UserName  string `json:"name"`
}

type UserName struct {
	Name string `json:"name"`
}

// TableName 明确绑定到 user_profiles 表
func (UserProfile) TableName() string {
	return "user_profiles"
}

//把这个表明确绑定到对应的数据库表
