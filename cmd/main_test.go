package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

var testStorageDir string

func TestAddChart(t *testing.T) {
	// Create a temporary directory for storage
	tmpDir, err := ioutil.TempDir("", "test-storage")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Assign the temporary directory path to the package-level testStorageDir variable
	testStorageDir = tmpDir

	// Clear the chart list
	err = clearChartList()
	if err != nil {
		t.Fatalf("Failed to clear chart list: %v", err)
	}

	// Call the addChart function
	chartLocation := "mycharts/mychart"
	err = addChart(chartLocation)
	if err != nil {
		t.Fatalf("Failed to add chart: %v", err)
	}

	// Load the chart list and verify the added chart
	chartList, err := loadChartList()
	if err != nil {
		t.Fatalf("Failed to load chart list: %v", err)
	}

	if len(chartList.Charts) != 1 {
		t.Fatalf("Unexpected number of charts in the chart list. Expected 1, got %d", len(chartList.Charts))
	}

	if chartList.Charts[0] != chartLocation {
		t.Fatalf("Unexpected chart location in the chart list. Expected %s, got %s", chartLocation, chartList.Charts[0])
	}
}

func clearChartList() error {
	chartList := &ChartList{Charts: []string{}}
	err := saveChartList(chartList)
	if err != nil {
		return fmt.Errorf("failed to clear chart list: %s", err.Error())
	}
	return nil
}


func loadTestChartList(storageDir string) (*ChartList, error) {
	chartList := &ChartList{}

	data, err := ioutil.ReadFile(filepath.Join(storageDir, "charts.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to load chart list: %s", err.Error())
	}

	err = yaml.Unmarshal(data, chartList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chart list: %s", err.Error())
	}

	return chartList, nil
}
