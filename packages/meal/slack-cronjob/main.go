package main

import (
	"context"
	"fmt"
	"log"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/slack-go/slack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type EmployeeRef struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

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

var api = slack.New(os.Getenv("SLACK_API_TOKEN"))

func readEmployeeRefs(dynamicClient dynamic.Interface, ctx context.Context, group, version, resource, namespace string) ([]EmployeeRef, error) {
	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	list, err := dynamicClient.Resource(resourceId).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting resources: %v", err)
	}

	var employeeRefs []EmployeeRef
	for _, item := range list.Items {
		meal := Meal{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &meal); err != nil {
			return nil, fmt.Errorf("error converting Unstructured to Meal: %v", err)
		}
		employeeRefs = append(employeeRefs, meal.Spec.EmployeeRefs...)
	}
	fmt.Println(employeeRefs)
	return employeeRefs, nil
}

func shouldSendMessage(employeeRefs []EmployeeRef, userName string) bool {
	// Check if the user should receive a message based on status
	for _, employeeRef := range employeeRefs {
		if employeeRef.Name == userName {
			return employeeRef.Status == ""
		}
	}
	return true
}

func main() {
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamicClient := dynamic.NewForConfigOrDie(config)

	employeeRefs, err := readEmployeeRefs(dynamicClient, ctx, "kndp.io", "v1alpha1", "meals", "")
	if err != nil {
		fmt.Printf("Error reading employee references: %s\n", err)
		return
	}

	options := []slack.GetUsersOption{
		slack.GetUsersOptionLimit(1000),
	}
	users, err := api.GetUsers(options...)
	if err != nil {
		log.Fatalf("Error fetching users: %s", err)
	}

	attachment := slack.Attachment{
		Color:         "#f9a41b",
		Fallback:      "",
		CallbackID:    "meal",
		AuthorID:      "",
		AuthorName:    "",
		AuthorSubname: "",
		AuthorLink:    "",
		AuthorIcon:    "",
		Title:         "Meals",
		TitleLink:     "Meals",
		Pretext:       "",
		Text:          "Do you want to order meal today? :poultry_leg:",
		ImageURL:      "",
		ThumbURL:      "",
		ServiceName:   "",
		ServiceIcon:   "",
		FromURL:       "",
		OriginalURL:   "",
		Fields:        []slack.AttachmentField{},
		Actions:       []slack.AttachmentAction{{Name: "actionSelect", Type: "select", Options: []slack.AttachmentActionOption{{Text: "Yes", Value: "Yes"}, {Text: "No", Value: "No"}}}, {Name: "actionCancel", Text: "Cancel", Type: "button", Style: "danger"}},
		MarkdownIn:    []string{},
		Blocks:        slack.Blocks{},
		Footer:        "",
		FooterIcon:    "",
		Ts:            "",
	}

	for _, user := range users {
		userID := user.ID
		userName := user.Name

		if shouldSendMessage(employeeRefs, userName) {
			channelID, _, err := api.PostMessage(
				userID,
				slack.MsgOptionText("", true),
				slack.MsgOptionAttachments(attachment),
				slack.MsgOptionAsUser(true),
			)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			fmt.Printf("Message sent to user %s (%s) in channel %s\n", userName, userID, channelID)
		}
	}
}
