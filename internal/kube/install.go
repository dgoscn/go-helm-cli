package kube

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/dgoscn/go-helm-cli/internal/charts"
)

const (
	kubectlCommand = "/usr/bin/kubectl"
)

func InstallChart(chartName string) error {
	chartList, err := charts.LoadChartList()
	if err != nil {
		return err
	}

	var chartLocation string
	for _, cl := range chartList.Charts {
		if strings.Contains(cl, chartName) {
			chartLocation = cl
			break
		}
	}

	if chartLocation == "" {
		return fmt.Errorf("chart '%s' not found", chartName)
	}

	cmd := exec.Command(kubectlCommand, "apply", "-f", filepath.Join(charts.StorageDir, chartLocation))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to install chart '%s': %s", chartName, err.Error())
	}

	return nil
}

