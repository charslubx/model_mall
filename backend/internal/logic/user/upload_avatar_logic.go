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

type UploadAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 上传用户头像
func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadAvatarLogic) UploadAvatar() (resp *types.UploadAvatarResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// TODO: 实际的文件上传逻辑
	// 这里模拟生成头像URL
	avatarUrl := fmt.Sprintf("https://cdn.example.com/avatars/u%d_%d.jpg", userId, time.Now().Unix())

	// 更新用户头像
	user, err := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	user.Avatar = avatarUrl
	err = l.svcCtx.Repos.UserRepo.Update(l.ctx, user)
	if err != nil {
		return nil, fmt.Errorf("更新头像失败: %v", err)
	}

	logx.Infof("用户 %d 上传头像成功: %s", userId, avatarUrl)

	resp = &types.UploadAvatarResponse{
		AvatarUrl: avatarUrl,
	}

	return resp, nil
}
