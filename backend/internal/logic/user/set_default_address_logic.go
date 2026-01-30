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

type SetDefaultAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设置默认地址
func NewSetDefaultAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultAddressLogic {
	return &SetDefaultAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDefaultAddressLogic) SetDefaultAddress(r *http.Request) (resp *types.SetDefaultAddressResponse, err error) {
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

	// 验证地址是否存在且属于当前用户
	_, err = l.svcCtx.Repos.AddressRepo.GetByIDAndUserID(l.ctx, addressId, userId)
	if err != nil {
		logx.Errorf("查询地址失败: %v", err)
		return nil, fmt.Errorf("地址不存在")
	}

	// 设置默认地址（内部会先清除其他默认地址）
	if err := l.svcCtx.Repos.AddressRepo.SetDefaultByIDAndUserID(l.ctx, addressId, userId); err != nil {
		logx.Errorf("设置默认地址失败: %v", err)
		return nil, fmt.Errorf("设置默认地址失败")
	}

	return &types.SetDefaultAddressResponse{
		Success: true,
	}, nil
}
