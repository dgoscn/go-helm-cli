package main

import (
	"path/filepath"
	"testing"
)

func TestAddChart(t *testing.T) {
	// Set up the test
	tempDir := t.TempDir()
	storageDir := filepath.Join(tempDir, "charts")

	// Override the original storageDir variable
	originalStorageDir := storageDir
	defer func() { storageDir = originalStorageDir }()

	// Define the expected chart location
	expectedChartLocation := "test/chart"

	// Run the addChart function with the expected chart location
	err := addChart(expectedChartLocation)
	if err != nil {
		t.Errorf("Failed to add chart: %s", err)
	}

	// Load the chart list
	chartList, err := loadChartList()
	if err != nil {
		t.Errorf("Failed to load chart list: %s", err)
	}

	// Check if the expected chart location is present in the chart list
	found := false
	for _, chartLocation := range chartList.Charts {
		if chartLocation == expectedChartLocation {
			found = true
			break
		}
	}

	// Verify the result
	if !found {
		t.Errorf("Unexpected chart location. Expected: %s, Actual: %v", expectedChartLocation, chartList.Charts)
	}

	// Test adding multiple charts
	expectedChartLocation2 := "test/chart2"
	err = addChart(expectedChartLocation2)
	if err != nil {
		t.Errorf("Failed to add chart: %s", err)
	}

	chartList, err = loadChartList()
	if err != nil {
		t.Errorf("Failed to load chart list: %s", err)
	}

	found = false
	for _, chartLocation := range chartList.Charts {
		if chartLocation == expectedChartLocation2 {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Unexpected chart location. Expected: %s, Actual: %v", expectedChartLocation2, chartList.Charts)
	}

	// Test adding a chart with an error
	err = addChart("invalid-chart")
	if err == nil {
		t.Errorf("Expected an error but got nil")
	}

}
