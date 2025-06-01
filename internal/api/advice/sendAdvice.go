package advice

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"log"
	config "xiaohaiyun/configs"
	cosFile "xiaohaiyun/internal/utils/cos"
)

type advice struct {
	content string
}

func SendAdvice(c *gin.Context) {
	var ad advice
	email := cosFile.GetID(c).Email
	err := c.ShouldBind(&ad)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    ad,
		})
		return
	}

	fmt.Println(ad.content)
	m := gomail.NewMessage()
	m.SetAddressHeader("From", config.Conf.Email.SmtpUser, "小海云后台")
	m.SetHeader("To", "3277975910@qq.com")
	m.SetHeader("Subject", "【小海云】用户建议收集")
	m.SetBody("text/html", fmt.Sprintf(`
        <div style="font-family: 'Microsoft YaHei', sans-serif; max-width: 600px; margin: 20px auto;">
            <h3 style="color: #27ba9b;">建议内容</h3>
            <p><strong style="font-size: 24px;">%s</strong></p>
            <p>来自%s的建议</p>
            <hr style="border-color: #eee;">
        </div>
    `, ad.content, email))

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
		"message": "已发送",
	})
}
