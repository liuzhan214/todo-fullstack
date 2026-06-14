package svc

import (
	"todo/backend/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config    config.Config
	TodoRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := redis.MustNewRedis(c.TodoRedis)
	return &ServiceContext{
		Config:    c,
		TodoRedis: rds,
	}
}
