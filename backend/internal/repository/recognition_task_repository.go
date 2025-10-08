package repository

import (
	"model_mall_backend/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type RecognitionTaskRepository struct {
	db *gorm.DB
}

func NewRecognitionTaskRepository(db *gorm.DB) *RecognitionTaskRepository {
	return &RecognitionTaskRepository{db: db}
}

// Create 创建识别任务
func (r *RecognitionTaskRepository) Create(task *models.RecognitionTask) error {
	return r.db.Create(task).Error
}

// GetByID 根据ID获取任务
func (r *RecognitionTaskRepository) GetByID(id int64) (*models.RecognitionTask, error) {
	var task models.RecognitionTask
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetByTaskID 根据TaskID获取任务
func (r *RecognitionTaskRepository) GetByTaskID(taskID string) (*models.RecognitionTask, error) {
	var task models.RecognitionTask
	err := r.db.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetByImageID 根据图片ID获取任务
func (r *RecognitionTaskRepository) GetByImageID(imageID int64) (*models.RecognitionTask, error) {
	var task models.RecognitionTask
	err := r.db.Where("image_id = ?", imageID).Order("created_at DESC").First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetByUserID 获取用户的任务列表
func (r *RecognitionTaskRepository) GetByUserID(userID int64, page, pageSize int) ([]models.RecognitionTask, int64, error) {
	var tasks []models.RecognitionTask
	var total int64

	offset := (page - 1) * pageSize
	
	if err := r.db.Model(&models.RecognitionTask{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&tasks).Error

	return tasks, total, err
}

// Update 更新任务
func (r *RecognitionTaskRepository) Update(task *models.RecognitionTask) error {
	return r.db.Save(task).Error
}

// UpdateStatus 更新任务状态
func (r *RecognitionTaskRepository) UpdateStatus(taskID string, status int16) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	
	if status == models.TaskStatusProcessing {
		updates["started_at"] = time.Now()
	} else if status == models.TaskStatusCompleted || status == models.TaskStatusFailed {
		updates["completed_at"] = time.Now()
	}
	
	return r.db.Model(&models.RecognitionTask{}).
		Where("task_id = ?", taskID).
		Updates(updates).Error
}

// UpdateProgress 更新任务进度
func (r *RecognitionTaskRepository) UpdateProgress(taskID string, progress int) error {
	return r.db.Model(&models.RecognitionTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"progress":   progress,
			"updated_at": time.Now(),
		}).Error
}

// UpdateError 更新错误信息
func (r *RecognitionTaskRepository) UpdateError(taskID string, errorMsg string) error {
	return r.db.Model(&models.RecognitionTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":        models.TaskStatusFailed,
			"error_message": errorMsg,
			"completed_at":  time.Now(),
			"updated_at":    time.Now(),
		}).Error
}
