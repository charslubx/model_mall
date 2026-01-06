// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"
	"time"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户信息
func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentUserLogic) GetCurrentUser() (resp *types.UserInfo, err error) {
	// 从context中获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 查询用户信息
	user, err := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 构造响应
	resp = &types.UserInfo{
		Id:        fmt.Sprintf("%d", user.ID),
		Name:      user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		UserType:  getUserType(user.RoleID),
		Status:    getStatusText(user.Status),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	if user.LastLoginAt != nil {
		resp.LastLoginAt = user.LastLoginAt.Format(time.RFC3339)
	}

	return resp, nil
}

func getStatusText(status int8) string {
	if status == 1 {
		return "active"
	}
	return "disabled"
}
