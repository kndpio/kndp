package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/crossplane/function-sdk-go/errors"
	"github.com/crossplane/function-sdk-go/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/resource/composed"
	"github.com/crossplane/function-sdk-go/response"
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

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer
	log logging.Logger
}

func readMealResource(dynamicClient dynamic.Interface, ctx context.Context) (string, error) {
	resourceId := schema.GroupVersionResource{
		Group:    "kndp.io",
		Version:  "v1alpha1",
		Resource: "meals",
	}

	list, err := dynamicClient.Resource(resourceId).List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("error getting resources: %v", err)
	}

	for _, item := range list.Items {
		meal := Meal{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &meal); err != nil {
			return "", fmt.Errorf("error converting Unstructured to Meal: %v", err)
		}

		dueOrderTime := meal.Spec.DueOrderTime
		fmt.Println(dueOrderTime)
		return dueOrderTime, nil
	}
	return "", fmt.Errorf("no resources found")
}

// RunFunction adds a Deployment and the new object template to the desired state.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	rsp := response.To(req, response.DefaultTTL)

	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamicClient := dynamic.NewForConfigOrDie(config)

	dueOrderTime, err := readMealResource(dynamicClient, ctx)
	if err != nil {
		fmt.Printf("Error reading employee references: %s\n", err)
		return nil, err
	}

	if dueOrderTime == "over" {
		unstructuredData := composed.Unstructured{}
		desired, err := request.GetDesiredComposedResources(req)
		if err != nil {
			response.Fatal(rsp, errors.Wrapf(err, "cannot get desired resources from %T", req))
			return rsp, nil
		}
		desired[resource.Name("")] = &resource.DesiredComposed{Resource: &unstructuredData}

		if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
			response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources in %T", rsp))
			return rsp, nil
		}

		fmt.Println("DueOrderTime is over, return empty object")
	} else {

		desired, err := request.GetDesiredComposedResources(req)
		if err != nil {
			response.Fatal(rsp, errors.Wrapf(err, "cannot get desired resources from %T", req))
			return rsp, nil
		}

		ingressTemplate := map[string]interface{}{
			"apiVersion": "kubernetes.crossplane.io/v1alpha1",
			"kind":       "Object",
			"metadata": map[string]interface{}{
				"name": "ingress-object",
			},
			"spec": map[string]interface{}{
				"forProvider": map[string]interface{}{
					"manifest": map[string]interface{}{
						"apiVersion": "networking.k8s.io/v1",
						"kind":       "Ingress",
						"metadata": map[string]interface{}{
							"name":      "meal-ingress",
							"namespace": "default",
						},
						"spec": map[string]interface{}{
							"rules": []map[string]interface{}{
								{
									"host": os.Getenv("HOST_NAME"),
									"http": map[string]interface{}{
										"paths": []map[string]interface{}{
											{
												"pathType": "Prefix",
												"path":     "/",
												"backend": map[string]interface{}{
													"service": map[string]interface{}{
														"name": "meal-service",
														"port": map[string]interface{}{
															"number": 80,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"managementPolicy":  "Default",
				"providerConfigRef": map[string]interface{}{"name": os.Getenv("PROVIDER_CONFIG_REF_NAME")},
			},
		}

		cronJobTemplate := map[string]interface{}{
			"apiVersion": "kubernetes.crossplane.io/v1alpha1",
			"kind":       "Object",
			"metadata": map[string]interface{}{
				"name": os.Getenv("CRONJOB_NAME"),
			},
			"spec": map[string]interface{}{
				"forProvider": map[string]interface{}{
					"manifest": map[string]interface{}{
						"apiVersion": "batch/v1",
						"kind":       "CronJob",
						"metadata": map[string]interface{}{
							"name":      os.Getenv("CRONJOB_NAME"),
							"namespace": "default",
						},
						"spec": map[string]interface{}{
							"schedule": os.Getenv("CRONJOB_SCHEDULE_TIME"),
							"jobTemplate": map[string]interface{}{
								"spec": map[string]interface{}{
									"template": map[string]interface{}{
										"spec": map[string]interface{}{
											"serviceAccountName": os.Getenv("CRONJOB_SERVICE_ACCOUNT_NAME"),
											"containers": []map[string]interface{}{
												{
													"name":  "meal-container",
													"image": os.Getenv("CRONJOB_IMAGE_NAME"),
													"envFrom": []map[string]interface{}{
														{"configMapRef": map[string]interface{}{"name": "meal-cm"}},
													},
												},
											},
											"restartPolicy": "OnFailure",
										},
									},
								},
							},
						},
					},
				},
				"managementPolicy":  "Default",
				"providerConfigRef": map[string]interface{}{"name": os.Getenv("PROVIDER_CONFIG_REF_NAME")},
			},
		}

		deploymentTemplate := map[string]interface{}{
			"apiVersion": "kubernetes.crossplane.io/v1alpha1",
			"kind":       "Object",
			"metadata": map[string]interface{}{
				"name": os.Getenv("DEPLOYMENT_NAME"),
			},
			"spec": map[string]interface{}{
				"forProvider": map[string]interface{}{
					"manifest": map[string]interface{}{
						"apiVersion": "apps/v1",
						"kind":       "Deployment",
						"metadata": map[string]interface{}{
							"name":      os.Getenv("DEPLOYMENT_NAME"),
							"namespace": "default",
						},
						"spec": map[string]interface{}{
							"replicas": 3,
							"selector": map[string]interface{}{
								"matchLabels": map[string]interface{}{
									"app": "meal",
								},
							},
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"labels": map[string]interface{}{
										"app": "meal",
									},
								},
								"spec": map[string]interface{}{
									"serviceAccountName": os.Getenv("DEPLOYMENT_SERVICE_ACCOUNT_NAME"),
									"containers": []map[string]interface{}{
										{
											"name":  "meal-container",
											"image": os.Getenv("DEPLOYMENT_IMAGE_NAME"),
											"envFrom": []map[string]interface{}{
												{"configMapRef": map[string]interface{}{"name": os.Getenv("CONFIG_MAP_NAME")}},
											},
											"ports": []map[string]interface{}{
												{
													"containerPort": 80,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"managementPolicy":  "Default",
				"providerConfigRef": map[string]interface{}{"name": os.Getenv("PROVIDER_CONFIG_REF_NAME")},
			},
		}

		// List of templates
		templates := []map[string]interface{}{deploymentTemplate, ingressTemplate, cronJobTemplate}

		// Process each template
		for _, template := range templates {
			unstructuredData := composed.Unstructured{}
			unstructuredDataByte, err := json.Marshal(template)
			if err != nil {
				response.Fatal(rsp, errors.Wrapf(err, "error marshaling Unstructured data: %s", err))
				return rsp, nil
			}
			json.Unmarshal(unstructuredDataByte, &unstructuredData)
			desired[resource.Name(template["metadata"].(map[string]interface{})["name"].(string))] = &resource.DesiredComposed{Resource: &unstructuredData}
		}

		if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
			response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources in %T", rsp))
			return rsp, nil
		}

		f.log.Info("Added Deployment, Ingress, and CronJob templates to desired state")

	}
	return rsp, nil
}
