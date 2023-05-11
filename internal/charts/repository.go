package charts

import (
	"fmt"
	"path/filepath"

	"github.com/dgoscn/go-helm-cli/internal/common"
	"github.com/dgoscn/go-helm-cli/internal/utils"
)

const (
	repoIndexFile = "index.yaml"
	repoURL       = "http://example.com/helm/repo"
)

func GenerateRepoIndex() error {
	chartList, err := charts.LoadChartList()
	if err != nil {
		return err
	}

	indexData := map[string]interface{}{
		"apiVersion": "v1",
		"entries": map[string]interface{}{
			"": []map[string]interface{}{},
		},
	}

	for _, chartLocation := range chartList.Charts {
		chartName := getChartName(chartLocation)
		entry := map[string]interface{}{
			"name": chartName,
			"url":  fmt.Sprintf("%s/%s-%s.tgz", repoURL, chartName, getVersionFromChartLocation(chartLocation)),
		}

		indexData["entries"].(map[string]interface{})[""] = append(indexData["entries"].(map[string]interface{})[""].([]map[string]interface{}), entry)
	}

	indexFile := filepath.Join(charts.StorageDir, repoIndexFile)
	err = utils.SaveIndexFile(indexData, indexFile)
	if err != nil {
		return err
	}

	return nil
}

func getChartName(chartLocation string) string {
	parts := strings.Split(chartLocation, "/")
	return parts[len(parts)-1]
}

func getVersionFromChartLocation(chartLocation string) string {
	parts := strings.Split(chartLocation, "/")
	return parts[len(parts)-2]
}
