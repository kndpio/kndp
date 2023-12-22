package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/slack-go/slack"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"net/url"
)

type SelectedOptionValue struct {
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Actions []struct {
		Name            string `json:"name"`
		Type            string `json:"type"`
		SelectedOptions []struct {
			Value string `json:"value"`
		} `json:"selected_options"`
	} `json:"actions"`
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

func patchEmployeeRefStatus(user, selectedOption string) error {
	// Read employeeRefs status from Meal object
	cmd := exec.Command("kubectl", "get", "meal", mealName, "-o", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running kubectl command: %v", err)
	}
	var meal Meal
	err = json.Unmarshal(output, &meal)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	foundUser := false
	for i := range meal.Spec.EmployeeRefs {
		if meal.Spec.EmployeeRefs[i].Name == user {
			meal.Spec.EmployeeRefs[i].Status = selectedOption
			foundUser = true
			break
		}
	}

	if !foundUser {
		newEmployeeRef := EmployeeRef{
			Name:   user,
			Status: selectedOption,
		}
		meal.Spec.EmployeeRefs = append(meal.Spec.EmployeeRefs, newEmployeeRef)
	}

	fmt.Println(meal.Spec.EmployeeRefs)

	updatedMealJSON, err := json.Marshal(meal)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	cmd = exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader(updatedMealJSON)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error applying kubectl patch: %v", err)
	}

	fmt.Println(string(output))

	return nil
}
func events() {

	http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
		payload, err := url.QueryUnescape(r.FormValue("payload"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error decoding payload:", err)
			return
		}
		var data SelectedOptionValue
		err = json.Unmarshal([]byte(payload), &data)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		selected_option := data.Actions[0].SelectedOptions[0].Value
		user := data.User.Name
		userID := data.User.ID
		respondMsg(userID, user)

		err = patchEmployeeRefStatus(user, selected_option)
		if err != nil {
			fmt.Println("Error patching EmployeeRef status:", err)
		}
	})
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}

func respondMsg(userID string, userName string) {

	attachment := slack.Attachment{
		Color:         "#f9a41b",
		Fallback:      "",
		CallbackID:    "meal",
		AuthorID:      "",
		AuthorName:    "",
		AuthorSubname: "",
		AuthorLink:    "",
		AuthorIcon:    "",
		Title:         "",
		TitleLink:     "",
		Pretext:       "",
		Text:          "Thank you",
		ImageURL:      "",
		ThumbURL:      "",
		ServiceName:   "",
		ServiceIcon:   "",
		FromURL:       "",
		OriginalURL:   "",
		Fields:        []slack.AttachmentField{},
		Actions:       []slack.AttachmentAction{},
		MarkdownIn:    []string{},
		Blocks:        slack.Blocks{},
		Footer:        "",
		FooterIcon:    "",
		Ts:            "",
	}

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

func main() {
	events()
}
