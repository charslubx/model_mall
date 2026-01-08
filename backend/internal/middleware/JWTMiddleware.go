package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type JwtMiddleware struct {
	Secret      string
	RedisClient *redis.Redis
}

func NewJwtMiddleware(secret string) *JwtMiddleware {
	return &JwtMiddleware{
		Secret:      secret,
		RedisClient: nil, // 将在ServiceContext中设置
	}
}

func NewJwtMiddlewareWithRedis(secret string, redisClient *redis.Redis) *JwtMiddleware {
	return &JwtMiddleware{
		Secret:      secret,
		RedisClient: redisClient,
	}
}

// UserClaims 从 JWT token 中提取的用户信息
type UserClaims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #region agent log
		func() {
			f, _ := os.OpenFile("/home/model_mall/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if f != nil {
				defer f.Close()
				data, _ := json.Marshal(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "E", "location": "JWTMiddleware.go:30", "message": "Custom JWTMiddleware called", "data": map[string]interface{}{"path": r.URL.Path, "method": r.Method}, "timestamp": time.Now().UnixMilli()})
				f.Write(append(data, '\n'))
			}
		}()
		// #endregion
		// 1. 获取 token
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			// #region agent log
			func() {
				f, _ := os.OpenFile("/home/model_mall/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if f != nil {
					defer f.Close()
					data, _ := json.Marshal(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "E", "location": "JWTMiddleware.go:36", "message": "Custom middleware: No authorization header", "data": map[string]interface{}{}, "timestamp": time.Now().UnixMilli()})
					f.Write(append(data, '\n'))
				}
			}()
			// #endregion
			Error(w, http.StatusUnauthorized, "未授权访问")
			return
		}

		// 2. 解析 Bearer token
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			Error(w, http.StatusUnauthorized, "无效的认证格式")
			return
		}

		tokenString := parts[1]

		// 3. 检查token是否在黑名单中
		if m.RedisClient != nil {
			blacklisted, err := m.RedisClient.GetCtx(r.Context(), "token_blacklist:"+tokenString)
			if err == nil && blacklisted == "1" {
				logx.Info("Token已在黑名单中")
				Error(w, http.StatusUnauthorized, "token已失效，请重新登录")
				return
			}
		}

		// 4. 验证 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.Secret), nil
		})

		if err != nil {
			logx.Errorf("Token验证失败: %v", err)
			Error(w, http.StatusUnauthorized, "无效的token")
			return
		}

		// 5. 获取 claims
		if !token.Valid {
			Error(w, http.StatusUnauthorized, "无效的token")
			return
		}

		// 6. 提取用户ID
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			Error(w, http.StatusUnauthorized, "无效的token格式")
			return
		}

		var userId int64
		if userIdFloat, ok := claims["userId"].(float64); ok {
			userId = int64(userIdFloat)
		} else {
			Error(w, http.StatusUnauthorized, "token中缺少userId")
			return
		}

		// 7. 将用户信息和token注入到请求上下文
		ctx := context.WithValue(r.Context(), "userId", userId)
		ctx = context.WithValue(ctx, "token", tokenString)
		r = r.WithContext(ctx)

		// 8. 继续处理请求
		next(w, r)
	}
}
