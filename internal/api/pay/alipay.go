package pay

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"xiaohaiyun/internal/models/pay"
)

// 会员类型常量
const (
	MemberTypeMonthly      = "monthly"
	MemberTypeQuarterly    = "quarterly"
	MemberTypeAnnually     = "annually"
	MemberTypeSubscription = "subscription"
)

// MemberProducts 会员商品列表
var MemberProducts = []pay.MemberProduct{
	{
		MemberID:      "mem_monthly_001",
		MemberName:    "月度会员",
		MemberType:    MemberTypeMonthly,
		Price:         25.00,
		Duration:      "1个月",
		Description:   "享受一个月会员特权",
		IsHot:         false,
		IsRecommended: false,
	},
	{
		MemberID:      "mem_quarterly_001",
		MemberName:    "季度会员",
		MemberType:    MemberTypeQuarterly,
		Price:         68.00,
		OriginalPrice: 75.00,
		Duration:      "3个月",
		Description:   "热门选择，超值季度会员",
		IsHot:         true,
		IsRecommended: false,
		DiscountInfo:  "节省7元",
	},
	{
		MemberID:      "mem_annually_001",
		MemberName:    "年度会员",
		MemberType:    MemberTypeAnnually,
		Price:         263.00,
		OriginalPrice: 300.00,
		Duration:      "12个月",
		Description:   "推荐选择，最优惠的年度会员",
		IsHot:         false,
		IsRecommended: true,
		DiscountInfo:  "节省37元",
	},
	{
		MemberID:      "mem_subscription_001",
		MemberName:    "连续包月",
		MemberType:    MemberTypeSubscription,
		Price:         18.00,
		OriginalPrice: 25.00,
		Duration:      "首月",
		Description:   "首月特惠，后续每月25元，可随时取消",
		IsHot:         false,
		IsRecommended: true,
		DiscountInfo:  "首月特惠",
	},
}
var client *alipay.Client

// NewAlipayService 初始化支付宝客户端
func NewAlipayService(appID, privateKey, publicKey string, isProduction bool) error {
	client, err := alipay.New(appID, privateKey, isProduction)
	if err != nil {
		return fmt.Errorf("初始化支付宝客户端失败: %v", err)
	}
	// 加载支付宝公钥
	err = client.LoadAliPayPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("加载支付宝公钥失败: %v", err)
	}
	return nil
}

func CreatePayment(orderID, subject string, amount float64) (string, error) {
	var p = alipay.TradePagePay{}
	p.NotifyURL = "https://your-domain.com/alipay/notify" // 支付宝异步通知地址
	p.ReturnURL = "https://your-domain.com/alipay/return" // 支付完成后返回地址
	p.Subject = subject                                   // 订单标题
	p.OutTradeNo = orderID                                // 商户订单号
	p.TotalAmount = fmt.Sprintf("%.2f", amount)           // 订单金额
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"              // 销售产品码，固定值

	// 生成支付URL
	url, err := client.TradePagePay(p)
	if err != nil {
		return "", fmt.Errorf("生成支付URL失败: %v", err)
	}

	return url.String(), nil
}

func FindProductByID(memberID string) (*pay.MemberProduct, error) {
	for _, product := range MemberProducts {
		if product.MemberID == memberID {
			return &product, nil
		}
	}
	return nil, fmt.Errorf("未找到会员商品: %s", memberID)
}

// AliPay 调用支付宝付款url
func AliPay(c *gin.Context) {
	appID := "你的支付宝应用APP_ID"
	privateKey := `-----BEGIN PRIVATE KEY-----
你的应用私钥
-----END PRIVATE KEY-----`
	publicKey := `-----BEGIN PUBLIC KEY-----
支付宝公钥
-----END PUBLIC KEY-----`
	isProduction := false //  true:生产环境 false:沙箱环境
	var Pay pay.MemberObject
	err := c.ShouldBind(&Pay)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
	}
	member, err := FindProductByID(Pay.MemberID)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		})
	}
	err = NewAlipayService(appID, privateKey, publicKey, isProduction)
	if err != nil {
		return
	}
	payment, err := CreatePayment(member.MemberID, member.MemberName, member.Price)
	if err != nil {
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": payment,
	})

}

// GetPayment 获得商品列表
func GetPayment(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": MemberProducts,
	})
}

// 支付宝支付通知处理 - 使用 DecodeNotification
//func alipayNotifyHandlerGin() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 读取请求体
//		body, err := io.ReadAll(c.Request.Body)
//		if err != nil {
//			log.Printf("读取请求体失败: %v", err)
//			c.JSON(http.StatusBadRequest, gin.H{
//				"error": "读取请求体失败",
//			})
//			return
//		}
//
//		// 使用 DecodeNotification 替代已弃用的 GetTradeNotification
//		notification, err := client.DecodeNotification(body)
//		if err != nil {
//			log.Printf("解析通知失败: %v", err)
//			c.JSON(http.StatusBadRequest, gin.H{
//				"error": "解析通知失败: " + err.Error(),
//			})
//			return
//		}
//
//		// 验证通知的合法性
//		if notification != nil {
//			// 通知验证成功
//			switch notification.TradeStatus {
//			case alipay.TradeStatusSuccess:
//				// 支付成功
//				log.Printf("订单 %s 支付成功，支付宝交易号: %s", notification.OutTradeNo, notification.TradeNo)
//				// 更新数据库中的订单状态为已支付
//				//err := updateOrderStatus(notification.OutTradeNo, "paid")
//				if err != nil {
//					log.Printf("更新订单状态失败: %v", err)
//				}
//
//				// 返回success告诉支付宝已成功接收通知
//				c.String(http.StatusOK, "success")
//				return
//
//			case alipay.TradeStatusFinished:
//				// 交易结束（退款或关闭）
//				log.Printf("订单 %s 交易结束", notification.OutTradeNo)
//				//err := updateOrderStatus(notification.OutTradeNo, "finished")
//				if err != nil {
//					log.Printf("更新订单状态失败: %v", err)
//				}
//
//			case alipay.TradeStatusClosed:
//				// 交易关闭
//				log.Printf("订单 %s 交易关闭", notification.OutTradeNo)
//				//err := updateOrderStatus(notification.OutTradeNo, "closed")
//				if err != nil {
//					log.Printf("更新订单状态失败: %v", err)
//				}
//
//			case alipay.TradeStatusWaitBuyerPay:
//				// 等待付款
//				log.Printf("订单 %s 等待付款", notification.OutTradeNo)
//				// 不需要更新状态，保持pending
//
//			default:
//				// 其他状态
//				log.Printf("订单 %s 状态: %s", notification.OutTradeNo, notification.TradeStatus)
//			}
//		} else {
//			// 验证失败，可能是非法请求
//			log.Printf("支付宝通知验证失败")
//			c.JSON(http.StatusBadRequest, gin.H{
//				"error": "验证失败",
//			})
//			return
//		}
//
//		// 返回success告诉支付宝已接收通知（即使处理失败也要返回success，否则支付宝会重发）
//		c.String(http.StatusOK, "success")
//	}
//}
