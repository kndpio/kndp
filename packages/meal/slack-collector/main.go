package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/slack-go/slack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
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

func handleEventsEndpoint(w http.ResponseWriter, r *http.Request, dynamicClient dynamic.Interface, ctx context.Context) {
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

	selectedOption := data.Actions[0].SelectedOptions[0].Value
	user := data.User.Name
	userID := data.User.ID
	respondMsg(userID, user)

	err = patchEmployeeRefStatus(user, selectedOption, dynamicClient, ctx)
	if err != nil {
		fmt.Println("Error patching EmployeeRef status:", err)
	}
}

func patchEmployeeRefStatus(user, selectedOption string, dynamicClient dynamic.Interface, ctx context.Context) error {
	meal, err := getMealResource(dynamicClient, ctx, mealName)
	if err != nil {
		return fmt.Errorf("error getting Meal resource: %v", err)
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

	obj := createMealUnstructuredObject(meal)
	applyOptions := metav1.ApplyOptions{
		Force:        true,
		FieldManager: "meal-system",
	}

	ApplyResource(ctx, dynamicClient, "kndp.io", "v1alpha1", "meals", "", mealName, obj, applyOptions)
	fmt.Println(err)

	return nil
}

func createMealUnstructuredObject(meal *Meal) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kndp.io/v1alpha1",
			"kind":       "Meal",
			"metadata": map[string]interface{}{
				"name": mealName,
			},
			"spec": map[string]interface{}{
				"employeeRefs": meal.Spec.EmployeeRefs,
			},
		},
	}
}

func getMealResource(dynamicClient dynamic.Interface, ctx context.Context, name string) (*Meal, error) {
	items, err := GetResources(dynamicClient, ctx, "kndp.io", "v1alpha1", "meals", "")
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.GetName() == name {
			meal := &Meal{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, meal); err != nil {
				return nil, fmt.Errorf("error converting Unstructured to Meal: %v", err)
			}
			return meal, nil
		}
	}

	return nil, fmt.Errorf("Meal resource with name %s not found", name)
}

func GetResources(dynamic dynamic.Interface, ctx context.Context, group string, version string, resource string, namespace string) ([]unstructured.Unstructured, error) {
	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	list, err := dynamic.Resource(resourceId).Namespace(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func ApplyResource(ctx context.Context, client dynamic.Interface, group, version, resource, namespace, name string, obj *unstructured.Unstructured, options metav1.ApplyOptions, subresources ...string) (*unstructured.Unstructured, error) {
	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	resourceClient := client.Resource(resourceId).Namespace(namespace)

	return resourceClient.Apply(ctx, name, obj, options, subresources...)
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
		Text:          os.Getenv("SLACK_COLLECTOR_MESSAGE"),
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
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamicClient := dynamic.NewForConfigOrDie(config)
	url := os.Getenv("SLACK_COLLECTOR_URL")
	port := os.Getenv("SLACK_COLLECTOR_PORT")
	if port == "" {
		port = "3000"
	}
	if url == "" {
		url = "/events"
	}
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		handleEventsEndpoint(w, r, dynamicClient, ctx)
	})

	fmt.Println("[INFO] Server listening on port:", port)
	http.ListenAndServe(":"+port, nil)
}
