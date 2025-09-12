# 循环导入问题修复说明

## 问题描述
之前的代码存在循环导入问题：
- `svc` 包导入了 `repository` 包
- `repository` 包导入了 `svc` 包
- 形成了循环依赖：`svc` → `repository` → `svc`

## 修复方案

### 1. 重构 ServiceContext
- 移除了 ServiceContext 中对具体 repository 实例的直接引用
- 改为使用 Repositories 结构体来管理所有 repository

### 2. 修改 Repository 构造函数
- 所有 repository 构造函数不再依赖 ServiceContext
- 改为直接接收 `*gorm.DB` 参数
- 移除了对 `svc` 包的导入

### 3. 创建 Repository 管理器
- 新增 `repository/repository.go` 文件
- 创建 `Repositories` 结构体来统一管理所有 repository
- 提供 `NewRepositories` 函数来初始化所有 repository

### 4. 更新 Logic 层
- 修改 LoginLogic 使用新的 repository 访问方式
- 通过 `svcCtx.Repos.User` 访问 UserRepository
- 添加了完整的用户登录验证逻辑

## 修复的文件

### 修改的文件：
- `backend/internal/svc/servicecontext.go`
- `backend/internal/repository/user_repository.go`
- `backend/internal/repository/role_repository.go`
- `backend/internal/repository/permission_repository.go`
- `backend/internal/repository/role_permission_repository.go`
- `backend/internal/logic/LoginLogic.go`

### 新增的文件：
- `backend/internal/repository/repository.go`

## 修复结果
- ✅ 消除了循环导入问题
- ✅ 代码能够正常编译
- ✅ 保持了原有的功能不变
- ✅ 改善了代码架构，降低了耦合度

## 使用方式
现在在 logic 层中使用 repository 的方式：
```go
// 之前（有循环导入问题）
user, err := l.svcCtx.UserRepo.GetByUsername(ctx, username)

// 现在（无循环导入）
user, err := l.svcCtx.Repos.User.GetByUsername(ctx, username)
```