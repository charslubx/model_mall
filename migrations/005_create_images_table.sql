-- 创建图片分类表
CREATE TABLE IF NOT EXISTS images (
    id BIGSERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL COMMENT '原始文件名',
    file_path VARCHAR(500) NOT NULL COMMENT '文件存储路径',
    file_size BIGINT NOT NULL COMMENT '文件大小(字节)',
    mime_type VARCHAR(100) NOT NULL COMMENT '文件MIME类型',
    width INTEGER COMMENT '图片宽度',
    height INTEGER COMMENT '图片高度',
    uploaded_by BIGINT COMMENT '上传用户ID',
    status SMALLINT DEFAULT 0 NOT NULL COMMENT '状态：0-处理中 1-已分类 2-失败',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建图片分类标签表
CREATE TABLE IF NOT EXISTS image_classifications (
    id BIGSERIAL PRIMARY KEY,
    image_id BIGINT NOT NULL COMMENT '图片ID',
    label VARCHAR(100) NOT NULL COMMENT '分类标签',
    confidence DECIMAL(5,4) NOT NULL COMMENT '置信度(0-1)',
    model_name VARCHAR(100) NOT NULL COMMENT '使用的模型名称',
    model_version VARCHAR(50) COMMENT '模型版本',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX idx_images_uploaded_by ON images(uploaded_by);
CREATE INDEX idx_images_status ON images(status);
CREATE INDEX idx_images_created_at ON images(created_at);
CREATE INDEX idx_image_classifications_image_id ON image_classifications(image_id);
CREATE INDEX idx_image_classifications_label ON image_classifications(label);
CREATE INDEX idx_image_classifications_confidence ON image_classifications(confidence DESC);

-- 添加外键约束
ALTER TABLE images ADD CONSTRAINT fk_images_uploaded_by 
    FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE SET NULL;

-- 添加注释
COMMENT ON TABLE images IS '图片信息表';
COMMENT ON TABLE image_classifications IS '图片分类标签表';
COMMENT ON COLUMN images.status IS '状态：0-处理中 1-已分类 2-失败';
COMMENT ON COLUMN image_classifications.confidence IS '置信度范围0-1，值越高表示模型越确信';