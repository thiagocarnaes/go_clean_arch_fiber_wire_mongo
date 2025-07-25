package acceptance

import (
	"testing"
)

func TestSetup(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Simple test to verify that the setup works
	if testApp.App == nil {
		t.Fatal("App should not be nil")
	}

	if testApp.DB == nil {
		t.Fatal("DB should not be nil")
	}

	if testApp.Container == nil {
		t.Fatal("Container should not be nil")
	}
}
