package profiles

import (
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/userData"
)

func UpUserName(email string, userName userData.UserName) error {
	type UpdateData struct {
		UserName string `xorm:"'username'"` // 映射到数据库的username列
	}
	_, err := app.Engine.Table("user_req").Where("email=?", email).Update(userName)
	_, err = app.Engine.Table("url_share").Where("email=?", email).Update(&UpdateData{UserName: userName.Name})
	if err != nil {
		return err
	}
	return nil
}
