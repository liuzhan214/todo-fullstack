import { useTodos } from "./hooks/useTodos";
import AddTodo from "./components/AddTodo";
import TodoList from "./components/TodoList";
import FilterBar from "./components/FilterBar";

function App() {
  const {
    todos,
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
  } = useTodos();

  return (
    <div className="max-w-[700px] mx-auto p-5">
      <header className="bg-gradient-to-r from-indigo-600 to-purple-600 rounded-t-2xl p-6 text-center">
        <h1 className="text-white text-2xl font-bold">📋 Todo App</h1>
      </header>

      <main className="bg-white rounded-b-2xl p-6 shadow-md">
        <AddTodo onAdd={addTodo} />
        <hr className="my-6 border-gray-200" />
        <TodoList
          todos={todos}
          onToggle={toggleTodo}
          onUpdate={updateTitle}
          onDelete={removeTodo}
          total={total}
          completedCount={completedCount}
          activeCount={activeCount}
        />
        <hr className="my-6 border-gray-200" />
        <FilterBar
          filter={filter}
          onFilterChange={setFilter}
          onClearCompleted={clearCompleted}
          hasCompleted={completedCount > 0}
        />

        <p className="mt-6 text-gray-400 text-xs text-center">
          💡 勾选任务 → 标记完成 &nbsp;|&nbsp; 双击标题 → 进入编辑 &nbsp;|&nbsp;
          实时同步到后端
        </p>
      </main>
    </div>
  );
}

export default App;
