package client

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err.Error())
    }

    clientSet, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    pods, err := clientSet.CoreV1().Pods("namespace").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }

		for _, pod := range pods.Items {
			fmt.Printf("Pod Name: %s, Status: %s\n", pod.Name, pod.Status.Phase)
		}
}
