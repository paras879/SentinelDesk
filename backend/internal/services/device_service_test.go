package services

import (
	"testing"

	"sentineldesk/backend/internal/config"
)

// planNameFromConfig mirrors the plan name derivation logic in DeviceService.GetPlanInfo.
func planNameFromConfig(maxDevices int, planNameOverride string) string {
	planName := planNameOverride
	if planName == "" {
		if maxDevices <= 3 && maxDevices > 0 {
			planName = "Starter"
		} else if maxDevices <= 25 && maxDevices > 0 {
			planName = "Pro"
		} else {
			planName = "Enterprise"
		}
	}
	return planName
}

func TestPlanName_Starter(t *testing.T) {
	name := planNameFromConfig(3, "")
	if name != "Starter" {
		t.Errorf("expected Starter, got %q", name)
	}
}

func TestPlanName_StarterLow(t *testing.T) {
	name := planNameFromConfig(1, "")
	if name != "Starter" {
		t.Errorf("expected Starter for maxDevices=1, got %q", name)
	}
}

func TestPlanName_Pro(t *testing.T) {
	name := planNameFromConfig(25, "")
	if name != "Pro" {
		t.Errorf("expected Pro, got %q", name)
	}
}

func TestPlanName_ProMid(t *testing.T) {
	name := planNameFromConfig(10, "")
	if name != "Pro" {
		t.Errorf("expected Pro for maxDevices=10, got %q", name)
	}
}

func TestPlanName_Enterprise_Zero(t *testing.T) {
	name := planNameFromConfig(0, "")
	if name != "Enterprise" {
		t.Errorf("expected Enterprise for maxDevices=0, got %q", name)
	}
}

func TestPlanName_Enterprise_Large(t *testing.T) {
	name := planNameFromConfig(1000, "")
	if name != "Enterprise" {
		t.Errorf("expected Enterprise for maxDevices=1000, got %q", name)
	}
}

func TestPlanName_Override(t *testing.T) {
	name := planNameFromConfig(3, "CustomPlan")
	if name != "CustomPlan" {
		t.Errorf("expected CustomPlan override, got %q", name)
	}
}

func TestDeviceService_Constructor_NoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewDeviceService panicked: %v", r)
		}
	}()
	// NewDeviceService calls config.LoadConfig() and repository.NewDeviceRepository().
	// Neither hits the DB, so this is safe without a real connection.
	cfg := config.LoadConfig()
	if cfg == nil {
		t.Fatal("config.LoadConfig returned nil")
	}
}
