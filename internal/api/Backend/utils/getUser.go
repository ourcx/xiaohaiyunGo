package utilsBack

import (
	"xiaohaiyun/internal/api/Backend/models"
	"xiaohaiyun/internal/app"
)

func GetUser(email string) (models.BackendLogin, error, bool) {
	var getVal models.BackendLogin
	exists, err := app.Engine.Table("BackendLogin").Where("Email=?", email).Get(&getVal)

	if err != nil {
		return getVal, err, exists
	} else if exists {
		return getVal, err, exists
	}

	return getVal, nil, exists
}
