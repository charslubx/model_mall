// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"
	"model_mall_backend/backend/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改密码
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest) (resp *types.BaseResponse, err error) {
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

	// 验证旧密码
	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		return nil, fmt.Errorf("旧密码错误")
	}

	// 验证新密码长度
	if len(req.NewPassword) < 8 {
		return nil, fmt.Errorf("新密码长度至少8位")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码
	user.Password = hashedPassword
	err = l.svcCtx.Repos.UserRepo.Update(l.ctx, user)
	if err != nil {
		return nil, fmt.Errorf("修改密码失败: %v", err)
	}

	logx.Infof("用户 %d 修改密码成功", userId)

	resp = &types.BaseResponse{
		Code:    200,
		Message: "修改成功",
		Data:    map[string]bool{"success": true},
	}

	return resp, nil
}
