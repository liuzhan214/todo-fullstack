// src/components/TodoItem.tsx
import { useState } from "react";
import type { Todo } from "../types";

interface Props {
  todo: Todo;
  onToggle: (id: string) => void;
  onUpdate: (id: string, title: string) => void;
  onDelete: (id: string) => void;
}

export default function TodoItem({
  todo,
  onToggle,
  onUpdate,
  onDelete,
}: Props) {
  const [editing, setEditing] = useState(false);
  const [editTitle, setEditTitle] = useState(todo.title);

  const handleSave = () => {
    if (editTitle.trim() && editTitle !== todo.title) {
      onUpdate(todo.id, editTitle);
    }
    setEditing(false);
  };

  return (
    <div
      className={`flex items-center gap-3 px-4 py-3 rounded-lg 
                ${
                  todo.completed
                    ? "bg-gray-50"
                    : "bg-white border border-gray-200"
                }`}
    >
      {/* 勾选框 */}
      <button
        onClick={() => onToggle(todo.id)}
        className={`w-6 h-6 rounded-full flex
                    items-center justify-center flex-shrink-0
                    cursor-pointer transition-colors 
                    ${
                      todo.completed
                        ? "bg-indigo-600 text-white"
                        : "border-2 border-gray-300 bg-white"
                    }`}
      >
        {todo.completed && (
          <svg
            className="w-4 h-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2.5}
              d="M5 13 l4 4 L19 7"
            />
          </svg>
        )}
      </button>

      {/* 标题区 */}
      {editing ? (
        <div className="flex-1 flex gap-2">
          <input
            type="text"
            value={editTitle}
            onChange={(e) => setEditTitle(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSave()}
            className="flex-1 px-3 py-1.5 border border-indigo-500 rounded text-sm focus:outline-none"
            autoFocus
          />
          <button
            onClick={handleSave}
            className="px-3 py-1.5 bg-indigo-600 text-white rounded text-xs font-medium cursor-pointer hover:bg-indigo-700"
          >
            保存
          </button>
          <button
            onClick={() => setEditing(false)}
            className="px-3 py-1.5 bg-gray-200 text-gray-600 rounded text-xs cursor-pointer hover:bg-gray-300"
          >
            取消
          </button>
        </div>
      ) : (
        <span
          onDoubleClick={() => setEditing(true)}
          className={`flex-1 text-sm cursor-pointer 
                ${
                  todo.completed
                    ? "text-gray-400 line-through"
                    : "text-gray-900"
                }`}
        >
          {todo.title}
        </span>
      )}

      {/* 删除按钮 */}
      <button
        onClick={() => onDelete(todo.id)}
        className="text-red-500 hover:text-red-700 cursor-pointer text-sm flex-shrink-0"
      >
        🗑
      </button>
    </div>
  );
}
