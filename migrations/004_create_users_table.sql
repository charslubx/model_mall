-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255),
    nickname VARCHAR(50),
    gender SMALLINT DEFAULT 0,
    birthday DATE,
    status SMALLINT DEFAULT 1,
    role_id BIGINT NOT NULL,
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 外键约束
    CONSTRAINT fk_users_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- 创建索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_created_at ON users(created_at);

-- 创建更新时间触发器
CREATE TRIGGER update_users_updated_at BEFORE UPDATE
    ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 插入系统默认用户数据
-- 注意：密码为明文 'admin123' 的 bcrypt 哈希值
INSERT INTO users (username, email, password, nickname, role_id, status) VALUES
('admin', 'admin@modelmall.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKXgwHNjFQSa0CzJdCQfKJHmqQf2', '系统管理员', 1, 1),
('manager', 'manager@modelmall.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKXgwHNjFQSa0CzJdCQfKJHmqQf2', '管理员', 2, 1),
('user', 'user@modelmall.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKXgwHNjFQSa0CzJdCQfKJHmqQf2', '普通用户', 3, 1);

COMMENT ON TABLE users IS '用户表';