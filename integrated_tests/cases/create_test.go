package cases

import (
	"testing"

	"integrated_tests/testutil"
)

// 验证创建一个 todo 后：
//  1. HTTP 响应返回正确的 ID、标题、完成状态
//  2. Redis hash 中持久化了对应的 JSON 数据
func TestCreateTodo(t *testing.T) {
	// ---- 准备 ----
	title := "买咖啡豆"

	// ---- 执行 ----
	resp := testutil.Post[testutil.CreateTodoResp](t, "/api/createTodo", testutil.CreateTodoReq{
		Title:         title,
		IsOnlyForTest: true,
	})

	// ---- 验证响应 ----
	testutil.PrintJSON(t, "创建响应", resp)
	testutil.AssertNotEmpty(t, "ID", resp.Todo.Id)
	testutil.AssertEqual(t, title, resp.Todo.Title)
	testutil.AssertEqual(t, false, resp.Todo.Completed)

	// ---- 验证 Redis 持久化 ----
	raw := testutil.HGet(t, "todos", resp.Todo.Id)
	t.Logf("Redis [HGET todos %s]: %s", resp.Todo.Id, raw)

	// ---- 清理 ----
	testutil.HDel(t, "todos", resp.Todo.Id)
}
