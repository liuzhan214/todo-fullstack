import type { Todo } from "../types";
import TodoItem from "./TodoItem";

export default function TodoList({
  todos,
  onToggle,
  onUpdate,
  onDelete,
  total,
  completedCount,
  activeCount,
}: {
  todos: Todo[];
  onToggle: (id: string) => void;
  onUpdate: (id: string, title: string) => void;
  onDelete: (id: string) => void;
  total: number;
  completedCount: number;
  activeCount: number;
}) {
  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-900">任务列表</h2>
        <span className="text-sm text-gray-500">
          {total} 项任务 | {completedCount} 项已完成 | {activeCount} 项待办
        </span>
      </div>

      <div className="space-y-2 mb-6">
        {todos.length === 0 ? (
          <p className="text-center text-gray-400 py-8 text-sm">暂无任务</p>
        ) : (
          todos.map((todo) => (
            <TodoItem
              key={todo.id}
              todo={todo}
              onToggle={onToggle}
              onUpdate={onUpdate}
              onDelete={onDelete}
            />
          ))
        )}
      </div>
    </div>
  );
}
