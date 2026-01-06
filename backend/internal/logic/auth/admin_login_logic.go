// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"
	"time"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"
	"model_mall_backend/backend/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 管理员登录
func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLoginLogic) AdminLogin(req *types.AdminLoginRequest) (resp *types.AdminLoginResponse, err error) {
	// 根据邮箱查询用户
	user, err := l.svcCtx.Repos.UserRepo.GetByEmail(l.ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("管理员不存在或密码错误")
	}

	// 检查是否是管理员角色
	if user.RoleID != 3 {
		return nil, fmt.Errorf("无管理员权限")
	}

	// 检查账户状态
	if user.Status == 0 {
		return nil, fmt.Errorf("账户已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("管理员不存在或密码错误")
	}

	// 生成JWT token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 更新最后登录时间
	_ = l.svcCtx.Repos.UserRepo.UpdateLastLogin(l.ctx, user.ID, "127.0.0.1")

	// 构造响应
	resp = &types.AdminLoginResponse{
		Token: accessToken,
		Admin: types.AdminInfo{
			Id:    fmt.Sprintf("%d", user.ID),
			Name:  user.Nickname,
			Email: user.Email,
			Role:  "admin",
			Permissions: []string{
				"user_manage",
				"order_manage",
				"product_manage",
				"system_config",
			},
		},
	}

	return resp, nil
}

func (l *AdminLoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64, username string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["username"] = username
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
