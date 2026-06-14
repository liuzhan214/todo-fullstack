# Todo App 实现计划

## Context
在 `/Users/liuzhan/full-stack/frontend/todo/` 下创建一个完整的 Todo 应用。采用 Go-Zero 微服务架构，HTTP 层做转发和聚合，业务逻辑在 RPC 层。所有服务通过 Docker Compose 编排部署。

## 部署架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                          宿主机 (macOS)                              │
│                                                                     │
│   浏览器访问 http://localhost:5173                                    │
│         │                                                           │
│         ▼                                                           │
│   ┌──────────────────────────────────────────────────────────────┐  │
│   │  Docker Compose 网络                                          │  │
│   │                                                              │  │
│   │  ┌──────────────┐                                            │  │
│   │  │ todo-frontend │                                            │  │
│   │  │ (Vite)       │                                            │  │
│   │  │ :5173        │                                            │  │
│   │  └──────┬───────┘                                            │  │
│   │         │ /api/*                                              │  │
│   │         ▼                                                     │  │
│   │  ┌──────────────┐   gRPC    ┌────────────────────────────┐  │  │
│   │  │ todo-gateway │ ────────▶ │ todo-backend               │  │  │
│   │  │ (Go-Zero     │           │ (Go-Zero zrpc)             │  │  │
│   │  │  HTTP)       │           │                            │  │  │
│   │  │              │           │ 业务逻辑层                   │  │  │
│   │  │ :20000       │           │ :20001                     │  │  │
│   │  │              │           │                            │  │  │
│   │  │ 只做转发和    │           │ GET /api/getTodos          │  │  │
│   │  │ 接口聚合      │           │ POST /api/createTodo       │  │  │
│   │  │              │           │ POST /api/updateTodo       │  │  │
│   │  │              │           │ POST /api/deleteTodo       │  │  │
│   │  └──────────────┘           └─────────────┬──────────────┘  │  │
│   │                                            │                 │  │
│   │                                            ▼                 │  │
│   │                                   ┌──────────────────────┐   │  │
│   │                                   │ redis-todo           │   │  │
│   │                                   │ (Redis 7.4)          │   │  │
│   │                                   │ :6379                │   │  │
│   │                                   │ 数据持久化 (volume)    │   │  │
│   │                                   └──────────────────────┘   │  │
│   └──────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
```

## 目录结构
```
frontend/todo/
├── docker-compose.yml       # 编排四个服务
├── plan.md
├── frontend/                # React 前端（待初始化）
│   ├── Dockerfile
├── gateway/                 # HTTP 网关（Go-Zero rest）
│   ├── Dockerfile
│   ├── todo.api
│   ├── todo.go
│   ├── utils/               # 本地工具包
│   ├── etc/
│   │   └── todo-api.yaml
│   └── internal/
│       ├── config/
│       ├── handler/
│       ├── logic/
│       └── svc/
├── backend/                 # RPC 服务（Go-Zero zrpc）
│   ├── Dockerfile
│   ├── todo.proto
│   ├── todo.go
│   ├── utils/               # 本地工具包
│   ├── etc/
│   │   └── todo.yaml
│   └── internal/
│       ├── config/
│       ├── server/
│       ├── logic/
│       └── svc/
└── integrated_tests/        # 集成测试
    ├── go.mod
    ├── testutil/
    └── cases/
```

## 技术选型
- **前端**: React + Vite（Node 20 Alpine）
- **网关**: Go-Zero rest（HTTP 层，只做转发和聚合）
- **后端**: Go-Zero zrpc（gRPC 服务，业务逻辑层）
- **存储**: Redis 7.4（数据持久化到 Docker Volume）
- **部署**: Docker Compose

## 接口契约

### HTTP 接口（todo.api）
| 方法 | 路径 | 请求体 | 功能 |
|------|------|--------|------|
| GET | /api/getTodos | 无 | 获取所有 todo |
| POST | /api/createTodo | `{ "title": "..." }` | 创建 todo |
| POST | /api/updateTodo | `{ "id": "...", "title": "...", "completed": true }` | 更新 todo |
| POST | /api/deleteTodo | `{ "id": "..." }` | 删除 todo |

### gRPC 接口（todo.proto）
| 方法 | 请求 | 响应 | 功能 |
|------|------|------|------|
| GetTodos | GetTodosReq | GetTodosResp{todos} | 获取所有 |
| CreateTodo | CreateTodoReq{title} | CreateTodoResp{todo} | 创建 |
| UpdateTodo | UpdateTodoReq{id,title,completed} | UpdateTodoResp | 更新 |
| DeleteTodo | DeleteTodoReq{id} | DeleteTodoResp | 删除 |

## Docker 服务说明

| 服务 | 镜像 | 端口 | 依赖 | 说明 |
|------|------|------|------|------|
| redis | redis:7.4 | 6379 | 无 | 数据存储 |
| backend | ./backend | 20001(gRPC) | redis | 业务逻辑 |
| gateway | ./gateway | 20000(HTTP) | backend | HTTP 转发 |
| frontend | ./frontend | 5173 | gateway | 前端界面 |

## 实现步骤

### 第一步：安装 goctl 工具
```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

### 第二步：用 goctl 生成后端代码骨架
```bash
# 生成 RPC 服务代码
cd /Users/liuzhan/full-stack/frontend/todo/backend 
goctl rpc protoc todo.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.

# 生成 Gateway 代码
cd /Users/liuzhan/full-stack/frontend/todo/gateway 
goctl api go --api todo.api --dir .
```

```bash
# 生成gateway下游依赖backend的代码
cd /Users/liuzhan/full-stack/frontend/todo/backend
goctl rpc protoc todo.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.
```

### 第三步：填充业务逻辑
- backend/internal/logic/ — 实现 Redis 操作
- gateway/internal/logic/ — 调用 RPC 方法

### 第四步：创建 React 前端
- npm create vite@latest 初始化
- vite.config.js 配置代理
- App.jsx 实现 Todo 界面

### 第五步：启动和测试
```bash
# 构建并启动
cd /Users/liuzhan/full-stack/frontend/todo
docker compose up --build -d

# 测试 API
curl http://127.0.0.1:20000/api/getTodos
curl -X POST http://127.0.0.1:20000/api/createTodo -H "Content-Type: application/json" -d '{"title":"test"}'
curl -X POST http://127.0.0.1:20000/api/updateTodo -H "Content-Type: application/json" -d '{"id":"1","title":"updated","completed":true}'
curl -X POST http://127.0.0.1:20000/api/deleteTodo -H "Content-Type: application/json" -d '{"id":"1"}'
```

访问 http://localhost:5173

## 命令速查

### Go-Zero 代码生成
```bash
# RPC 服务 — 从 proto 生成 server/client 代码
goctl rpc protoc todo.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.

# Gateway — 从 api 定义生成 handler/logic 骨架
goctl api go --api todo.api --dir .

# 整理依赖
go mod tidy
```

### Docker 操作
```bash
# 构建镜像并启动
docker compose up --build -d

# 停止并移除容器
docker compose down

# 查看服务日志
docker compose logs <service>
docker compose logs --tail=50 <service>

# 查看运行状态
docker compose ps

# 进入容器
docker exec -it <container> sh

# 只重建特定服务
docker compose up --build -d <service>
```

## Config 结构体注意事项

### RpcServerConf 字段冲突
```go
type Config struct {
    zrpc.RpcServerConf              // 匿名嵌入，内部有 Redis redis.RedisKeyConf
    TodoRedis redis.RedisConf       // ❌ 不能叫 Redis，会和 RpcServerConf.Redis 冲突
}
```

**规则**: `RpcServerConf` 匿名嵌入后，其所有子字段提升到顶层。如果自定义字段和提升的字段同名，Go-Zero 加载配置时报 `conflict key` 错误。必须重命名自定义字段。

### RedisConf Key 字段
```go
type RedisConf struct {
    Host string
    Type string       // node / cluster
    Key  string       // Redis 密码，为空表示无密码
    ...
}
```

Go-Zero 的配置验证器要求所有字段都有显式值，即使不需要密码也要设 `Key: ""`。

## 踩坑记录

| # | 问题 | 错误现象 | 原因 | 解决方案 |
|---|------|----------|------|----------|
| 1 | Gateway 端口不通 | `curl: exit code 7`，容器显示 "Up" 但进程持续重启 | 配置文件 `etc/todo-api.yaml` 没复制到 Docker 最终镜像，`conf.MustLoad` 用默认值启动后立即退出 | Dockerfile 加 `COPY --from=builder /app/etc /etc` |
| 2 | GOPROXY 网络超时 | `go mod download` / `go build` 超时 120s | 国内网络连不上 Docker Hub / Go 官方 proxy | Dockerfile 加 `ENV GOPROXY=https://goproxy.cn,direct` |
| 3 | Config 字段冲突 | `conflict key redis, pay attention to anonymous fields` | `Config` 嵌入 `RpcServerConf`（含 `Redis` 字段），又自定义 `Redis redis.RedisConf` | 自定义字段改名：`Redis` → `TodoRedis` |
| 4 | Redis 连接被拒绝 | `dial tcp 172.19.0.x:6380: connect: connection refused` | 端口映射 `6380:6380`，但 Redis 容器内默认监听 6379，不是 6380 | 映射改为 `6380:6379`，backend 连 `redis:6379` |
| 5 | Docker 构建缓存 | 代码改了但容器行为不变 | `docker compose up` 复用缓存层 | 重建时加 `--build` 或清理 `docker builder prune` |
| 6 | 镜像版本不存在 | `golang:1.26-alpine: failed to resolve source metadata` | Docker Hub 上不可用指定 Go 镜像 tag | 确认可用 tag 后再写 |
| 7 | 配置文件未更新 | `docker compose restart` 后配置没变 | Docker 镜像构建时打包配置，重启不重建 | 改配置后必须 `docker compose up --build -d` |

## 网络通信流程
```
用户浏览器
  │
  │ http://localhost:5173
  ▼
todo-frontend (Vite 开发服务器)
  │
  │ /api/* → vite proxy → http://gateway:20000
  ▼
todo-gateway (Go-Zero HTTP)
  │
  │ gRPC → backend:20001
  ▼
todo-backend (Go-Zero zrpc)
  │
  │ redis:6379
  ▼
redis-todo (Redis)
```