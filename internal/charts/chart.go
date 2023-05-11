package charts

import (
	"fmt"
	"os"

	"github.com/dgoscn/go-helm-cli/internal/common"
)

const (
	storageDir = "./charts"
)

type ChartList struct {
	Charts []string `yaml:"charts"`
}

func AddChart(chartLocation string) error {
	// Create the charts directory if it does not exist
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		err := os.Mkdir(storageDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create charts directory: %s", err.Error())
		}
	}

	// Load the existing chart list or create a new one if it does not exist
	chartList, err := loadChartList()
	if err != nil {
		chartList = &ChartList{}
	}

	// Append the chart location to the chart list
	chartList.Charts = append(chartList.Charts, chartLocation)

	// Save the updated chart list
	err = saveChartList(chartList)
	if err != nil {
		return fmt.Errorf("failed to save chart list: %s", err.Error())
	}

	return nil
}

func loadChartList() (*ChartList, error) {
	chartList := &ChartList{}

	data, err := utils.ReadFile(filepath.Join(storageDir, "charts.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to load chart list: %s", err.Error())
	}

	err = utils.UnmarshalYAML(data, chartList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chart list: %s", err.Error())
	}

	return chartList, nil
}

func saveChartList(chartList *ChartList) error {
	data, err := utils.MarshalYAML(chartList)
	if err != nil {
		return fmt.Errorf("failed to serialize chart list: %s", err.Error())
	}

	err = utils.WriteFile(filepath.Join(storageDir, "charts.yaml"), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save chart list: %s", err.Error())
	}

	return nil
}

