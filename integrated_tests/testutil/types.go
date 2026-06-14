package testutil

// ---------- HTTP 响应类型（匹配 Gateway JSON） ----------

type Todo struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Completed     bool   `json:"completed"`
	IsOnlyForTest bool   `json:"is_only_for_test"`
}

type CreateTodoResp struct {
	Todo Todo `json:"todo"`
}

type GetTodosResp struct {
	Todos []Todo `json:"todos"`
}

type UpdateTodoResp struct{}

type DeleteTodoResp struct{}

// ---------- 请求类型 ----------

type CreateTodoReq struct {
	Title         string `json:"title"`
	IsOnlyForTest bool   `json:"is_only_for_test"`
}

type DeleteTodoReq struct {
	Id string `json:"id"`
}

type UpdateTodoReq struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
