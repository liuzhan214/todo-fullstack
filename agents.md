# Goctl 开发指南

## 前置条件

```bash
# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 安装 protoc 及 Go 插件
# protoc 需手动安装（brew install protobuf 或官网下载）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

确认工具链：

```bash
protoc --version           # 需要 v3+
protoc-gen-go --version    # 需要 v1.30+
protoc-gen-go-grpc --version
goctl --version
```

---

## 项目结构

```
todo/
├── backend/                 # RPC 服务
│   ├── todo.proto           # proto 定义（唯一源）
│   ├── todo.go              # 服务入口
│   ├── pb/pb/               # 由 goctl rpc protoc 生成
│   └── internal/
│       ├── config/
│       ├── server/
│       ├── logic/           # 业务逻辑（手动维护）
│       └── svc/
├── gateway/                 # HTTP 网关
│   ├── todo.api             # API 定义
│   ├── todo.go              # 网关入口
│   ├── pb/                  # 从 backend/todo.proto 用 protoc 生成
│   └── internal/
│       ├── config/
│       ├── handler/         # 由 goctl api go 生成
│       ├── types/           # 由 goctl api go 生成
│       ├── logic/           # 业务逻辑（手动维护）
│       └── svc/
└── integrated_tests/
```

---

## 新增字段：完整工作流

以本次添加 `is_only_for_test` 字段为例。

### 第 1 步：修改 proto

```bash
# 编辑 backend/todo.proto
message Todo {
    string id = 1;
    string title = 2;
    bool completed = 3;
    bool is_only_for_test = 4;   # ← 新增
}

message CreateTodoReq {
    string title = 1;
    bool is_only_for_test = 2;   # ← 新增
}
```

### 第 2 步：再生 backend pb 代码

```bash
cd /Users/liuzhan/full-stack/frontend/todo/backend
goctl rpc protoc todo.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.
```

作用：
- `--go_out=./pb` — 生成 `pb/pb/todo.pb.go`（Go struct + 序列化）
- `--go-grpc_out=./pb` — 生成 `pb/pb/todo_grpc.pb.go`（gRPC client/server 接口）
- `--zrpc_out=.` — 生成 `internal/server/todoserviceserver.go`（Go-Zero 服务注册）
- 已存在的 logic 文件不会被覆盖

涉及改动文件：
- `pb/pb/todo.pb.go` ✅ 由 protoc-gen-go 生成
- `pb/pb/todo_grpc.pb.go` ✅ 由 protoc-gen-go-grpc 生成
- `internal/server/todoserviceserver.go` ✅ 由 goctl 生成

### 第 3 步：再生 gateway pb 代码（从 backend proto）

```bash
cd /Users/liuzhan/full-stack/frontend/todo/gateway

protoc \
  -I/Users/liuzhan/full-stack/frontend/todo/backend \
  --go_out=./pb \
  --go-grpc_out=./pb \
  /Users/liuzhan/full-stack/frontend/todo/backend/todo.proto
```

作用：直接从 backend 的 proto 文件生成 pb，不复制 proto，不依赖 gateway 自己的 proto。

### 第 4 步：更新 todo.api

```bash
# 编辑 gateway/todo.api
type Todo {
    Id            string `json:"id"`
    Title         string `json:"title"`
    Completed     bool   `json:"completed"`
    IsOnlyForTest bool   `json:"is_only_for_test"`  # ← 新增
}

type CreateTodoReq {
    Title         string `json:"title"`
    IsOnlyForTest bool   `json:"is_only_for_test"`  # ← 新增
}
```

### 第 5 步：再生 gateway types + handler

```bash
cd /Users/liuzhan/full-stack/frontend/todo/gateway
goctl api go --api todo.api --dir .
```

作用：
- `types/types.go` — 生成 Go 类型（新增字段自动包含）
- `handler/*handler.go` — 生成 handler 方法
- 已存在的 `logic/*.go` 不会被覆盖（会提示 `exists, ignored generation`）
- 已存在的 `routes.go`、`config.go`、`svc/servicecontext.go` 也不会被覆盖

**注意**：如果 goctl 生成的 `types.go` 不包含新字段（如旧版本 goctl 的 bug），需要手动补。

### 第 6 步：修改业务逻辑

#### Backend：CreateTodoLogic

```go
todo := &pb.Todo{
    Id:            id,
    Title:         in.Title,
    IsOnlyForTest: in.IsOnlyForTest,  // 透传
}
```

#### Backend：UpdateTodoLogic

更新时会覆盖 Redis 中的整条 JSON，需要保留 `is_only_for_test` 标记：

```go
// 读取旧数据，保留 is_only_for_test 标记
oldRaw, err := l.svcCtx.TodoRedis.Hget(todoKey, in.Id)
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
    IsOnlyForTest: oldIsOnlyForTest,  // 保留原值
}
```

#### Gateway：CreateTodoLogic

```go
rpcResp, err := l.svcCtx.TodoRpc.CreateTodo(l.ctx, &pb.CreateTodoReq{
    Title:         req.Title,
    IsOnlyForTest: req.IsOnlyForTest,  // 透传
})
```

#### Gateway：GetTodosLogic

```go
for _, t := range rpcResp.Todos {
    todos = append(todos, types.Todo{
        Id:            t.Id,
        Title:         t.Title,
        Completed:     t.Completed,
        IsOnlyForTest: t.IsOnlyForTest,  // 透传
    })
}
```

### 第 7 步：更新集成测试

```go
// testutil/types.go — 添加字段
type CreateTodoReq struct {
    Title         string `json:"title"`
    IsOnlyForTest bool   `json:"is_only_for_test"`
}

type Todo struct {
    Id            string `json:"id"`
    Title         string `json:"title"`
    Completed     bool   `json:"completed"`
    IsOnlyForTest bool   `json:"is_only_for_test"`
}

// cases/*.go — 创建测试数据时标记
resp := testutil.Post[testutil.CreateTodoResp](t, "/api/createTodo", testutil.CreateTodoReq{
    Title:         "买咖啡豆",
    IsOnlyForTest: true,  // 标记为纯测试数据
})
```

### 第 8 步：验证编译

```bash
cd /Users/liuzhan/full-stack/frontend/todo/backend && go build ./...
cd /Users/liuzhan/full-stack/frontend/todo/gateway && go build ./...
```

### 第 9 步：重建 Docker

```bash
cd /Users/liuzhan/full-stack/frontend/todo
sudo docker compose down
sudo docker compose up --build -d
# 无需手动清理数据：测试数据标记了 is_only_for_test: true，
# 且每个测试 case 在清理阶段会 HDel 自己的数据。
# 如需手动清理测试数据（不删真实数据）：
#   redis-cli EVAL "for _,k in ipairs(redis.call('keys','*')) do
#     local v=redis.call('hget','todos',k)
#     if v and cjson.decode(v).is_only_for_test then redis.call('hdel','todos',k) end
#   end" 0
```

### 第 10 步：跑集成测试

```bash
cd /Users/liuzhan/full-stack/frontend/todo/integrated_tests
go test -v ./cases/...
```

---

## 常见坑

| 问题 | 原因 | 解决 |
|------|------|------|
| `goctl rpc protoc` 生成的 pb 无新字段 | proto 文件没保存/语法错误 | 检查 `.proto` 格式，确认字段编号不重复 |
| `goctl api go` 生成的 types 不包含新字段 | goctl 版本 bug | 手动在 `types/types.go` 添加 |
| protoc `Missing input file` | proto 路径不对或跨目录 | 用绝对路径或正确相对路径 |
| goctl 生成 logic 骨架把自定义代码清掉 | `goctl api go` 在文件不存在时创建新文件 | 已存在的 logic 文件会自动跳过 |
| Docker 构建后代码没更新 | Docker cache | 用 `--build` 强制重建 |
| Config 字段冲突 `conflict key redis` | `RpcServerConf` 匿名嵌入后同名冲突 | 自定义字段改名（如 `Redis` → `TodoRedis`）|
| Redis 连接被拒绝 | 端口映射配错 | Docker 内部通信用容器端口（6379），非宿主机映射端口 |
