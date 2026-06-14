package cases

import (
	"testing"

	"integrated_tests/testutil"
)

// 验证删除一个 todo 后：
//  1. HTTP 响应状态正常（200 OK）
//  2. Redis 中对应的 hash field 已被删除
func TestDeleteTodo(t *testing.T) {
	// ---- 准备：创建一个 todo ----
	created := testutil.Post[testutil.CreateTodoResp](t, "/api/createTodo", testutil.CreateTodoReq{
		Title:         "待删除任务",
		IsOnlyForTest: true,
	})
	t.Logf("已创建 todo ID=%s", created.Todo.Id)

	// ---- 执行：删除 ----
	testutil.Post[testutil.DeleteTodoResp](t, "/api/deleteTodo", testutil.DeleteTodoReq{
		Id: created.Todo.Id,
	})

	// ---- 验证 Redis 中已删除 ----
	all := testutil.HGetAll(t, "todos")
	if _, exists := all[created.Todo.Id]; exists {
		t.Fatalf("Redis 中 todo %s 仍然存在", created.Todo.Id)
	}
	t.Log("Redis 中该 todo 已被正确删除")
}
