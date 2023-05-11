package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	storageDir     = "./charts"
	repoIndexFile  = "index.yaml"
	repoURL        = "http://example.com/helm/repo"
	kubectlCommand = "/usr/bin/kubectl"
)

type ChartList struct {
	Charts []string `yaml:"charts"`
}

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	indexCmd := flag.NewFlagSet("index", flag.ExitOnError)
	installCmd := flag.NewFlagSet("install", flag.ExitOnError)
	imagesCmd := flag.NewFlagSet("images", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Usage: <binary> <command> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if addCmd.NArg() < 1 {
			fmt.Println("Usage: <binary> add <chart-location>")
			os.Exit(1)
		}
		chartLocation := addCmd.Arg(0)
		err := addChart(chartLocation)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Chart added successfully.")

	case "index":
		indexCmd.Parse(os.Args[2:])
		err := generateRepoIndex()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Repository index generated successfully.")

	case "install":
		installCmd.Parse(os.Args[2:])
		args := installCmd.Args()
		if len(args) < 1 {
			fmt.Println("Usage: <binary> install chart <chart-name>")
			os.Exit(1)
		}
		chartName := args[0]
		err := installChart(chartName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Chart '%s' installed successfully.\n", chartName)

	case "images":
		imagesCmd.Parse(os.Args[2:])
		imageList, err := getImageList()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Container Images:")
		for _, image := range imageList {
			fmt.Println(image)
		}

	default:
		fmt.Println("Invalid command. Supported commands: add, index, install, images")
		os.Exit(1)
	}
}

func addChart(chartLocation string) error {
	// Create the charts directory if it does not exist
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		os.Mkdir(storageDir, 0755)
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

func generateRepoIndex() error {
	chartList, err := loadChartList()
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

	indexFile := filepath.Join(storageDir, repoIndexFile)
	err = saveIndexFile(indexData, indexFile)
	if err != nil {
		return err
	}

	return nil
}

func installChart(chartName string) error {
	chartList, err := loadChartList()
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

	cmd := exec.Command(kubectlCommand, "apply", "-f", filepath.Join(storageDir, chartLocation))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to install chart: %s", err.Error())
	}

	return nil
}

func getImageList() ([]string, error) {
	chartList, err := loadChartList()
	if err != nil {
		return nil, err
	}

	imageList := make([]string, 0)

	for _, chartLocation := range chartList.Charts {
		chartPath := filepath.Join(storageDir, chartLocation)
		valuesFile := filepath.Join(chartPath, "values.yaml")

		valuesData, err := ioutil.ReadFile(valuesFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read values file for chart '%s': %s", getChartName(chartLocation), err.Error())
		}

		var values map[string]interface{}
		err = yaml.Unmarshal(valuesData, &values)
		if err != nil {
			return nil, fmt.Errorf("failed to parse values file for chart '%s': %s", getChartName(chartLocation), err.Error())
		}

		containers := getContainersFromValues(values)
		imageList = append(imageList, containers...)
	}

	return imageList, nil
}

func loadChartList() (*ChartList, error) {
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

func saveChartList(chartList *ChartList) error {
	data, err := yaml.Marshal(chartList)
	if err != nil {
		return fmt.Errorf("failed to serialize chart list: %s", err.Error())
	}

	err = ioutil.WriteFile(filepath.Join(storageDir, "charts.yaml"), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save chart list: %s", err.Error())
	}

	return nil
}

func saveIndexFile(indexData map[string]interface{}, indexFile string) error {
	data, err := yaml.Marshal(indexData)
	if err != nil {
		return fmt.Errorf("failed to serialize index data: %s", err.Error())
	}

	err = ioutil.WriteFile(indexFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save index file: %s", err.Error())
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

func getContainersFromValues(values map[string]interface{}) []string {
	containers := make([]string, 0)

	for _, v := range values {
		if m, ok := v.(map[string]interface{}); ok {
			if image, ok := m["image"]; ok {
				if imageStr, ok := image.(string); ok {
					containers = append(containers, imageStr)
				}
			}
		}
	}

	return containers
}
