// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户资料
func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile() (resp *types.UserProfile, err error) {
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

	// TODO: 查询用户地址列表
	// 这里简化处理，返回空地址列表
	addresses := []types.Address{}

	resp = &types.UserProfile{
		Id:        fmt.Sprintf("%d", user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Birthday:  user.Birthday,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Addresses: addresses,
	}

	return resp, nil
}
