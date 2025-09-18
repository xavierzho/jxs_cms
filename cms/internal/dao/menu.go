package dao

type Menu struct {
	Title      string  `json:"title"`
	Name       string  `json:"name"`
	Path       string  `json:"path"`
	Permission string  `json:"-"`
	Show       bool    `json:"show"`
	Children   []*Menu `json:"children,omitempty"`
}
