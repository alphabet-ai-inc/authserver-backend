package models

import "fmt"

type ThisApp struct {
	ID      int    `json:"id"`
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
func (t ThisApp) Error() string {
	return fmt.Sprintf("ThisApp error: ID: %d, Name: %s, Release: %s, Path: %s, Init: %s, Web: %s, Title: %s, Created: %d, Updated: %d", t.ID, t.Name, t.Release, t.Path, t.Init, t.Web,
		t.Title, t.Created, t.Updated)
}

func (t ThisApp) String() string {
	return fmt.Sprintf("ThisApp{ID: %d, Name: %s, Release: %s, Path: %s, Init: %s, Web: %s, Title: %s, Created: %d, Updated: %d}", t.ID, t.Name, t.Release, t.Path, t.Init, t.Web,
		t.Title, t.Created, t.Updated)
}
