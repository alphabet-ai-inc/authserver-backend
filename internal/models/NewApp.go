package models

type NewApp struct {
	Name    string `json:"name"`
	Release string `json:"release"`
	Path    string `json:"path"`
	Init    string `json:"init"`
	Web     string `json:"web"`
	Title   string `json:"title"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
}

// Error implements error.
func (t NewApp) Error() string {
	panic("unimplemented")
}
