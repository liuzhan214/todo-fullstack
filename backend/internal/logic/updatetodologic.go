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

type UpdateTodoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTodoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTodoLogic {
	return &UpdateTodoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTodoLogic) UpdateTodo(in *pb.UpdateTodoReq) (*pb.UpdateTodoResp, error) {
	l.Infof("UpdateTodo request: %s", utils.FormatJson(in))

	if in.Id == "" {
		err := status.Error(codes.InvalidArgument, "id is required")
		l.Errorf("UpdateTodo error: %v", err)
		return nil, err
	}

	// 读取旧数据，保留 is_only_for_test 标记
	oldRaw, err := l.svcCtx.TodoRedis.Hget(todoKey, in.Id)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		l.Errorf("UpdateTodo read old error: %v", err)
		return nil, err
	}
	oldIsOnlyForTest := false
	if oldRaw != "" {
		var oldTodo pb.Todo
		if err := json.Unmarshal([]byte(oldRaw), &oldTodo); err == nil {
			oldIsOnlyForTest = oldTodo.IsOnlyForTest
		}
	}

	todo := &pb.Todo{
		Id:            in.Id,
		Title:         in.Title,
		Completed:     in.Completed,
		IsOnlyForTest: oldIsOnlyForTest,
	}

	data, _ := json.Marshal(todo)
	if err := l.svcCtx.TodoRedis.Hset(todoKey, in.Id, string(data)); err != nil {
		err = status.Error(codes.Internal, err.Error())
		l.Errorf("UpdateTodo error: %v", err)
		return nil, err
	}

	resp := &pb.UpdateTodoResp{}
	l.Infof("UpdateTodo response: %s", utils.FormatJson(resp))
	return resp, nil
}
