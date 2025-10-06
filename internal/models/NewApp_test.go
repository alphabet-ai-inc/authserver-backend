package models_test

import (
	"authserver-backend/internal/models"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewApp tests the NewApp struct for JSON serialization and deserialization.
func TestNewApp(t *testing.T) {
	// Create an instance of NewApp with assigned values
	app := models.NewApp{
		Name:                    "TestApp",
		Release:                 "1.0.0",
		Path:                    "/test/path",
		Init:                    "init.sh",
		Web:                     "http://testapp.com",
		Title:                   "Test AuthServerApp",
		Created:                 1660000000,
		Updated:                 1660000001,
		Description:             "A test app",
		PositioningStmt:         "Positioning statement",
		Logo:                    "http://logo.com",
		Category:                "Test",
		Platform:                []string{"Web"},
		Developer:               "TestDev",
		LicenseType:             "MIT",
		Size:                    1024,
		Compatibility:           []string{"v1.0"},
		IntegrationCapabilities: []string{"API"},
		DevelopmentStack:        []string{"Go"},
		APIDocumentation:        "http://api.com",
		SecurityFeatures:        []string{"Auth"},
		RegulatoryCompliance:    []string{"GDPR"},
		RevenueStreams:          []string{"Subscription"},
		CustomerSegments:        []string{"Developers"},
		Channels:                []string{"Web"},
		ValueProposition:        "Value prop",
		PricingTiers:            []string{"Free", "Pro"},
		Partnerships:            []string{"Partner1"},
		CostStructure:           []string{"Hosting"},
		CustomerRelationships:   []string{"Support"},
		UnfairAdvantage:         "Unique feature",
		Roadmap:                 "http://roadmap.com",
		VersionControl:          "http://git.com",
		ErrorRate:               0.01,
		AverageResponseTime:     0.5,
		UptimePercentage:        99.9,
		KeyActivities:           []string{"Development"},
		ActiveUsers:             100,
		UserRetentionRate:       80.0,
		UserAcquisitionCost:     10.0,
		ChurnRate:               5.0,
		MonthlyRecurringRevenue: 1000.0,
		UserFeedback:            []string{"Good"},
		BackupRecoveryOptions:   []string{"Daily"},
		LocalizationSupport:     []string{"EN", "ES"},
		AccessibilityFeatures:   []string{"WCAG"},
		TeamStructure:           []string{"Dev", "QA"},
		DataBackupLocation:      "Cloud",
		EnvironmentalImpact:     "Low",
		SocialImpact:            "Positive",
		IntellectualProperty:    []string{"Patent"},
		FundingsInvestment:      50000.0,
		ExitStrategy:            "IPO",
		AnalyticsTools:          []string{"Google Analytics"},
		KeyMetrics:              []string{"Users"},
		URL:                     "http://app.com",
		LandingPage:             "http://landing.com",
	}

	// Serialize to JSON
	data, err := json.Marshal(app)
	assert.NoError(t, err)

	expectedJSON := `{"name":"TestApp","release":"1.0.0","path":"/test/path","init":"init.sh","web":"http://testapp.com","title":"Test AuthServerApp","created":1660000000,"updated":1660000001,"description":"A test app","positioning_stmt":"Positioning statement","logo":"http://logo.com","category":"Test","platform":["Web"],"developer":"TestDev","license_type":"MIT","size":1024,"compatibility":["v1.0"],"integration_capabilities":["API"],"development_stack":["Go"],"api_documentation":"http://api.com","security_features":["Auth"],"regulatory_compliance":["GDPR"],"revenue_streams":["Subscription"],"customer_segments":["Developers"],"channels":["Web"],"value_proposition":"Value prop","pricing_tiers":["Free","Pro"],"partnerships":["Partner1"],"cost_structure":["Hosting"],"customer_relationships":["Support"],"unfair_advantage":"Unique feature","roadmap":"http://roadmap.com","version_control":"http://git.com","error_rate":0.01,"average_response_time":0.5,"uptime_percentage":99.9,"key_activities":["Development"],"active_users":100,"user_retention_rate":80,"user_acquisition_cost":10,"churn_rate":5,"monthly_recurring_revenue":1000,"user_feedback":["Good"],"backup_recovery_options":["Daily"],"localization_support":["EN","ES"],"accessibility_features":["WCAG"],"team_structure":["Dev","QA"],"data_backup_location":"Cloud","environmental_impact":"Low","social_impact":"Positive","intellectual_property":["Patent"],"fundings_investment":50000,"exit_strategy":"IPO","analytics_tools":["Google Analytics"],"key_metrics":["Users"],"url":"http://app.com","landing_page":"http://landing.com"}`
	assert.JSONEq(t, expectedJSON, string(data))

	// Test deserialization
	var newAppFromJSON models.NewApp
	err = json.Unmarshal(data, &newAppFromJSON)
	assert.NoError(t, err)
	assert.Equal(t, app, newAppFromJSON)
}

func TestNewApp_Error(t *testing.T) {
	app := models.NewApp{
		Name:    "Test",
		Release: "1.0",
		Path:    "/path",
		Init:    "init",
		Web:     "web",
		Title:   "Title",
		Created: time.Unix(1660000000, 0).Unix(),
		Updated: time.Unix(1660000000, 0).Unix(),
	}

	// The Error method returns a string, not panics
	errMsg := app.Error()
	assert.Contains(t, errMsg, "Test")
	assert.Contains(t, errMsg, "1.0")
}
