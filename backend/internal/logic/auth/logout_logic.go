// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"

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
	// 从context中获取token（在中间件中设置）
	token, ok := l.ctx.Value("token").(string)
	if !ok || token == "" {
		logx.Error("无法获取token")
		return nil, fmt.Errorf("无效的token")
	}

	// 获取用户ID用于日志记录
	userId, _ := l.ctx.Value("userId").(int64)

	// 将token加入黑名单（Redis）
	// 设置过期时间为token的剩余有效期（这里设置为24小时）
	err = l.svcCtx.RedisHelper.GetClient().SetexCtx(l.ctx, "token_blacklist:"+token, "1", 86400)
	if err != nil {
		logx.Errorf("将token加入黑名单失败: %v", err)
		return nil, fmt.Errorf("登出失败")
	}

	logx.Infof("用户 %d 登出成功，token已加入黑名单", userId)

	resp = &types.BaseResponse{
		Code:    200,
		Message: "登出成功",
		Data:    map[string]bool{"success": true},
	}

	return resp, nil
}
