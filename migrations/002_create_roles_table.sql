-- 创建角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '角色代码',
    description VARCHAR(255) COMMENT '角色描述',
    status SMALLINT DEFAULT 1 COMMENT '状态 0-禁用 1-正常',
    sort INTEGER DEFAULT 0 COMMENT '排序',
    is_system BOOLEAN DEFAULT FALSE COMMENT '是否系统角色',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间'
);

-- 创建索引
CREATE INDEX idx_roles_status ON roles(status);
CREATE INDEX idx_roles_sort ON roles(sort);
CREATE INDEX idx_roles_is_system ON roles(is_system);

-- 创建更新时间触发器
CREATE TRIGGER update_roles_updated_at BEFORE UPDATE
    ON roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 插入系统默认角色数据
INSERT INTO roles (name, code, description, sort, is_system) VALUES
('超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1, TRUE),
('管理员', 'admin', '系统管理员，拥有大部分权限', 2, TRUE),
('普通用户', 'user', '普通用户，拥有基础权限', 3, TRUE),
('访客', 'guest', '访客用户，只有查看权限', 4, TRUE);

COMMENT ON TABLE roles IS '角色表';