package logic

import (
	"context"
	"time"

	"model_mall_backend/backend/internal/middleware"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// TODO: 验证用户名密码
	userId := int64(1) // 这里应该是从数据库获取的用户ID

	// 生成 JWT token
	token, err := l.generateToken(userId)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		AccessToken: token,
		ExpireTime:  time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

// 生成 JWT token
func (l *LoginLogic) generateToken(userId int64) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &middleware.UserClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	return claims.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
}
