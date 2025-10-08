-- 识别任务表
CREATE TABLE IF NOT EXISTS recognition_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_id VARCHAR(100) UNIQUE NOT NULL, -- 任务唯一标识
    image_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    model_name VARCHAR(100), -- 使用的模型名称
    status SMALLINT DEFAULT 0, -- 0-待处理 1-处理中 2-已完成 3-失败
    progress INTEGER DEFAULT 0, -- 进度 0-100
    result_count INTEGER DEFAULT 0, -- 识别结果数量
    error_message TEXT, -- 错误信息
    started_at TIMESTAMP, -- 开始时间
    completed_at TIMESTAMP, -- 完成时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_recognition_tasks_image_id FOREIGN KEY (image_id) REFERENCES images(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_recognition_tasks_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX idx_recognition_tasks_task_id ON recognition_tasks(task_id);
CREATE INDEX idx_recognition_tasks_image_id ON recognition_tasks(image_id);
CREATE INDEX idx_recognition_tasks_user_id ON recognition_tasks(user_id);
CREATE INDEX idx_recognition_tasks_status ON recognition_tasks(status);
CREATE INDEX idx_recognition_tasks_created_at ON recognition_tasks(created_at);

COMMENT ON TABLE recognition_tasks IS '识别任务表';
COMMENT ON COLUMN recognition_tasks.id IS '任务ID';
COMMENT ON COLUMN recognition_tasks.task_id IS '任务唯一标识';
COMMENT ON COLUMN recognition_tasks.image_id IS '图片ID';
COMMENT ON COLUMN recognition_tasks.user_id IS '用户ID';
COMMENT ON COLUMN recognition_tasks.model_name IS '使用的模型名称';
COMMENT ON COLUMN recognition_tasks.status IS '状态 0-待处理 1-处理中 2-已完成 3-失败';
COMMENT ON COLUMN recognition_tasks.progress IS '进度 0-100';
COMMENT ON COLUMN recognition_tasks.result_count IS '识别结果数量';
COMMENT ON COLUMN recognition_tasks.error_message IS '错误信息';
COMMENT ON COLUMN recognition_tasks.started_at IS '开始时间';
COMMENT ON COLUMN recognition_tasks.completed_at IS '完成时间';
COMMENT ON COLUMN recognition_tasks.created_at IS '创建时间';
COMMENT ON COLUMN recognition_tasks.updated_at IS '更新时间';
