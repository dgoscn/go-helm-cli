package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dgoscn/go-helm-cli/internal/charts"
	"github.com/dgoscn/go-helm-cli/internal/repository"
	"github.com/dgoscn/go-helm-cli/internal/kube"
	"github.com/dgoscn/go-helm-cli/internal/common/utils"
)

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
		err := charts.AddChart(chartLocation)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Chart added successfully.")

	case "index":
		indexCmd.Parse(os.Args[2:])
		err := repository.GenerateRepoIndex()
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
		err := kube.InstallChart(chartName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Chart '%s' installed successfully.\n", chartName)

	case "images":
		imagesCmd.Parse(os.Args[2:])
		imageList, err := charts.GetImageList()
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
