package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"operator/controller"

	"github.com/fsnotify/fsnotify"
)

func main() {
	controllerInstance := controller.NewController()

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

				objects := controllerInstance.GetObjects()

				jsonData, err := json.Marshal(objects)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error marshaling data: %s", err)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonData)
			} else if r.Method == http.MethodPost {

				var reqBody struct {
					ID string `json:"id"`
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

				manifest := controller.Manifest{
					Content: content,
				}

				if err := controllerInstance.ApplyManifest(manifest); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error applying manifest: %s", err)
					return
				}

				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "Manifest applied successfully")
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprint(w, "Unsupported method")
			}
		})

		corsWrappedHandler := corsHandler(http.DefaultServeMux)

		// Start the HTTP server with the CORS middleware
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

	// Define the folder path to watch
	folderPath := "./.platform/stacks"
	err = watcher.Add(folderPath)
	if err != nil {
		fmt.Printf("Error watching folder: %s\n", err)
		os.Exit(1)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					// Process the ApplicationSet manifest in the file
					filePath := filepath.Join(folderPath, event.Name)
					if err := controllerInstance.ProcessApplicationSetFile(filePath); err != nil {
						fmt.Printf("Error processing file %s: %s\n", event.Name, err)
					} else {
						fmt.Printf("Processed file: %s\n", event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("Watcher error: %s\n", err)
			}
		}
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	<-terminate

}
