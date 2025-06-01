package repositories

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"xiaohaiyun/internal/models"
	"xiaohaiyun/internal/models/RelationShip"
	"xiaohaiyun/internal/models/userData"
)

type UserRepository struct {
	engine *xorm.Engine
}

func NewUserRepository(engine *xorm.Engine) *UserRepository {
	return &UserRepository{engine: engine}
}

// GetUsers 获取所有用户
func (r *UserRepository) GetUsers() ([]*models.User, error) {
	var users []*models.User
	err := r.engine.Table(models.User{}.TableName()).Find(&users)
	return users, err
}

// GetUserByEmail 获取用户的邮箱对应的用户在登录表中
func (r *UserRepository) GetUserByEmail(email string) (*models.UserReq, error) {
	userReq := &models.UserReq{}
	has, err := r.engine.Where("email = ?", email).Get(userReq)
	if !has {
		return nil, nil // 用户不存在
	}
	return userReq, err
}

func (r *UserRepository) GetUserByID(id int) (*models.UserReq, error) {
	userReq := &models.UserReq{}
	has, err := r.engine.Where("ID = ?", id).Get(userReq)
	fmt.Println(has)
	if !has {
		return nil, nil // 用户不存在
	}
	return userReq, err
}

// GetTableByEmail 获得用户信息，自己传入邮箱
func (r *UserRepository) GetTableByEmail(email string) (any, error) {
	table := userData.UserProfile{}
	// 使用 JOIN 查询
	has, err := r.engine.Table("user_req").
		Join("INNER", "user_profiles", "user_req.ID = user_profiles.user_id").
		Where("user_req.email = ?", email).
		Get(&table)
	//根据req表查询个人信息表

	//has是true或者false

	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	if !has {
		return nil, fmt.Errorf("用户资料不存在")
	}

	return &table, nil
}

// GetUserRelationSHipByEmail 根据邮箱获取用户所有好友关系
func (r *UserRepository) GetUserRelationSHipByEmail(email string) ([]*RelationShip.Friend, error) {
	var friends []*RelationShip.Friend
	//拿到多条数据使用的是切片，不是单个的
	// 使用正确的字段进行查询（假设user_req表有email字段）
	err := r.engine.Table("user_req").
		Join("INNER", "friends", "user_req.ID = friends.user_id").
		Where("user_req.email = ?", email). // 修正where条件
		Find(&friends)

	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	// 当没有找到记录时返回空数组而不是错误
	if len(friends) == 0 {
		return []*RelationShip.Friend{}, nil
	}

	return friends, nil
}
