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

type DeleteAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除收货地址
func NewDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddressLogic {
	return &DeleteAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAddressLogic) DeleteAddress(r *http.Request) (resp *types.BaseResponse, err error) {
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

	// 删除地址（确保只能删除自己的地址）
	if err := l.svcCtx.Repos.AddressRepo.DeleteByIDAndUserID(l.ctx, addressId, userId); err != nil {
		logx.Errorf("删除地址失败: %v", err)
		return nil, fmt.Errorf("删除地址失败")
	}

	return &types.BaseResponse{
		Code:    200,
		Message: "删除成功",
	}, nil
}
