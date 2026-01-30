// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新收货地址
func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAddressLogic) UpdateAddress(r *http.Request, req *types.UpdateAddressRequest) (resp *types.UpdateAddressResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 获取地址ID
	addressIdStr := r.URL.Query().Get(":id")
	if addressIdStr == "" {
		addressIdStr = r.PathValue("id")
	}

	addressId, err := strconv.ParseInt(addressIdStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的地址ID")
	}

	// 查询地址是否存在且属于当前用户
	address, err := l.svcCtx.Repos.AddressRepo.GetByIDAndUserID(l.ctx, addressId, userId)
	if err != nil {
		logx.Errorf("查询地址失败: %v", err)
		return nil, fmt.Errorf("地址不存在")
	}

	// 如果要设置为默认地址，先清除该用户的其他默认地址
	if req.IsDefault && !address.IsDefault {
		if err := l.svcCtx.Repos.AddressRepo.ClearDefaultByUserID(l.ctx, userId); err != nil {
			logx.Errorf("清除默认地址失败: %v", err)
			return nil, fmt.Errorf("操作失败")
		}
	}

	// 更新地址信息
	if req.Name != "" {
		address.Name = req.Name
	}
	if req.Phone != "" {
		address.Phone = req.Phone
	}
	if req.Province != "" {
		address.Province = req.Province
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.District != "" {
		address.District = req.District
	}
	if req.Detail != "" {
		address.Address = req.Detail
	}
	address.IsDefault = req.IsDefault
	address.UpdatedAt = time.Now()

	if err := l.svcCtx.Repos.AddressRepo.Update(l.ctx, address); err != nil {
		logx.Errorf("更新地址失败: %v", err)
		return nil, fmt.Errorf("更新地址失败")
	}

	return &types.UpdateAddressResponse{
		Id:        strconv.FormatInt(address.ID, 10),
		UpdatedAt: address.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
