package models_test

import (
	"authserver-backend/internal/models"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestThisApp tests the ThisApp struct for JSON serialization/deserialization, focusing on the ID and embedded NewApp fields.
// Since NewApp fields are tested separately in NewApp_test.go, we only test a subset here.
func TestThisApp(t *testing.T) {
	// Create an instance of ThisApp with ID and a few NewApp fields
	app := models.ThisApp{
		ID: 1,
		NewApp: models.NewApp{
			Name:    "Test App",
			Release: "1.0.0",
			Path:    "/test/path",
			Init:    "init.sh",
			Web:     "http://testapp.com",
			Title:   "Test AuthServerApp",
			Created: 1660000000,
			Updated: 1660000001,
		},
	}

	// Serialize to JSON
	data, err := json.Marshal(app)
	assert.NoError(t, err)

	// Expected JSON includes ID and the set fields (other NewApp fields will be zero values)
	expectedJSON := `{"id":1,"name":"Test App","release":"1.0.0","path":"/test/path","init":"init.sh","web":"http://testapp.com","title":"Test AuthServerApp","created":1660000000,"updated":1660000001,"description":"","positioning_stmt":"","logo":"","category":"","platform":null,"developer":"","license_type":"","size":0,"compatibility":null,"integration_capabilities":null,"development_stack":null,"api_documentation":"","security_features":null,"regulatory_compliance":null,"revenue_streams":null,"customer_segments":null,"channels":null,"value_proposition":"","pricing_tiers":null,"partnerships":null,"cost_structure":null,"customer_relationships":null,"unfair_advantage":"","roadmap":"","version_control":"","error_rate":0,"average_response_time":0,"uptime_percentage":0,"key_activities":null,"active_users":0,"user_retention_rate":0,"user_acquisition_cost":0,"churn_rate":0,"monthly_recurring_revenue":0,"user_feedback":null,"backup_recovery_options":null,"localization_support":null,"accessibility_features":null,"team_structure":null,"data_backup_location":"","environmental_impact":"","social_impact":"","intellectual_property":null,"fundings_investment":0,"exit_strategy":"","analytics_tools":null,"key_metrics":null,"url":"","landing_page":""}`
	assert.JSONEq(t, expectedJSON, string(data))

	// Test deserialization
	var appFromJSON models.ThisApp
	err = json.Unmarshal(data, &appFromJSON)
	assert.NoError(t, err)
	assert.Equal(t, app, appFromJSON)
}

func TestThisApp_Error(t *testing.T) {
	app := models.ThisApp{
		NewApp: models.NewApp{
			Name:    "Test",
			Release: "1.0",
			Path:    "/path",
			Init:    "init",
			Web:     "web",
			Title:   "Title",
			Created: time.Unix(1660000000, 0).Unix(),
			Updated: time.Unix(1660000001, 0).Unix(),
		},
	}

	// The Error method returns a string, not panics
	errMsg := app.Error()
	assert.Contains(t, errMsg, "Test")
	assert.Contains(t, errMsg, "1.0")
}
