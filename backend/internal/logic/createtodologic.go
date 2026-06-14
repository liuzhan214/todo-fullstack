package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"todo/backend/utils"
	"todo/backend/internal/svc"
	"todo/backend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const todoKey = "todos"

type CreateTodoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTodoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTodoLogic {
	return &CreateTodoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTodoLogic) CreateTodo(in *pb.CreateTodoReq) (*pb.CreateTodoResp, error) {
	l.Infof("CreateTodo request: %s", utils.FormatJson(in))

	if in.Title == "" {
		err := status.Error(codes.InvalidArgument, "title is required")
		l.Errorf("CreateTodo error: %v", err)
		return nil, err
	}

	n, err := l.svcCtx.TodoRedis.Incr("todo_id")
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		l.Errorf("CreateTodo error: %v", err)
		return nil, err
	}
	id := fmt.Sprintf("%d", n)
	todo := &pb.Todo{
		Id:            id,
		Title:         in.Title,
		IsOnlyForTest: in.IsOnlyForTest,
	}

	data, _ := json.Marshal(todo)
	if err := l.svcCtx.TodoRedis.Hset(todoKey, id, string(data)); err != nil {
		err = status.Error(codes.Internal, err.Error())
		l.Errorf("CreateTodo error: %v", err)
		return nil, err
	}

	resp := &pb.CreateTodoResp{Todo: todo}
	l.Infof("CreateTodo response: %s", utils.FormatJson(resp))
	return resp, nil
}
