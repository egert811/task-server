package storage

type TaskDBItem struct {
	ID  int    `json:"id"`
	CMD string `json:"cmd"`
}

type TaskOutputDBItem struct {
	ID     int    `json:"-"`
	Output string `json:"output"`
}