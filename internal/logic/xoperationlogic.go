package logic

import (
	"context"

	"x_operation/internal/svc"
	"x_operation/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type X_operationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewX_operationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *X_operationLogic {
	return &X_operationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *X_operationLogic) X_operation(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
