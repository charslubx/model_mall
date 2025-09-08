# 数据库设计文档

## 概述

本项目采用基于角色的访问控制（RBAC）模型，设计了用户、角色、权限三个核心表，以及角色权限关联表。

## 表结构设计

### 1. 用户表 (users)

用户基础信息表，与角色表为一对一关系。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | 主键 | 用户ID |
| username | VARCHAR(50) | 唯一,非空 | 用户名 |
| email | VARCHAR(100) | 唯一,非空 | 邮箱 |
| phone | VARCHAR(20) | 唯一 | 手机号 |
| password | VARCHAR(255) | 非空 | 密码哈希 |
| avatar | VARCHAR(255) | | 头像URL |
| nickname | VARCHAR(50) | | 昵称 |
| gender | SMALLINT | 默认0 | 性别 0-未知 1-男 2-女 |
| birthday | DATE | | 生日 |
| status | SMALLINT | 默认1 | 状态 0-禁用 1-正常 |
| role_id | BIGINT | 非空,外键 | 角色ID |
| last_login_at | TIMESTAMP | | 最后登录时间 |
| last_login_ip | VARCHAR(45) | | 最后登录IP |
| created_at | TIMESTAMP | 默认当前时间 | 创建时间 |
| updated_at | TIMESTAMP | 默认当前时间 | 更新时间 |

**索引：**
- username (唯一索引)
- email (唯一索引)
- phone (唯一索引)
- status
- role_id
- created_at

**外键约束：**
- role_id → roles.id (ON UPDATE CASCADE, ON DELETE RESTRICT)

### 2. 角色表 (roles)

角色定义表，定义系统中的各种角色。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | 主键 | 角色ID |
| name | VARCHAR(50) | 唯一,非空 | 角色名称 |
| code | VARCHAR(50) | 唯一,非空 | 角色代码 |
| description | VARCHAR(255) | | 角色描述 |
| status | SMALLINT | 默认1 | 状态 0-禁用 1-正常 |
| sort | INTEGER | 默认0 | 排序 |
| is_system | BOOLEAN | 默认false | 是否系统角色 |
| created_at | TIMESTAMP | 默认当前时间 | 创建时间 |
| updated_at | TIMESTAMP | 默认当前时间 | 更新时间 |

**索引：**
- name (唯一索引)
- code (唯一索引)
- status
- sort
- is_system

**默认数据：**
- 超级管理员 (super_admin)
- 管理员 (admin)
- 普通用户 (user)
- 访客 (guest)

### 3. 权限表 (permissions)

权限定义表，支持树形结构的权限体系。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | 主键 | 权限ID |
| name | VARCHAR(50) | 非空 | 权限名称 |
| code | VARCHAR(100) | 唯一,非空 | 权限代码 |
| type | VARCHAR(20) | 非空 | 权限类型 menu/button/api |
| parent_id | BIGINT | 默认0 | 父权限ID |
| path | VARCHAR(255) | | 路径/接口地址 |
| method | VARCHAR(10) | | 请求方法 |
| icon | VARCHAR(100) | | 图标 |
| component | VARCHAR(255) | | 组件路径 |
| sort | INTEGER | 默认0 | 排序 |
| status | SMALLINT | 默认1 | 状态 0-禁用 1-正常 |
| is_system | BOOLEAN | 默认false | 是否系统权限 |
| description | VARCHAR(255) | | 权限描述 |
| created_at | TIMESTAMP | 默认当前时间 | 创建时间 |
| updated_at | TIMESTAMP | 默认当前时间 | 更新时间 |

**索引：**
- code (唯一索引)
- parent_id
- type
- status
- sort

**权限类型说明：**
- `menu`: 菜单权限，用于前端菜单显示
- `button`: 按钮权限，用于页面按钮控制
- `api`: 接口权限，用于API访问控制

### 4. 角色权限关联表 (role_permissions)

角色与权限的多对多关联表。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | 主键 | 关联ID |
| role_id | BIGINT | 非空,索引 | 角色ID |
| permission_id | BIGINT | 非空,索引 | 权限ID |
| created_at | TIMESTAMP | 默认当前时间 | 创建时间 |

**索引：**
- role_id
- permission_id
- (role_id, permission_id) 联合唯一索引

**外键约束：**
- role_id → roles.id (ON UPDATE CASCADE, ON DELETE CASCADE)
- permission_id → permissions.id (ON UPDATE CASCADE, ON DELETE CASCADE)

## 关系设计

### 1. 用户与角色 (1:1)
- 每个用户只能有一个角色
- 一个角色可以被多个用户使用
- 外键约束确保数据完整性

### 2. 角色与权限 (N:M)
- 一个角色可以有多个权限
- 一个权限可以被多个角色使用
- 通过中间表 role_permissions 建立关联

### 3. 权限树形结构
- 权限表支持自关联，通过 parent_id 字段构建树形结构
- 根权限的 parent_id 为 0
- 支持无限级权限嵌套

## 默认数据

### 默认角色
1. **超级管理员 (super_admin)**: 拥有所有权限
2. **管理员 (admin)**: 拥有大部分权限，不包括删除和权限管理
3. **普通用户 (user)**: 只有个人中心相关权限
4. **访客 (guest)**: 只有查看个人信息权限

### 默认用户
1. **admin**: 超级管理员账号，密码: admin123
2. **manager**: 管理员账号，密码: admin123
3. **user**: 普通用户账号，密码: admin123

### 系统权限
预置了完整的系统管理权限：
- 系统管理模块
- 用户管理 (增删改查、重置密码)
- 角色管理 (增删改查、分配权限)
- 权限管理 (增删改查)
- 个人中心 (查看信息、修改信息、修改密码)

## 使用说明

### 1. 数据库迁移
```bash
cd migrations
chmod +x run_migrations.sh
./run_migrations.sh
```

### 2. 手动执行SQL
```bash
cd migrations
go run migrate.go
```

### 3. 配置数据库连接
修改 `backend/etc/backend-api.yaml` 中的数据库配置：
```yaml
PostgreSQL:
  Host: localhost
  Port: 5432
  Username: postgres
  Password: your-password
  Database: model_mall
  SSLMode: disable
```

## 安全考虑

1. **密码安全**: 使用 bcrypt 哈希存储密码
2. **权限控制**: 基于角色的多层权限控制
3. **系统保护**: 系统角色和权限不可删除
4. **数据完整性**: 外键约束确保数据一致性
5. **软删除**: 可考虑实现软删除避免数据丢失

## 扩展性

1. **多租户**: 可通过添加 tenant_id 字段支持多租户
2. **权限缓存**: 可使用 Redis 缓存用户权限信息
3. **审计日志**: 可添加操作日志表记录用户行为
4. **权限继承**: 可实现角色继承机制
5. **动态权限**: 支持运行时动态添加权限