-- 分类标签表
CREATE TABLE IF NOT EXISTS classification_labels (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL,
    image_id BIGINT NOT NULL,
    label_name VARCHAR(100) NOT NULL, -- 标签名称
    label_code VARCHAR(100), -- 标签代码
    confidence DECIMAL(5,4), -- 置信度 0.0000-1.0000
    bbox_x INTEGER, -- 边界框X坐标
    bbox_y INTEGER, -- 边界框Y坐标
    bbox_width INTEGER, -- 边界框宽度
    bbox_height INTEGER, -- 边界框高度
    extra_data JSONB, -- 额外数据（JSON格式）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_classification_labels_task_id FOREIGN KEY (task_id) REFERENCES recognition_tasks(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_classification_labels_image_id FOREIGN KEY (image_id) REFERENCES images(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX idx_classification_labels_task_id ON classification_labels(task_id);
CREATE INDEX idx_classification_labels_image_id ON classification_labels(image_id);
CREATE INDEX idx_classification_labels_label_name ON classification_labels(label_name);
CREATE INDEX idx_classification_labels_label_code ON classification_labels(label_code);
CREATE INDEX idx_classification_labels_confidence ON classification_labels(confidence);
CREATE INDEX idx_classification_labels_created_at ON classification_labels(created_at);

COMMENT ON TABLE classification_labels IS '分类标签表';
COMMENT ON COLUMN classification_labels.id IS '标签ID';
COMMENT ON COLUMN classification_labels.task_id IS '任务ID';
COMMENT ON COLUMN classification_labels.image_id IS '图片ID';
COMMENT ON COLUMN classification_labels.label_name IS '标签名称';
COMMENT ON COLUMN classification_labels.label_code IS '标签代码';
COMMENT ON COLUMN classification_labels.confidence IS '置信度 0.0000-1.0000';
COMMENT ON COLUMN classification_labels.bbox_x IS '边界框X坐标';
COMMENT ON COLUMN classification_labels.bbox_y IS '边界框Y坐标';
COMMENT ON COLUMN classification_labels.bbox_width IS '边界框宽度';
COMMENT ON COLUMN classification_labels.bbox_height IS '边界框高度';
COMMENT ON COLUMN classification_labels.extra_data IS '额外数据（JSON格式）';
COMMENT ON COLUMN classification_labels.created_at IS '创建时间';
COMMENT ON COLUMN classification_labels.updated_at IS '更新时间';
