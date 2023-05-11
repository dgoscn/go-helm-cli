package charts

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetImageList() ([]string, error) {
	cmd := exec.Command("helm", "list", "-a", "-q")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get chart list: %s", err.Error())
	}

	chartNames := strings.Split(strings.TrimSpace(string(output)), "\n")
	imageList := make([]string, 0)

	for _, chartName := range chartNames {
		cmd := exec.Command("helm", "get", "all", chartName, "-o", "yaml")
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get chart details for '%s': %s", chartName, err.Error())
		}

		images := parseImagesFromChartDetails(output)
		imageList = append(imageList, images...)
	}

	return imageList, nil
}

func parseImagesFromChartDetails(chartDetails []byte) []string {
	var data map[string]interface{}
	err := yaml.Unmarshal(chartDetails, &data)
	if err != nil {
		log.Printf("Failed to parse chart details YAML: %v", err)
		return []string{}
	}

	containers := make([]string, 0)

	if values, ok := data["values"].(map[interface{}]interface{}); ok {
		// Traverse the values to find container images
		traverseValues(values, &containers)
	}

	return containers
}

func traverseValues(values map[interface{}]interface{}, containers *[]string) {
	for _, v := range values {
		switch val := v.(type) {
		case map[interface{}]interface{}:
			traverseValues(val, containers)
		case []interface{}:
			for _, item := range val {
				if subValues, ok := item.(map[interface{}]interface{}); ok {
					traverseValues(subValues, containers)
				}
			}
		default:
			// Check if the value is an image field
			if key, ok := v.(string); ok && strings.Contains(key, "image") {
				*containers = append(*containers, key)
			}
		}
	}
}
