package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Meal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              struct {
		DeliveryTime string        `json:"deliveryTime"`
		DueOrderTime string        `json:"dueOrderTime"`
		DueTakeTime  string        `json:"dueTakeTime"`
		EmployeeRefs []EmployeeRef `json:"employeeRefs"`
		Status       string        `json:"status"`
	} `json:"spec"`
}

type EmployeeRef struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Unstructured struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
	Spec       map[string]interface{} `json:"spec"`
}

var api = slack.New(os.Getenv("SLACK_API_TOKEN"))

func countUsers(employeeRefs []EmployeeRef) int {
	count := 0
	for _, employeeRef := range employeeRefs {
		if strings.ToLower(employeeRef.Status) == "yes" {
			count++
		}
	}
	return count
}

func readMealResource(dynamicClient dynamic.Interface, ctx context.Context) ([]EmployeeRef, error) {
	resourceId := schema.GroupVersionResource{
		Group:    "kndp.io",
		Version:  "v1alpha1",
		Resource: "meals",
	}

	list, err := dynamicClient.Resource(resourceId).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting resources: %v", err)
	}

	for _, item := range list.Items {
		meal := Meal{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &meal); err != nil {
			return nil, fmt.Errorf("error converting Unstructured to Meal: %v", err)
		}
		return meal.Spec.EmployeeRefs, nil
	}
	return nil, fmt.Errorf("no resources found")
}

func main() {
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamicClient := dynamic.NewForConfigOrDie(config)

	employeeRefs, err := readMealResource(dynamicClient, ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	textContent := countUsers(employeeRefs)
	fmt.Println(textContent)

	attachment := slack.Attachment{
		Color:         "#f9a41b",
		Fallback:      "",
		CallbackID:    "meal",
		AuthorID:      "",
		AuthorName:    "",
		AuthorSubname: "",
		AuthorLink:    "",
		AuthorIcon:    "",
		Title:         strconv.Itoa(textContent),
		TitleLink:     "Meals",
		Pretext:       "",
		Text:          os.Getenv("SLACK_NOTIFY_MESSAGE"),
		ImageURL:      "",
		ThumbURL:      "",
		MarkdownIn:    []string{},
		Footer:        "",
		FooterIcon:    "",
		Ts:            "",
	}

	channelID, timestamp, err := api.PostMessage(
		os.Getenv("SLACK_ORDER_CHANNEL_ID"),
		slack.MsgOptionText("Ordered meals for today:", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
