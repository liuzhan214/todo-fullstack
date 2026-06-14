package cases

import (
	"testing"

	"integrated_tests/testutil"
)

// 验证列出所有 todo 时：
//  1. HTTP 响应返回完整的 todo 列表
//  2. 列表包含之前创建的所有待办项
func TestGetTodos(t *testing.T) {
	// ---- 准备：创建 2 个 todo ----
	t1 := testutil.Post[testutil.CreateTodoResp](t, "/api/createTodo", testutil.CreateTodoReq{
		Title:         "写周报",
		IsOnlyForTest: true,
	})
	t2 := testutil.Post[testutil.CreateTodoResp](t, "/api/createTodo", testutil.CreateTodoReq{
		Title:         "回复邮件",
		IsOnlyForTest: true,
	})

	// ---- 执行：获取列表 ----
	resp := testutil.Get[testutil.GetTodosResp](t, "/api/getTodos")
	testutil.PrintJSON(t, "列表响应", resp)

	// ---- 验证：列表包含刚才创建的任务 ----
	found := make(map[string]bool)
	for _, todo := range resp.Todos {
		if todo.Id == t1.Todo.Id {
			found["t1"] = true
			testutil.AssertEqual(t, "写周报", todo.Title)
		}
		if todo.Id == t2.Todo.Id {
			found["t2"] = true
			testutil.AssertEqual(t, "回复邮件", todo.Title)
		}
	}
	if !found["t1"] {
		t.Fatal("列表中未找到刚创建的 todo (写周报)")
	}
	if !found["t2"] {
		t.Fatal("列表中未找到刚创建的 todo (回复邮件)")
	}

	// ---- 清理 ----
	testutil.HDel(t, "todos", t1.Todo.Id, t2.Todo.Id)
}
