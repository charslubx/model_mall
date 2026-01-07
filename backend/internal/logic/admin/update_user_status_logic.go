// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户状态
func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserStatusLogic) UpdateUserStatus(req *types.UpdateUserStatusRequest, userId string) (resp *types.UpdateUserStatusResponse, err error) {
	// 获取管理员ID
	adminId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 转换用户ID
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID")
	}

	// 查询用户信息
	user, err := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 验证状态值
	if req.Status != "active" && req.Status != "disabled" {
		return nil, fmt.Errorf("无效的状态值")
	}

	// 更新状态
	var statusInt int8 = 1
	if req.Status == "disabled" {
		statusInt = 0
	}
	user.Status = statusInt
	err = l.svcCtx.Repos.UserRepo.Update(l.ctx, user)
	if err != nil {
		return nil, fmt.Errorf("更新用户状态失败: %v", err)
	}

	logx.Infof("管理员 %d 更新用户 %s 状态为 %s，原因: %s", adminId, userId, req.Status, req.Reason)

	resp = &types.UpdateUserStatusResponse{
		Success: true,
		Status:  req.Status,
	}

	return resp, nil
}
