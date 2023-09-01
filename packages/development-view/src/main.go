package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/rs/cors"
)

type Service struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

func convertToHelmSetFormat(prefix string, values map[string]interface{}, setValues *[]string) {
	for key, value := range values {
		switch v := value.(type) {
		case map[string]interface{}:
			convertToHelmSetFormat(fmt.Sprintf("%s.%s", prefix, key), v, setValues)
		default:
			*setValues = append(*setValues, fmt.Sprintf("%s=%v", key, value))
		}
	}
}

func handleGetChartValues(w http.ResponseWriter, r *http.Request) {
	// Extract the chart name from the URL path
	chartName := strings.TrimPrefix(r.URL.Path, "/api/stacks/")
	chartURL := fmt.Sprintf("http://localhost:5000/charts/%s.tgz", chartName)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Run the helm command to fetch chart metadata (YAML format)
	cmd := exec.Command("helm", "show", "chart", chartURL)
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to fetch chart metadata", http.StatusInternalServerError)
		return
	}

	// Parse YAML chart metadata into JSON
	var chartMetadata map[string]interface{}
	err = yaml.Unmarshal(output, &chartMetadata)
	if err != nil {
		http.Error(w, "Failed to parse chart metadata", http.StatusInternalServerError)
		return
	}

	var chartNameFromMetadata = chartMetadata["name"].(string)
	// Run the helm command to fetch chart values with the custom path
	cmd = exec.Command("helm", "show", "values", chartURL, "--jsonpath", "{}")
	output, err = cmd.Output()
	if err != nil {
		http.Error(w, "Failed to fetch chart values", http.StatusInternalServerError)
		return
	}

	// Parse Helm values JSON
	var helmValues map[string]interface{}
	err = json.Unmarshal(output, &helmValues)
	if err != nil {
		http.Error(w, "Failed to parse Helm values", http.StatusInternalServerError)
		return
	}

	// Construct the response in the desired format
	var services []Service
	service := Service{
		Name:       chartNameFromMetadata,
		Parameters: helmValues,
	}
	services = append(services, service)

	response := map[string][]Service{
		"services": services,
	}

	// Write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func handleUpdateChartValues(w http.ResponseWriter, r *http.Request) {
	// Extract the chart name from the URL path
	chartName := strings.TrimPrefix(r.URL.Path, "/api/update/")
	chartURL := fmt.Sprintf("http://localhost:5000/charts/%s.tgz", chartName)

	// Run the helm command to fetch chart metadata (YAML format)
	cmd := exec.Command("helm", "show", "chart", chartURL)
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to fetch chart metadata", http.StatusInternalServerError)
		return
	}

	// Parse YAML chart metadata into JSON
	var chartMetadata map[string]interface{}
	err = yaml.Unmarshal(output, &chartMetadata)
	if err != nil {
		http.Error(w, "Failed to parse chart metadata", http.StatusInternalServerError)
		return
	}

	chartNameFromMetadata, ok := chartMetadata["name"].(string)
	if !ok {
		http.Error(w, "Failed to retrieve chart name from metadata", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var updatedValues map[string]interface{}
	err = json.Unmarshal(payload, &updatedValues)
	if err != nil {
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}

	parameters, ok := updatedValues["parameters"].(map[string]interface{})
	if !ok {
		http.Error(w, "Failed to retrieve parameters from payload", http.StatusBadRequest)
		return
	}

	var setValues []string
	convertToHelmSetFormat("", parameters, &setValues)

	cmdArgs := []string{"upgrade", "--install", chartNameFromMetadata, chartURL}
	for _, setValue := range setValues {
		cmdArgs = append(cmdArgs, "--set", setValue)
	}

	cmd = exec.Command("helm", cmdArgs...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update chart values: %s", output), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully applied the changes"))
}

func main() {
	corsHandler := cors.Default()

	http.Handle("/api/stacks/", corsHandler.Handler(http.HandlerFunc(handleGetChartValues)))
	http.Handle("/api/update/", corsHandler.Handler(http.HandlerFunc(handleUpdateChartValues)))

	port := 8080
	fmt.Printf("Microservice listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
