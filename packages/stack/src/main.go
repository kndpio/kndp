package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"sigs.k8s.io/yaml"
)

type Manifest struct {
	Content []byte `json:"content"`
}

type ManifestMetadata struct {
	Name string `yaml:"name"`
}
// applyManifest applies or deletes the Kubernetes manifest content.
// If delete is true, it deletes the resources using `kubectl delete`.
// If delete is false, it applies the resources using `kubectl apply`.
func applyManifest(content []byte, delete bool) error {
	if !delete {
		cmd := exec.Command("kubectl", "delete", "-f", "-")
		cmd.Stdin = bytes.NewReader(content)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to delete resource: %s", err)
		}

		fmt.Println("Deleted manifest successfully")
	} else {
		
		yamlDocs := bytes.Split(content, []byte("---"))

		for _, doc := range yamlDocs {
			if len(bytes.TrimSpace(doc)) == 0 {
				continue
			}

			cmd := exec.Command("kubectl", "apply", "-f", "-")
			cmd.Stdin = bytes.NewReader(doc)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to apply resource: %s", err)
			}

			fmt.Println("Applied manifest successfully")
		}
	}

	return nil
}
// getManifestMetadata extracts metadata information from the manifest content.
// It returns the name, description, and status of the manifest.
func getManifestMetadata(content []byte) (string, string, string) {
	var data map[string]interface{}
	err := yaml.Unmarshal(content, &data)
	if err != nil {
		return "", "", ""
	}

	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		if name, ok := metadata["name"].(string); ok && name != "" {
			if annotations, ok := metadata["annotations"].(map[string]interface{}); ok {
				if description, ok := annotations["description"].(string); ok {
					status := getStatus(name)
					return name, description, status
				}
			}
			return name, "", ""
		}
	}

	return "", "", "" 
}
// getStatus checks the status of an application set using `kubectl describe`.
// It returns "Activated" if the status contains "Status: True", otherwise "Not Activated".
func getStatus(applicationName string) string {
	cmd := exec.Command("kubectl", "describe", "applicationset", applicationName)
	output, err := cmd.Output()
	if err != nil {
		return "Deactivated" 
	}

	if matched, _ := regexp.MatchString(`(?m)^Status:\s*True$`, string(output)); matched {
		return "Activated"
	}

	return "Activated"
}
// getFilesInDirectory returns a list of file names in the specified directory with the given suffix.
func getFilesInDirectory(dirPath, suffix string) []string {
	var fileNames []string

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading folder: %s\n", err)
		return fileNames
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), suffix) {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames
}
// main is the entry point of the application.
// It reads environment variables, starts the HTTP server, and watches the stacks directory for changes.
func main() {
	stacksDirectory := os.Getenv("STACKS_DIRECTORY")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}
	if stacksDirectory == "" {
		stacksDirectory = ".platform/stacks" 
	}
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	go func() {
		http.HandleFunc("/api/stacks", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				fileNames := getFilesInDirectory(stacksDirectory, ".yaml")
				var response []map[string]string
	
				for _, fileName := range fileNames {
					id := strings.TrimSuffix(fileName, ".yaml")
					filePath := filepath.Join(stacksDirectory, fileName)
					content, err := ioutil.ReadFile(filePath)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Error reading file: %s", err)
						return
					}
					name, description, status := getManifestMetadata(content)
					response = append(response, map[string]string{
						"id":   id,
						"name": name,
						"description": description,
						"status": status,
					})
				}
				responseJSON, err := json.Marshal(response)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error marshaling data: %s", err)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(responseJSON)
			} else if r.Method == http.MethodPost {
				var reqBody struct {
					ID       string `json:"id"`
					Activate bool   `json:"activate"`
				}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Error decoding JSON: %s", err)
					return
				}

				filePath := filepath.Join(stacksDirectory, reqBody.ID+".yaml")
				content, err := ioutil.ReadFile(filePath)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Error reading file: %s", err)
					return
				}
				err = applyManifest(content, reqBody.Activate)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error applying: %s", err)
					return
					}
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, "StatusOK")
				} else {
					w.WriteHeader(http.StatusMethodNotAllowed)
					fmt.Fprint(w, "Unsupported method")
				}
			})
			corsWrappedHandler := corsHandler(http.DefaultServeMux)

			if err := http.ListenAndServe(":"+port, corsWrappedHandler); err != nil {
				fmt.Printf("Error starting HTTP server: %s\n", err)
				os.Exit(1)
			}
		}()
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Printf("Error creating watcher: %s\n", err)
			os.Exit(1)
		}
		defer watcher.Close()

		folderPath := stacksDirectory
		err = watcher.Add(folderPath)
		if err != nil {
			fmt.Printf("Error watching folder: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Operator is running on port...", port)
		terminate := make(chan os.Signal, 1)
		signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
		<-terminate
		fmt.Println(" Stopping the operator.")
	}

