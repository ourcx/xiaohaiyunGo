package profiles

import (
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/userData"
)

func UpUserName(email string, userName userData.UserName) error {
	_, err := app.Engine.Table("user_req").Where("email=?", email).Update(userName)
	if err != nil {
		return err
	}
	return nil
}
