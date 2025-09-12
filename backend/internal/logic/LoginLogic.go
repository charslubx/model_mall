package logic

import (
	"context"
	"errors"
	"time"

	"model_mall_backend/backend/internal/middleware"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
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
	// 验证用户名密码
	user, err := l.svcCtx.Repos.User.GetByUsername(l.ctx, req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录信息
	if err := l.svcCtx.Repos.User.UpdateLastLogin(l.ctx, user.ID, ""); err != nil {
		l.Logger.Errorf("更新最后登录信息失败: %v", err)
	}

	// 生成 JWT token
	token, err := l.generateToken(user.ID)
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
