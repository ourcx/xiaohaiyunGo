package services

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models"
	"xiaohaiyun/internal/repositories"
	"xiaohaiyun/internal/utils"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(engine *xorm.Engine) *UserService {
	return &UserService{userRepo: repositories.NewUserRepository(engine)}
}

func (us *UserService) GetUsers() ([]*models.User, error) {
	return us.userRepo.GetUsers()
}

func CheckPassword(email, inputPassword string) (bool, error) {
	// 在业务层（如 Service 或 Handler）中调用
	s := repositories.NewUserRepository(app.Engine)
	user, err := s.GetUserByEmail(email)
	if err != nil {
		fmt.Print("查询用户失败")
		return false, err
	}
	if user == nil {
		fmt.Print("用户缺失")
		return false, nil // 用户不存在
	}
	//fmt.Print(inputPassword)
	//fmt.Print("-------")
	//fmt.Print(user.Password + "\n")
	// 校验密码哈希
	isValid := utils.CheckPassword(inputPassword, user.Password)
	if !isValid {
		fmt.Print("失败,密码错误")
		// 密码错误
		return false, nil
	}
	return true, nil
}
