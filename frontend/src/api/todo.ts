// src/api/todo.ts — 存储层：所有 HTTP 请求

export interface ApiTodo {
  id: string;
  title: string;
  completed: boolean;
  isOnlyForTest: boolean;
}

export interface GetTodosResp {
  todos: ApiTodo[];
}

export async function getTodos(): Promise<ApiTodo[]> {
  const res = await fetch("/api/getTodos");
  const data: GetTodosResp = await res.json();
  return data.todos || [];
}

export async function createTodo(title: string) {
  await fetch("/api/createTodo", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, is_only_for_test: false }),
  });
}

export async function updateTodo(
  id: string,
  title: string,
  completed: boolean,
) {
  await fetch("/api/updateTodo", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id, title, completed }),
  });
}

export async function deleteTodo(id: string) {
  await fetch("/api/deleteTodo", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id }),
  });
}
