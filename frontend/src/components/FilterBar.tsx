export default function FilterBar({
  filter,
  onFilterChange,
  onClearCompleted,
  hasCompleted,
}: {
  filter: "all" | "active" | "completed";
  onFilterChange: (f: "all" | "active" | "completed") => void;
  onClearCompleted: () => void;
  hasCompleted: boolean;
}) {
  const btn = (active: boolean) =>
    `px-4 py-2 rounded-lg text-sm font-medium cursor-pointer transition-colors ${
      active
        ? "bg-indigo-600 text-white"
        : "bg-gray-100 text-gray-700 hover:bg-gray-200"
    }`;

  return (
    <div>
      <h2 className="text-lg font-semibold text-gray-900 mb-3">筛选</h2>
      <div className="flex gap-2 mb-4">
        <button
          className={btn(filter === "all")}
          onClick={() => onFilterChange("all")}
        >
          全部
        </button>
        <button
          className={btn(filter === "active")}
          onClick={() => onFilterChange("active")}
        >
          待办
        </button>
        <button
          className={btn(filter === "completed")}
          onClick={() => onFilterChange("completed")}
        >
          已完成
        </button>

        {hasCompleted && (
          <button
            onClick={onClearCompleted}
            className="ml-auto px-4 py-2 bg-red-50 text-red-500 rounded-lg text-sm font-medium cursor-pointer hover:bg-red-100 transition-colors border border-red-200"
          >
            🗑 清空已完成
          </button>
        )}
      </div>
    </div>
  );
}
