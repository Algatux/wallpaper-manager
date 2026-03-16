package setup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfiguration_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.json")

	expectedConfig := MonitorList{
		Monitors: []MonitorConfig{
			{
				Name:  "DP-1",
				Path:  "/home/user/wallpapers/main",
				Flags: "",
			},
			{
				Name:  "DP-2",
				Path:  "/home/user/wallpapers/side",
				Flags: "-m 16:9",
			},
		},
	}

	data, err := json.Marshal(expectedConfig)
	if err != nil {
		t.Fatalf("Failed to marshal expected config: %v", err)
	}
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	inputs := Inputs{ConfigFile: configFile}
	result, err := inputs.GetConfiguration()

	if err != nil {
		t.Errorf("GetConfiguration returned unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("GetConfiguration returned nil result")
	}

	if len(result.Monitors) != len(expectedConfig.Monitors) {
		t.Errorf("Expected %d monitors, got %d", len(expectedConfig.Monitors), len(result.Monitors))
	}

	for i, monitor := range result.Monitors {
		if monitor.Name != expectedConfig.Monitors[i].Name {
			t.Errorf("Monitor %d: expected name %q, got %q", i, expectedConfig.Monitors[i].Name, monitor.Name)
		}
		if monitor.Path != expectedConfig.Monitors[i].Path {
			t.Errorf("Monitor %d: expected path %q, got %q", i, expectedConfig.Monitors[i].Path, monitor.Path)
		}
	}
}

func TestGetConfiguration_FileNotFound(t *testing.T) {
	inputs := Inputs{ConfigFile: "/nonexistent/path/config.json"}
	result, err := inputs.GetConfiguration()

	if err == nil {
		t.Error("GetConfiguration should return error for non-existent file")
	}

	if result != nil {
		t.Error("GetConfiguration should return nil result on error")
	}
}

func TestGetConfiguration_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid_config.json")

	invalidJSON := []byte(`{invalid json content}`)
	if err := os.WriteFile(configFile, invalidJSON, 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	inputs := Inputs{ConfigFile: configFile}
	result, err := inputs.GetConfiguration()

	if err == nil {
		t.Error("GetConfiguration should return error for invalid JSON")
	}

	if result != nil {
		t.Error("GetConfiguration should return nil result on error")
	}
}

func TestGetConfiguration_EmptyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "empty_config.json")

	emptyConfig := MonitorList{Monitors: []MonitorConfig{}}
	data, _ := json.Marshal(emptyConfig)
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write empty config file: %v", err)
	}

	inputs := Inputs{ConfigFile: configFile}
	result, err := inputs.GetConfiguration()

	if err != nil {
		t.Errorf("GetConfiguration returned unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("GetConfiguration returned nil result")
	}

	if len(result.Monitors) != 0 {
		t.Errorf("Expected 0 monitors, got %d", len(result.Monitors))
	}
}

func TestMonitorConfigStruct(t *testing.T) {
	monitor := MonitorConfig{
		Name:  "DP-1",
		Path:  "/home/user/wallpapers",
		Flags: "-s fit",
	}

	if monitor.Name != "DP-1" {
		t.Errorf("Expected name DP-1, got %s", monitor.Name)
	}
	if monitor.Path != "/home/user/wallpapers" {
		t.Errorf("Expected path /home/user/wallpapers, got %s", monitor.Path)
	}
	if monitor.Flags != "-s fit" {
		t.Errorf("Expected flags -s fit, got %s", monitor.Flags)
	}
}

func TestMonitorListStruct(t *testing.T) {
	monitors := []MonitorConfig{
		{Name: "DP-1", Path: "/path1"},
		{Name: "DP-2", Path: "/path2"},
	}

	list := MonitorList{Monitors: monitors}

	if len(list.Monitors) != 2 {
		t.Errorf("Expected 2 monitors, got %d", len(list.Monitors))
	}

	if list.Monitors[0].Name != "DP-1" {
		t.Errorf("Expected first monitor name DP-1, got %s", list.Monitors[0].Name)
	}
}
