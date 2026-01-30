// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAddressDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取地址详情
func NewGetAddressDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressDetailLogic {
	return &GetAddressDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAddressDetailLogic) GetAddressDetail(r *http.Request) (resp *types.GetAddressDetailResponse, err error) {
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

	// 查询地址信息（确保只能查询自己的地址）
	address, err := l.svcCtx.Repos.AddressRepo.GetByIDAndUserID(l.ctx, addressId, userId)
	if err != nil {
		logx.Errorf("查询地址失败: %v", err)
		return nil, fmt.Errorf("地址不存在")
	}

	return &types.GetAddressDetailResponse{
		Address: types.AddressInfo{
			Id:        strconv.FormatInt(address.ID, 10),
			Name:      address.Name,
			Phone:     address.Phone,
			Province:  address.Province,
			City:      address.City,
			District:  address.District,
			Detail:    address.Address,
			IsDefault: address.IsDefault,
		},
	}, nil
}
