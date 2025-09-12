-- 创建图片分类记录表
CREATE TABLE IF NOT EXISTS image_classifications (
    id BIGSERIAL PRIMARY KEY,
    image_path VARCHAR(500) NOT NULL,
    image_name VARCHAR(255) NOT NULL,
    image_size BIGINT,
    image_format VARCHAR(20),
    model_name VARCHAR(100) NOT NULL,
    model_version VARCHAR(50),
    process_time BIGINT,
    confidence DECIMAL(5,4),
    status SMALLINT DEFAULT 1 CHECK (status IN (0, 1, 2)),
    user_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建图片分类标签表
CREATE TABLE IF NOT EXISTS image_classification_labels (
    id BIGSERIAL PRIMARY KEY,
    classification_id BIGINT NOT NULL,
    label_name VARCHAR(100) NOT NULL,
    label_code VARCHAR(100),
    confidence DECIMAL(5,4) NOT NULL,
    bounding_box TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_image_classifications_user_id ON image_classifications(user_id);
CREATE INDEX idx_image_classifications_model_name ON image_classifications(model_name);
CREATE INDEX idx_image_classifications_status ON image_classifications(status);
CREATE INDEX idx_image_classifications_created_at ON image_classifications(created_at);

CREATE INDEX idx_image_classification_labels_classification_id ON image_classification_labels(classification_id);
CREATE INDEX idx_image_classification_labels_label_name ON image_classification_labels(label_name);
CREATE INDEX idx_image_classification_labels_confidence ON image_classification_labels(confidence);

-- 添加外键约束
ALTER TABLE image_classifications 
ADD CONSTRAINT fk_image_classifications_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON UPDATE CASCADE ON DELETE SET NULL;

ALTER TABLE image_classification_labels 
ADD CONSTRAINT fk_image_classification_labels_classification_id 
FOREIGN KEY (classification_id) REFERENCES image_classifications(id) 
ON UPDATE CASCADE ON DELETE CASCADE;

-- 添加注释
COMMENT ON TABLE image_classifications IS '图片分类记录表';
COMMENT ON COLUMN image_classifications.id IS '分类记录ID';
COMMENT ON COLUMN image_classifications.image_path IS '图片路径';
COMMENT ON COLUMN image_classifications.image_name IS '图片名称';
COMMENT ON COLUMN image_classifications.image_size IS '图片大小(字节)';
COMMENT ON COLUMN image_classifications.image_format IS '图片格式';
COMMENT ON COLUMN image_classifications.model_name IS '使用的模型名称';
COMMENT ON COLUMN image_classifications.model_version IS '模型版本';
COMMENT ON COLUMN image_classifications.process_time IS '处理耗时(毫秒)';
COMMENT ON COLUMN image_classifications.confidence IS '总体置信度';
COMMENT ON COLUMN image_classifications.status IS '状态 0-失败 1-成功 2-处理中';
COMMENT ON COLUMN image_classifications.user_id IS '用户ID';
COMMENT ON COLUMN image_classifications.created_at IS '创建时间';
COMMENT ON COLUMN image_classifications.updated_at IS '更新时间';

COMMENT ON TABLE image_classification_labels IS '图片分类标签表';
COMMENT ON COLUMN image_classification_labels.id IS '标签ID';
COMMENT ON COLUMN image_classification_labels.classification_id IS '分类记录ID';
COMMENT ON COLUMN image_classification_labels.label_name IS '标签名称';
COMMENT ON COLUMN image_classification_labels.label_code IS '标签代码';
COMMENT ON COLUMN image_classification_labels.confidence IS '置信度';
COMMENT ON COLUMN image_classification_labels.bounding_box IS '边界框信息(JSON格式)';
COMMENT ON COLUMN image_classification_labels.created_at IS '创建时间';