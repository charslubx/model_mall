// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAddressesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户地址列表
func NewGetAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressesLogic {
	return &GetAddressesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAddressesLogic) GetAddresses() (resp *types.GetAddressesResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 查询用户的所有地址
	addresses, err := l.svcCtx.Repos.AddressRepo.GetByUserID(l.ctx, userId)
	if err != nil {
		logx.Errorf("查询用户地址失败: %v", err)
		return nil, fmt.Errorf("查询地址失败")
	}

	// 转换为响应格式
	addressList := make([]types.AddressInfo, 0, len(addresses))
	for _, addr := range addresses {
		addressList = append(addressList, types.AddressInfo{
			Id:        strconv.FormatInt(addr.ID, 10),
			Name:      addr.Name,
			Phone:     addr.Phone,
			Province:  addr.Province,
			City:      addr.City,
			District:  addr.District,
			Detail:    addr.Address,
			IsDefault: addr.IsDefault,
		})
	}

	return &types.GetAddressesResponse{
		Addresses: addressList,
		Total:     len(addressList),
	}, nil
}
