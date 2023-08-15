package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// applyCluster applies the cluster configuration YAML using `kubectl`.
func applyCluster(filePath string) {
	cmd := exec.Command("kubectl", "apply", "-f", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error applying cluster: %s\n", err)
	}
}
// main is the entry point of the application.
// It watches the ".platform/clusters" directory for changes and applies the YAML files using `kubectl`.

func main() {
	clustersDirectory := os.Getenv("CLUSTERS_DIRECTORY")
	if clustersDirectory == "" {
		clustersDirectory = ".platform/clusters"
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error creating watcher: %s\n", err)
		os.Exit(1)
	}
	defer watcher.Close()

	err = watcher.Add(clustersDirectory)
	if err != nil {
		fmt.Printf("Error watching directory: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Cluster Deployment Operator is running...")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
				// Cluster configuration file changed, apply it using `kubectl`.
				applyCluster(filepath.Join(clustersDirectory, filepath.Base(event.Name)))
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Error watching directory: %s\n", err)
		}
	}
}