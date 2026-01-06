// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package merchant

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMerchantProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新商户信息
func NewUpdateMerchantProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMerchantProfileLogic {
	return &UpdateMerchantProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMerchantProfileLogic) UpdateMerchantProfile(req *types.UpdateMerchantProfileRequest) (resp *types.BaseResponse, err error) {
	// 获取商户ID
	merchantId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 查询商户信息
	merchant, err := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, merchantId)
	if err != nil {
		return nil, fmt.Errorf("商户不存在")
	}

	// 更新字段
	if req.Name != "" {
		merchant.MerchantName = req.Name
	}
	if req.Description != "" {
		merchant.Description = req.Description
	}
	if req.Avatar != "" {
		merchant.Avatar = req.Avatar
	}

	// 保存更新
	err = l.svcCtx.Repos.UserRepo.Update(l.ctx, merchant)
	if err != nil {
		return nil, fmt.Errorf("更新商户信息失败: %v", err)
	}

	logx.Infof("商户 %d 更新信息成功", merchantId)

	resp = &types.BaseResponse{
		Code:    200,
		Message: "更新成功",
		Data:    map[string]bool{"success": true},
	}

	return resp, nil
}
