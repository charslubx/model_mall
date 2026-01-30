// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加收货地址
func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAddressLogic) AddAddress(req *types.AddAddressRequest) (resp *types.AddAddressResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 如果设置为默认地址，先清除该用户的其他默认地址
	if req.IsDefault {
		if err := l.svcCtx.Repos.AddressRepo.ClearDefaultByUserID(l.ctx, userId); err != nil {
			logx.Errorf("清除默认地址失败: %v", err)
			return nil, fmt.Errorf("操作失败")
		}
	}

	// 创建地址记录
	now := time.Now()
	address := &models.Address{
		UserID:    userId,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Address:   req.Detail,
		IsDefault: req.IsDefault,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := l.svcCtx.Repos.AddressRepo.Create(l.ctx, address); err != nil {
		logx.Errorf("创建地址失败: %v", err)
		return nil, fmt.Errorf("添加地址失败")
	}

	return &types.AddAddressResponse{
		Id:        strconv.FormatInt(address.ID, 10),
		CreatedAt: address.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
