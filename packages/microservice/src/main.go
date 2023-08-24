package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func findChartYamlFiles(rootDir string) ([]string, error) {
	var chartYamlFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "Chart.yaml" {
			chartYamlFiles = append(chartYamlFiles, path)
		}
		return nil
	})

	return chartYamlFiles, err
}

func packageChart(chartPath string, outputDir string) error {
	cmd := exec.Command("helm", "package", chartPath, "--destination", filepath.Join(outputDir, "charts"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateRepoIndex(outputDir string) error {
	cmd := exec.Command("helm", "repo", "index", outputDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	rootDir := "/storage"
	outputDir := "./packages"

	chartYamlFiles, err := findChartYamlFiles(rootDir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(chartYamlFiles) == 0 {
		fmt.Println("No Chart.yaml files found.")
		return
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	fmt.Println("Packaging Helm charts...")
	for _, chartPath := range chartYamlFiles {
		fmt.Printf("Packaging %s...\n", chartPath)
		if err := packageChart(filepath.Dir(chartPath), outputDir); err != nil {
			fmt.Printf("Error packaging %s: %v\n", chartPath, err)
		}
	}

	fmt.Println("Generating repository index...")
	if err := generateRepoIndex(outputDir); err != nil {
		fmt.Println("Error generating repository index:", err)
		return
	}

	fmt.Println("Packaging and index generation complete. TGZ files and index stored in", outputDir)

	port := 5000

	http.HandleFunc("/index.yaml", func(w http.ResponseWriter, r *http.Request) {
		indexFilePath := filepath.Join(outputDir, "index.yaml")
		http.ServeFile(w, r, indexFilePath)
	})

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(outputDir))))

	fmt.Printf("Starting server on port %d...\n", port)
	fmt.Printf("To access the repository index, visit http://localhost:%d/index.yaml\n", port)
	fmt.Printf("To access the charts, use URLs like http://localhost:%d/charts/your-chart-version.tgz\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
