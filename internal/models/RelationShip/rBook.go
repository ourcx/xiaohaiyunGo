package RelationShip

type Contact struct {
	Src  string `json:"src"`  // 资源URL
	Name string `json:"name"` // 联系人名称
	Date string `json:"date"` // 时间描述
	ID   string `json:"id"`   // 用户ID（邮箱格式）
	//GroupID string `json:"group_id"` // 群组ID（UUID）
	//暂时不开放群组的功能
}
