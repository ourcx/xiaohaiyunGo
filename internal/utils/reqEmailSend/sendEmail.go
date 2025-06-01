package reqEmailSend

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"log"
	config "xiaohaiyun/configs"
	legislations "xiaohaiyun/internal/utils/legislation"
)

type EmailRequest struct {
	Email string `json:"email" binding:"required"` // 添加JSON tag
}

// SendReqEmail 邮箱授权码syafuhdndrnpdbdi
func SendReqEmail(c *gin.Context) {
	// 生成6位数字验证码
	var req EmailRequest

	// 推荐使用ShouldBindJSON处理JSON数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "参数错误: " + err.Error(),
		})
		return
	}
	email := req.Email
	//根据jwt解析邮箱
	code, _ := legislations.GenerateSecureCode(6)
	err := legislations.CodeStoreFromEmail(code, email)
	//把当前的code储存到后台的map里面

	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
		})
		return
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", config.Conf.Email.SmtpUser, "小海云后台")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "【小海云】注册账号验证码")
	m.SetBody("text/html", fmt.Sprintf(`
        <div style="font-family: 'Microsoft YaHei', sans-serif; max-width: 600px; margin: 20px auto;">
            <h3 style="color: #27ba9b;">注册账号</h3>
            <p>您的验证码：<strong style="font-size: 24px;">%s</strong></p>
            <p style="color: #ff4d4f;">有效期10分钟，请勿泄露给他人</p>
            <hr style="border-color: #eee;">
            <p style="color: #999; font-size: 12px;">
                ※ 本邮件由系统自动发送，请勿直接回复<br>
                © 2024 小海云 后台<br>
                _ 请确认本操作是你本人产生的
            </p>
        </div>
    `, code))

	d := gomail.NewDialer(
		config.Conf.Email.SmtpHost,
		config.Conf.Email.SmtpPort, // 必须使用465或587
		config.Conf.Email.SmtpUser,
		config.Conf.Email.SmtpPassword,
	)
	//创建自己的客户端
	d.TLSConfig = &tls.Config{
		ServerName:         config.Conf.Email.SmtpHost,
		InsecureSkipVerify: false,
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("邮件发送失败: %v", err)
		c.JSON(500, gin.H{
			"code":    500,
			"message": "邮件发送服务暂不可用",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "验证码已发送至绑定邮箱",
	})
}
