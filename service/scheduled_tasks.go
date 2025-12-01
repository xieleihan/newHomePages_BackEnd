package service

import (
	"fmt"
	"gin/db"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

var scheduler *gocron.Scheduler

// InitScheduledTasks 初始化定时任务
func InitScheduledTasks() error {
	// 创建调度器（使用本地时区）
	scheduler = gocron.NewScheduler(time.Local)

	// 每隔 24 小时清空 encryptionmessage 表
	_, err := scheduler.Every(24).Hours().Do(ClearEncryptionMessages)
	if err != nil {
		return fmt.Errorf("添加清空加密消息任务失败: %v", err)
	}

	log.Println("✓ 定时任务已初始化")
	log.Println("✓ 任务: 每隔 24 小时清空 encryptionmessage 表")

	// 启动调度器
	scheduler.StartAsync()

	return nil
}

// ClearEncryptionMessages 清空加密消息表
func ClearEncryptionMessages() {
	log.Printf("[定时任务] 开始清空 encryptionmessage 表 (执行时间: %s)\n", time.Now().Format("2006-01-02 15:04:05"))

	// 删除所有数据
	if err := db.DB.Exec("DELETE FROM encryptionmessage").Error; err != nil {
		log.Printf("✗ 清空 encryptionmessage 表失败: %v\n", err)
		return
	}

	log.Printf("✓ 成功清空 encryptionmessage 表\n")
}

// StopScheduler 停止调度器（优雅关闭）
func StopScheduler() {
	if scheduler != nil {
		scheduler.Stop()
		log.Println("✓ 定时任务调度器已停止")
	}
}

// GetSchedulerStatus 获取调度器状态
func GetSchedulerStatus() map[string]interface{} {
	if scheduler == nil {
		return map[string]interface{}{
			"status": "未初始化",
			"jobs":   0,
		}
	}

	return map[string]interface{}{
		"status": "运行中",
		"jobs":   len(scheduler.Jobs()),
	}
}
