package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"todo/gateway/internal/svc"
	"todo/gateway/internal/types"
	"todo/gateway/pb"
	"todo/gateway/utils"
)

type UpdateTodoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTodoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTodoLogic {
	return &UpdateTodoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTodoLogic) UpdateTodo(req *types.UpdateTodoReq) (resp *types.UpdateTodoResp, err error) {
	l.Infof("UpdateTodo request: %s", utils.FormatJson(req))

	_, err = l.svcCtx.TodoRpc.UpdateTodo(l.ctx, &pb.UpdateTodoReq{
		Id:        req.Id,
		Title:     req.Title,
		Completed: req.Completed,
	})
	if err != nil {
		l.Errorf("UpdateTodo error: %v", err)
		return nil, err
	}

	resp = &types.UpdateTodoResp{}
	l.Infof("UpdateTodo response: %s", utils.FormatJson(resp))
	return resp, nil
}
