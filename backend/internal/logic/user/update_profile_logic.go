// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"
	"time"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户资料
func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileRequest) (resp *types.UpdateProfileResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 查询用户信息
	user, err := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 更新字段
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// 检查邮箱是否已被使用
		existing, _ := l.svcCtx.Repos.UserRepo.GetByEmail(l.ctx, req.Email)
		if existing != nil && existing.ID != userId {
			return nil, fmt.Errorf("邮箱已被使用")
		}
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Gender != "" {
		switch req.Gender {
		case "male":
			user.Gender = 1
		case "female":
			user.Gender = 2
		default:
			user.Gender = 0
		}
	}
	if req.Birthday != "" {
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			user.Birthday = &birthday
		}
	}

	// 保存更新
	err = l.svcCtx.Repos.UserRepo.Update(l.ctx, user)
	if err != nil {
		return nil, fmt.Errorf("更新用户资料失败: %v", err)
	}

	logx.Infof("用户 %d 更新资料成功", userId)

	resp = &types.UpdateProfileResponse{
		Success: true,
		User: types.UserInfo{
			Id:    fmt.Sprintf("%d", user.ID),
			Name:  user.Name,
			Phone: user.Phone,
		},
	}

	return resp, nil
}
