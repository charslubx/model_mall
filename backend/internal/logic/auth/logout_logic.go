// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登出
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() (resp *types.BaseResponse, err error) {
	// 登出逻辑
	// 在实际应用中,可以将token加入黑名单(Redis)
	// 这里简化处理,直接返回成功

	resp = &types.BaseResponse{
		Code:    200,
		Message: "登出成功",
		Data:    map[string]bool{"success": true},
	}

	return resp, nil
}
