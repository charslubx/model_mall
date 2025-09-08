-- 创建权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限代码',
    type VARCHAR(20) NOT NULL COMMENT '权限类型 menu-菜单 button-按钮 api-接口',
    parent_id BIGINT DEFAULT 0 COMMENT '父权限ID',
    path VARCHAR(255) COMMENT '路径/接口地址',
    method VARCHAR(10) COMMENT '请求方法',
    icon VARCHAR(100) COMMENT '图标',
    component VARCHAR(255) COMMENT '组件路径',
    sort INTEGER DEFAULT 0 COMMENT '排序',
    status SMALLINT DEFAULT 1 COMMENT '状态 0-禁用 1-正常',
    is_system BOOLEAN DEFAULT FALSE COMMENT '是否系统权限',
    description VARCHAR(255) COMMENT '权限描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间'
);

-- 创建索引
CREATE INDEX idx_permissions_parent_id ON permissions(parent_id);
CREATE INDEX idx_permissions_type ON permissions(type);
CREATE INDEX idx_permissions_status ON permissions(status);
CREATE INDEX idx_permissions_sort ON permissions(sort);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_permissions_updated_at BEFORE UPDATE
    ON permissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 插入系统默认权限数据
INSERT INTO permissions (name, code, type, parent_id, path, method, icon, sort, is_system, description) VALUES
-- 系统管理
('系统管理', 'system', 'menu', 0, '/system', '', 'system', 1000, TRUE, '系统管理模块'),

-- 用户管理
('用户管理', 'system:user', 'menu', 1, '/system/user', '', 'user', 1100, TRUE, '用户管理菜单'),
('用户列表', 'system:user:list', 'api', 2, '/api/users', 'GET', '', 1101, TRUE, '获取用户列表'),
('用户详情', 'system:user:detail', 'api', 2, '/api/users/:id', 'GET', '', 1102, TRUE, '获取用户详情'),
('创建用户', 'system:user:create', 'api', 2, '/api/users', 'POST', '', 1103, TRUE, '创建用户'),
('更新用户', 'system:user:update', 'api', 2, '/api/users/:id', 'PUT', '', 1104, TRUE, '更新用户'),
('删除用户', 'system:user:delete', 'api', 2, '/api/users/:id', 'DELETE', '', 1105, TRUE, '删除用户'),
('重置密码', 'system:user:reset_password', 'button', 2, '', '', '', 1106, TRUE, '重置用户密码'),

-- 角色管理
('角色管理', 'system:role', 'menu', 1, '/system/role', '', 'role', 1200, TRUE, '角色管理菜单'),
('角色列表', 'system:role:list', 'api', 9, '/api/roles', 'GET', '', 1201, TRUE, '获取角色列表'),
('角色详情', 'system:role:detail', 'api', 9, '/api/roles/:id', 'GET', '', 1202, TRUE, '获取角色详情'),
('创建角色', 'system:role:create', 'api', 9, '/api/roles', 'POST', '', 1203, TRUE, '创建角色'),
('更新角色', 'system:role:update', 'api', 9, '/api/roles/:id', 'PUT', '', 1204, TRUE, '更新角色'),
('删除角色', 'system:role:delete', 'api', 9, '/api/roles/:id', 'DELETE', '', 1205, TRUE, '删除角色'),
('分配权限', 'system:role:assign_permission', 'button', 9, '', '', '', 1206, TRUE, '为角色分配权限'),

-- 权限管理
('权限管理', 'system:permission', 'menu', 1, '/system/permission', '', 'permission', 1300, TRUE, '权限管理菜单'),
('权限列表', 'system:permission:list', 'api', 16, '/api/permissions', 'GET', '', 1301, TRUE, '获取权限列表'),
('权限详情', 'system:permission:detail', 'api', 16, '/api/permissions/:id', 'GET', '', 1302, TRUE, '获取权限详情'),
('创建权限', 'system:permission:create', 'api', 16, '/api/permissions', 'POST', '', 1303, TRUE, '创建权限'),
('更新权限', 'system:permission:update', 'api', 16, '/api/permissions/:id', 'PUT', '', 1304, TRUE, '更新权限'),
('删除权限', 'system:permission:delete', 'api', 16, '/api/permissions/:id', 'DELETE', '', 1305, TRUE, '删除权限'),

-- 个人中心
('个人中心', 'profile', 'menu', 0, '/profile', '', 'profile', 2000, TRUE, '个人中心模块'),
('个人信息', 'profile:info', 'api', 22, '/api/profile', 'GET', '', 2001, TRUE, '获取个人信息'),
('更新个人信息', 'profile:update', 'api', 22, '/api/profile', 'PUT', '', 2002, TRUE, '更新个人信息'),
('修改密码', 'profile:change_password', 'api', 22, '/api/profile/password', 'PUT', '', 2003, TRUE, '修改密码');

COMMENT ON TABLE permissions IS '权限表';