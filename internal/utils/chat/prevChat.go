package chat

import (
	"github.com/go-xorm/xorm"
	"time"
	"xiaohaiyun/internal/models"
)

// HasMessageByCreatedAt 检查是否存在特定 CreatedAt 的记录
func HasMessageByCreatedAt(engine *xorm.Engine, email string) bool {
	exists, _ := engine.
		Table("pending_messages").
		Where("to_user = ?", email).
		Exist(&models.PendingMessage{})
	return exists
	//判断数据库是不是有未发送的信息
}

// GetMessagesByToUser 查询对应toUser的值
// 直接获取 Message 结构列表（从 PendingMessage 表中提取 Message 部分）
func GetMessagesByToUser(engine *xorm.Engine, toUser string) ([]models.Message, error) {
	var pendingMessages []models.PendingMessage
	err := engine.
		Where("to_user = ?", toUser). // 根据 to_user 字段筛选
		Find(&pendingMessages)
	if err != nil {
		return nil, err
	}

	// 提取嵌入的 Message 结构
	messages := make([]models.Message, len(pendingMessages))
	for i, pm := range pendingMessages {
		messages[i] = pm.Message
	}

	_, err = UpdateUserMessagesStatus(engine, toUser, "pending", "sent")
	if err != nil {
		return nil, err
	}
	_, err = DeleteUserMessagesByStatus(engine, toUser, "sent")
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// SaveUnsentMessage 将未成功发送的消息存入数据库
func SaveUnsentMessage(engine *xorm.Engine, msg models.Message) error {
	pendingMsg := models.PendingMessage{
		Message:   msg,       // 嵌入原始消息内容
		Status:    "pending", // 明确设置状态（即使有默认值）
		CreatedAt: time.Now(),
	}

	// 执行插入操作
	_, err := engine.Insert(&pendingMsg)
	return err
}

// UpdateUserMessagesStatus 更改status的样子
func UpdateUserMessagesStatus(engine *xorm.Engine, toUser, oldStatus, newStatus string) (int64, error) {
	return engine.
		Where("to_user = ? AND status = ?", toUser, oldStatus).
		Update(&models.PendingMessage{Status: newStatus})
}

// DeleteUserMessagesByStatus 删除指定用户的状态消息
func DeleteUserMessagesByStatus(engine *xorm.Engine, toUser, status string) (int64, error) {
	return engine.
		Where("to_user = ? AND status = ?", toUser, status).
		Delete(&models.PendingMessage{})
}

// GetAllToUsers 获取所有唯一的 toUser 列表
func GetAllToUsers(engine *xorm.Engine) ([]string, error) {
	var toUsers []string

	// 使用 DISTINCT 查询去重
	err := engine.
		Table("pending_messages").
		Distinct("to_user").
		Select("to_user").
		Find(&toUsers)

	if err != nil {
		return nil, nil
	}

	return toUsers, nil
}
