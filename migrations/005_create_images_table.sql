-- 图片表
CREATE TABLE IF NOT EXISTS images (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_url VARCHAR(500),
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    width INTEGER,
    height INTEGER,
    md5 VARCHAR(32),
    status SMALLINT DEFAULT 1, -- 0-已删除 1-正常
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_images_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX idx_images_user_id ON images(user_id);
CREATE INDEX idx_images_status ON images(status);
CREATE INDEX idx_images_md5 ON images(md5);
CREATE INDEX idx_images_created_at ON images(created_at);

COMMENT ON TABLE images IS '图片表';
COMMENT ON COLUMN images.id IS '图片ID';
COMMENT ON COLUMN images.user_id IS '用户ID';
COMMENT ON COLUMN images.filename IS '文件名';
COMMENT ON COLUMN images.original_name IS '原始文件名';
COMMENT ON COLUMN images.file_path IS '文件路径';
COMMENT ON COLUMN images.file_url IS '文件URL';
COMMENT ON COLUMN images.file_size IS '文件大小（字节）';
COMMENT ON COLUMN images.mime_type IS 'MIME类型';
COMMENT ON COLUMN images.width IS '图片宽度';
COMMENT ON COLUMN images.height IS '图片高度';
COMMENT ON COLUMN images.md5 IS '文件MD5';
COMMENT ON COLUMN images.status IS '状态 0-已删除 1-正常';
COMMENT ON COLUMN images.created_at IS '创建时间';
COMMENT ON COLUMN images.updated_at IS '更新时间';
