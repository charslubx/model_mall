package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type JwtMiddleware struct {
	Secret string
}

func NewJwtMiddleware(secret string) *JwtMiddleware {
	return &JwtMiddleware{
		Secret: secret,
	}
}

// UserClaims 从 JWT token 中提取的用户信息
type UserClaims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 获取 token
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			Error(w, http.StatusUnauthorized, "未授权访问")
			return
		}

		// 2. 解析 Bearer token
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			Error(w, http.StatusUnauthorized, "无效的认证格式")
			return
		}

		// 3. 验证 token
		token, err := jwt.ParseWithClaims(parts[1], &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.Secret), nil
		})

		if err != nil {
			logx.Errorf("Token验证失败: %v", err)
			Error(w, http.StatusUnauthorized, "无效的token")
			return
		}

		// 4. 获取 claims
		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			Error(w, http.StatusUnauthorized, "无效的token")
			return
		}

		// 5. 将用户信息注入到请求上下文
		ctx := context.WithValue(r.Context(), "userId", claims.UserId)
		r = r.WithContext(ctx)

		// 6. 继续处理请求
		next(w, r)
	}
}

// GetUserIdFromCtx 从上下文中获取用户ID
func GetUserIdFromCtx(ctx context.Context) (int64, bool) {
	userId, ok := ctx.Value("userId").(int64)
	return userId, ok
}
