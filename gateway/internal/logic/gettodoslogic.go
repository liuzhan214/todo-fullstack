package logic

import (
	"context"

	"todo/gateway/internal/svc"
	"todo/gateway/internal/types"
	"todo/gateway/pb"
	"todo/gateway/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTodosLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTodosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTodosLogic {
	return &GetTodosLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTodosLogic) GetTodos(req *types.GetTodosReq) (resp *types.GetTodosResp, err error) {
	l.Infof("GetTodos request: %s", utils.FormatJson(req))

	rpcResp, err := l.svcCtx.TodoRpc.GetTodos(l.ctx, &pb.GetTodosReq{})
	if err != nil {
		l.Errorf("GetTodos error: %v", err)
		return nil, err
	}

	var todos []types.Todo
	for _, t := range rpcResp.Todos {
		todos = append(todos, types.Todo{
			Id:            t.Id,
			Title:         t.Title,
			Completed:     t.Completed,
			IsOnlyForTest: t.IsOnlyForTest,
		})
	}

	resp = &types.GetTodosResp{Todos: todos}
	l.Infof("GetTodos response: %s", utils.FormatJson(resp))
	return resp, nil
}
