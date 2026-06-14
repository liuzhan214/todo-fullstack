package logic

import (
	"context"
	"encoding/json"

	"todo/backend/utils"
	"todo/backend/internal/svc"
	"todo/backend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetTodosLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTodosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTodosLogic {
	return &GetTodosLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTodosLogic) GetTodos(in *pb.GetTodosReq) (*pb.GetTodosResp, error) {
	l.Infof("GetTodos request: %s", utils.FormatJson(in))

	vals, err := l.svcCtx.TodoRedis.Hvals(todoKey)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		l.Errorf("GetTodos error: %v", err)
		return nil, err
	}

	var todos []*pb.Todo
	for _, v := range vals {
		var todo pb.Todo
		if err := json.Unmarshal([]byte(v), &todo); err != nil {
			l.Errorf("GetTodos unmarshal error: %v", err)
			continue
		}
		todos = append(todos, &todo)
	}

	resp := &pb.GetTodosResp{Todos: todos}
	l.Infof("GetTodos response: %s", utils.FormatJson(resp))
	return resp, nil
}
