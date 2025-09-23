package models

import "fmt"

// ThisApp represents the application itself, including its ID and embedded NewApp details.
// JSON tags are included for serialization/deserialization.
type ThisApp struct {
	ID     int `json:"id"`
	NewApp     // `json:"new_app"`
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
