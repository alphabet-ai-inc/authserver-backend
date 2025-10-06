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
	return fmt.Sprintf(`ThisApp error:
		ID: %d,
		%s`, t.ID, t.NewApp.Error())
}
