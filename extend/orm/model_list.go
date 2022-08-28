package orm

// ModelList 通用模型列表
type ModelList[T IModel] struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"pageSize"`
	PageCount int   `json:"pageCount"`
	Total     int64 `json:"total"`
	Items     []T   `json:"items"`
}
