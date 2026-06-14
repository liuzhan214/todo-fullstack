package logic

import (
	"context"

	"todo/gateway/internal/svc"
	"todo/gateway/internal/types"
	"todo/gateway/pb"
	"todo/gateway/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTodoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTodoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTodoLogic {
	return &DeleteTodoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTodoLogic) DeleteTodo(req *types.DeleteTodoReq) (resp *types.DeleteTodoResp, err error) {
	l.Infof("DeleteTodo request: %s", utils.FormatJson(req))

	_, err = l.svcCtx.TodoRpc.DeleteTodo(l.ctx, &pb.DeleteTodoReq{
		Id: req.Id,
	})
	if err != nil {
		l.Errorf("DeleteTodo error: %v", err)
		return nil, err
	}

	resp = &types.DeleteTodoResp{}
	l.Infof("DeleteTodo response: %s", utils.FormatJson(resp))
	return resp, nil
}
