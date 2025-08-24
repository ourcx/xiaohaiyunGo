package pay

import "github.com/smartwalle/alipay/v3"

type AlipayService struct {
	Client *alipay.Client
}

type MemberObject struct {
	MemberID   string `json:"memberId"`
	MemberName string `json:"memberName"`
}

// 会员商品结构体
type MemberProduct struct {
	MemberID      string  `json:"memberId"`
	MemberName    string  `json:"memberName"`
	MemberType    string  `json:"memberType"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"originalPrice,omitempty"` // 原价（可选）
	Duration      string  `json:"duration"`                // 持续时间
	Description   string  `json:"description,omitempty"`   // 描述（可选）
	IsHot         bool    `json:"isHot"`                   // 是否热门
	IsRecommended bool    `json:"isRecommended"`           // 是否推荐
	DiscountInfo  string  `json:"discountInfo,omitempty"`  // 折扣信息（可选）
}
