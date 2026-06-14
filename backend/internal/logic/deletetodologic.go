package logic

import (
	"context"

	"todo/backend/utils"
	"todo/backend/internal/svc"
	"todo/backend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteTodoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTodoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTodoLogic {
	return &DeleteTodoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteTodoLogic) DeleteTodo(in *pb.DeleteTodoReq) (*pb.DeleteTodoResp, error) {
	l.Infof("DeleteTodo request: %s", utils.FormatJson(in))

	if in.Id == "" {
		err := status.Error(codes.InvalidArgument, "id is required")
		l.Errorf("DeleteTodo error: %v", err)
		return nil, err
	}

	if _, err := l.svcCtx.TodoRedis.Hdel(todoKey, in.Id); err != nil {
		err = status.Error(codes.Internal, err.Error())
		l.Errorf("DeleteTodo error: %v", err)
		return nil, err
	}

	resp := &pb.DeleteTodoResp{}
	l.Infof("DeleteTodo response: %s", utils.FormatJson(resp))
	return resp, nil
}
