package controller

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

type Manifest struct {
	Content []byte `json:"content"`
}

type MyAppController struct {
}

func NewController() *MyAppController {
	return &MyAppController{}
}

func (c *MyAppController) GetObjects() []string {
	folderPath := "./.platform/stacks"
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("Error reading folder: %s\n", err)
		return []string{}
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames
}

func (c *MyAppController) ApplyManifest(manifest Manifest) error {
	decode := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(manifest.Content), 4096)
	var object unstructured.Unstructured
	if err := decode.Decode(&object); err != nil {
		return fmt.Errorf("error decoding file manifest: %v", err)
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return fmt.Errorf("error creating Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		pod := &corev1.Pod{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, pod); err != nil {
			return fmt.Errorf("error converting to Pod: %v", err)
		}

		_, err := clientset.CoreV1().Pods(object.GetNamespace()).Create(context.TODO(), pod, metav1.CreateOptions{})
		return err
	})

	if retryErr != nil {
		return fmt.Errorf("apply failed: %v", retryErr)
	}
	return nil
}

func (c *MyAppController) ProcessApplicationSetFile(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	yamlDocs := bytes.Split(content, []byte("---"))

	for _, doc := range yamlDocs {
		if len(bytes.TrimSpace(doc)) == 0 {
			continue
		}

		manifest := Manifest{
			Content: doc,
		}

		if err := c.ApplyManifest(manifest); err != nil {
			fmt.Printf("Error applying ApplicationSet manifest: %v\n", err)
	
		} else {
			fmt.Println("Applied ApplicationSet manifest successfully")
		}
	}

	return nil
}
