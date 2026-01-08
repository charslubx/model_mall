// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"
	"model_mall_backend/backend/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 根据邮箱查询用户
	user, err := l.svcCtx.Repos.UserRepo.GetByEmail(l.ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("用户不存在或密码错误")
	}

	// 检查账户状态
	if user.Status == 0 {
		return nil, fmt.Errorf("账户已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("用户不存在或密码错误")
	}

	// 获取用户类型
	userType := user.UserType
	if userType == "" {
		userType = getUserType(user.RoleID)
	}

	// 生成JWT token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 生成refresh token (7天有效期)
	refreshToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, 7*24*3600, user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成refresh token失败: %v", err)
	}

	// 更新最后登录时间
	_ = l.svcCtx.Repos.UserRepo.UpdateLastLogin(l.ctx, user.ID, getClientIP(l.ctx))

	// 构造响应
	resp = &types.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User: types.UserInfo{
			Id:        fmt.Sprintf("%d", user.ID),
			Name:      user.Nickname,
			Email:     user.Email,
			Avatar:    user.Avatar,
			UserType:  userType,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	return resp, nil
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64, username string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["username"] = username
	// #region agent log
	func() {
		f, _ := os.OpenFile("/home/model_mall/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if f != nil {
			defer f.Close()
			data, _ := json.Marshal(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "B,D", "location": "login_logic.go:96", "message": "Generating JWT token", "data": map[string]interface{}{"secretKeyPrefix": secretKey[:20], "userId": userId, "username": username, "claims": claims}, "timestamp": time.Now().UnixMilli()})
			f.Write(append(data, '\n'))
		}
	}()
	// #endregion
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	signedToken, err := token.SignedString([]byte(secretKey))
	// #region agent log
	func() {
		f, _ := os.OpenFile("/home/model_mall/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if f != nil {
			defer f.Close()
			data, _ := json.Marshal(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "B", "location": "login_logic.go:103", "message": "JWT token signed", "data": map[string]interface{}{"tokenPrefix": signedToken[:50], "error": err}, "timestamp": time.Now().UnixMilli()})
			f.Write(append(data, '\n'))
		}
	}()
	// #endregion
	return signedToken, err
}

func getUserType(roleID int64) string {
	// 1: customer, 2: merchant, 3: admin
	switch roleID {
	case 1:
		return "customer"
	case 2:
		return "merchant"
	case 3:
		return "admin"
	default:
		return "customer"
	}
}

func getClientIP(ctx context.Context) string {
	// 从context中获取客户端IP
	// 这里简化处理,实际应该从HTTP请求头中获取
	return "127.0.0.1"
}
