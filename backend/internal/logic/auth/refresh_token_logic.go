// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"
	"time"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新令牌
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRequest) (resp *types.RefreshTokenResponse, err error) {
	// 解析refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.Auth.AccessSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("无效的refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("无效的token claims")
	}

	// 获取用户ID
	userIdFloat, ok := claims["userId"].(float64)
	if !ok {
		return nil, fmt.Errorf("无效的用户ID")
	}
	userId := int64(userIdFloat)

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("无效的用户名")
	}

	// 生成新的access token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, userId, username)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 生成新的refresh token
	refreshToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, 7*24*3600, userId, username)
	if err != nil {
		return nil, fmt.Errorf("生成refresh token失败: %v", err)
	}

	resp = &types.RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	return resp, nil
}

func (l *RefreshTokenLogic) getJwtToken(secretKey string, iat, seconds, userId int64, username string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["username"] = username
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
