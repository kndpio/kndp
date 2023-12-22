package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/slack-go/slack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	slackBot()
}

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
var mealName = os.Getenv("MEAL_NAME")

func slackBot() {
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

	employeeRefs, err := readEmployeeRefs()
	if err != nil {
		fmt.Printf("Error reading employee references: %s\n", err)
		return
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
func readEmployeeRefs() ([]EmployeeRef, error) {
	// Read employeeRefs status from Meal object
	cmd := exec.Command("/usr/local/bin/kubectl", "get", "meal", mealName, "-o", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running kubectl command: %v", err)
	}

	var meal Meal
	err = json.Unmarshal(output, &meal)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	fmt.Println(meal.Spec.EmployeeRefs)
	return meal.Spec.EmployeeRefs, nil
}

func shouldSendMessage(employeeRefs []EmployeeRef, userName string) bool {
	// Check if user should receive a message based on status
	for _, employeeRef := range employeeRefs {
		if employeeRef.Name == userName {
			return employeeRef.Status == ""
		}
	}
	return true
}
