package cases

import (
	"encoding/json"
	"testing"

	"integrated_tests/testutil"
)

// 验证更新一个 todo 后：
//  1. HTTP 响应状态正常（200 OK）
//  2. Redis 中的 JSON 数据已更新为新内容
func TestUpdateTodo(t *testing.T) {
	// ---- 准备：创建一个 todo ----
	created := testutil.Post[testutil.CreateTodoResp](t, "/api/createTodo", testutil.CreateTodoReq{
		Title:         "旧标题",
		IsOnlyForTest: true,
	})
	t.Logf("已创建 todo ID=%s", created.Todo.Id)

	// ---- 执行：更新标题 + 标记完成 ----
	testutil.Post[testutil.UpdateTodoResp](t, "/api/updateTodo", testutil.UpdateTodoReq{
		Id:        created.Todo.Id,
		Title:     "新标题",
		Completed: true,
	})

	// ---- 验证 Redis 数据已更新 ----
	raw := testutil.HGet(t, "todos", created.Todo.Id)

	var stored testutil.Todo
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		t.Fatalf("解析 Redis JSON 失败: %v", err)
	}
	testutil.AssertEqual(t, created.Todo.Id, stored.Id)
	testutil.AssertEqual(t, "新标题", stored.Title)
	testutil.AssertEqual(t, true, stored.Completed)
	t.Logf("Redis 数据已更新: %+v", stored)

	// ---- 清理 ----
	testutil.HDel(t, "todos", created.Todo.Id)
}
