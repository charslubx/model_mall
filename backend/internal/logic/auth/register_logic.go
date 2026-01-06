// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"
	"time"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"
	"model_mall_backend/backend/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.LoginResponse, err error) {
	// 检查邮箱是否已存在
	exists, err := l.svcCtx.Repos.UserRepo.ExistsByEmail(l.ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("邮箱已被注册")
	}

	// 如果提供了手机号,检查是否已存在
	if req.Phone != "" {
		exists, err = l.svcCtx.Repos.UserRepo.ExistsByPhone(l.ctx, req.Phone)
		if err != nil {
			return nil, fmt.Errorf("检查手机号失败: %v", err)
		}
		if exists {
			return nil, fmt.Errorf("手机号已被注册")
		}
	}

	// 如果是商户注册,验证必填字段
	if req.UserType == "merchant" && req.MerchantName == "" {
		return nil, fmt.Errorf("商户名称不能为空")
	}

	// 对密码进行哈希加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 确定角色ID
	roleID := getRoleIDByType(req.UserType)

	// 创建用户
	user := &models.User{
		Username: req.Email, // 使用邮箱作为用户名
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: req.Name,
		Status:   1, // 默认启用
		RoleID:   roleID,
	}

	err = l.svcCtx.Repos.UserRepo.Create(l.ctx, user)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	// 生成JWT token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 生成refresh token
	refreshToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, 7*24*3600, user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成refresh token失败: %v", err)
	}

	// 构造响应
	resp = &types.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User: types.UserInfo{
			Id:        fmt.Sprintf("%d", user.ID),
			Name:      user.Nickname,
			Email:     user.Email,
			Avatar:    user.Avatar,
			UserType:  req.UserType,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	return resp, nil
}

func (l *RegisterLogic) getJwtToken(secretKey string, iat, seconds, userId int64, username string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["username"] = username
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func getRoleIDByType(userType string) int64 {
	// 1: customer, 2: merchant, 3: admin
	switch userType {
	case "customer":
		return 1
	case "merchant":
		return 2
	case "admin":
		return 3
	default:
		return 1
	}
}
