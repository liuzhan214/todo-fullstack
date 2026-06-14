// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"todo/gateway/internal/config"
	"todo/gateway/pb"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	TodoRpc pb.TodoServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := zrpc.MustNewClient(c.Backend)
	return &ServiceContext{
		Config:  c,
		TodoRpc: pb.NewTodoServiceClient(conn.Conn()),
	}
}
