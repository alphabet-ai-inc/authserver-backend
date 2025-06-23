package models

import "fmt"

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
	return fmt.Sprintf("NewApp error: Name: %s, Release: %s, Path: %s, Init: %s, Web: %s, Title: %s, Created: %d, Updated: %d", t.Name, t.Release, t.Path, t.Init, t.Web, t.Title, t.Created, t.Updated)
}
