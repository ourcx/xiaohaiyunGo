package locationFIleName

import "time"

type Folder struct {
	Id        int       `xorm:"pk autoincr  'id'"`         // 主键自增长
	UserId    int       `xorm:"not null index  'user_id'"` // 外键关联user.id
	FileName  []string  `xorm:"json not null 'file_name'"` // 带扩展名的文件名
	FileType  string    `xorm:"varchar(50) 'file_type'"`   // 自动提取的扩展名
	CreatedAt time.Time `xorm:"created 'created_at'"`      // 自动记录创建时间
	UpdatedAt time.Time `xorm:"updated 'updated_at'"`      // 自动记录更新时间
}

// 表名映射
func (f *Folder) TableName() string {
	return "folder"
}
