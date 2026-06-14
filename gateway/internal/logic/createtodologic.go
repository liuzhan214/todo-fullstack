package logic

import (
	"context"

	"todo/gateway/internal/svc"
	"todo/gateway/internal/types"
	"todo/gateway/pb"
	"todo/gateway/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTodoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTodoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTodoLogic {
	return &CreateTodoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTodoLogic) CreateTodo(req *types.CreateTodoReq) (resp *types.CreateTodoResp, err error) {
	l.Infof("CreateTodo request: %s", utils.FormatJson(req))

	rpcResp, err := l.svcCtx.TodoRpc.CreateTodo(l.ctx, &pb.CreateTodoReq{
		Title:         req.Title,
		IsOnlyForTest: req.IsOnlyForTest,
	})
	if err != nil {
		l.Errorf("CreateTodo error: %v", err)
		return nil, err
	}

	resp = &types.CreateTodoResp{
		Todo: types.Todo{
			Id:            rpcResp.Todo.Id,
			Title:         rpcResp.Todo.Title,
			Completed:     rpcResp.Todo.Completed,
			IsOnlyForTest: rpcResp.Todo.IsOnlyForTest,
		},
	}
	l.Infof("CreateTodo response: %s", utils.FormatJson(resp))
	return resp, nil
}
