// src/hooks/useTodos.ts — 业务层：状态管理 + 业务逻辑

import { useState, useEffect, useCallback } from "react";
import type { Todo } from "../types";
import * as api from "../api/todo";

export function useTodos() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [filter, setFilter] = useState<"all" | "active" | "completed">("all");

  // 加载数据
  const loadTodos = useCallback(async () => {
    const data = await api.getTodos();
    setTodos(data);
  }, []);

  useEffect(() => {
    loadTodos();
  }, [loadTodos]);

  // 统计数据
  const total = todos.length;
  const completedCount = todos.filter((t) => t.completed).length;
  const activeCount = total - completedCount;

  // 筛选
  const filteredTodos = todos.filter((t) => {
    if (filter === "active") return !t.completed;
    if (filter === "completed") return t.completed;
    return true;
  });

  // CRUD 操作
  const addTodo = async (title: string) => {
    await api.createTodo(title);
    await loadTodos();
  };

  const toggleTodo = async (id: string) => {
    const todo = todos.find((t) => t.id === id);
    if (!todo) return;
    await api.updateTodo(id, todo.title, !todo.completed);
    await loadTodos();
  };

  const updateTitle = async (id: string, title: string) => {
    const todo = todos.find((t) => t.id === id);
    if (!todo) return;
    await api.updateTodo(id, title, todo.completed);
    await loadTodos();
  };

  const removeTodo = async (id: string) => {
    await api.deleteTodo(id);
    await loadTodos();
  };

  const clearCompleted = async () => {
    const completed = todos.filter((t) => t.completed);
    await Promise.all(completed.map((t) => api.deleteTodo(t.id)));
    await loadTodos();
  };

  return {
    todos: filteredTodos,
    filter,
    setFilter,
    total,
    completedCount,
    activeCount,
    addTodo,
    toggleTodo,
    updateTitle,
    removeTodo,
    clearCompleted,
  };
}
