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

func applyManifest(content []byte, delete bool) error {
	if delete {
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

func getStatus(applicationName string) string {
	cmd := exec.Command("kubectl", "describe", "applicationset", applicationName)
	output, err := cmd.Output()
	if err != nil {
		return "Not Activated" 
	}

	if strings.Contains(string(output), "Status:                True") {
		return "Activated"
	}

	return "Not Activated"
}







func main() {
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
		http.HandleFunc("/api/objects", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				fileNames := getFilesInDirectory("./.platform/stacks", ".yaml")
				var response []map[string]string
	
				for _, fileName := range fileNames {
					id := strings.TrimSuffix(fileName, ".yaml")
					filePath := filepath.Join("./.platform/stacks", fileName)
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

				filePath := filepath.Join("./.platform/stacks", reqBody.ID+".yaml")
				content, err := ioutil.ReadFile(filePath)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Error reading file: %s", err)
					return
				}

				if reqBody.Activate {
					err = applyManifest(content, false)
				} else {
					err = applyManifest(content, true)
				}

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error applying: %s", err)
					return
					}

					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, "Deactivated")
				} else {
					w.WriteHeader(http.StatusMethodNotAllowed)
					fmt.Fprint(w, "Unsupported method")
				}
			})

			corsWrappedHandler := corsHandler(http.DefaultServeMux)

			if err := http.ListenAndServe(":8080", corsWrappedHandler); err != nil {
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

		folderPath := "./.platform/stacks"
		err = watcher.Add(folderPath)
		if err != nil {
			fmt.Printf("Error watching folder: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Operator is running...")

		terminate := make(chan os.Signal, 1)
		signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
		<-terminate

		fmt.Println("Received termination signal. Stopping the operator.")
	}

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
