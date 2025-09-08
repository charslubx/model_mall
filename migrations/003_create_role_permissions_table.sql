-- 创建角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL COMMENT '角色ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    -- 外键约束
    CONSTRAINT fk_role_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON UPDATE CASCADE ON DELETE CASCADE,
    
    -- 唯一约束，防止重复分配
    CONSTRAINT uk_role_permission UNIQUE (role_id, permission_id)
);

-- 创建索引
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- 为超级管理员分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions WHERE is_system = TRUE;

-- 为管理员分配部分权限（除了删除权限）
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions 
WHERE is_system = TRUE 
AND code NOT LIKE '%:delete%' 
AND code NOT IN ('system:permission:create', 'system:permission:update', 'system:permission:delete');

-- 为普通用户分配基础权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions 
WHERE code IN (
    'profile',
    'profile:info',
    'profile:update', 
    'profile:change_password'
);

-- 为访客分配查看权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 4, id FROM permissions 
WHERE code IN (
    'profile',
    'profile:info'
);

COMMENT ON TABLE role_permissions IS '角色权限关联表';