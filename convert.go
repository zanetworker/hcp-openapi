package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	yaml "sigs.k8s.io/yaml"
)

func InitializeKubernetesClientForExternalUse() (*kubernetes.Clientset, *apiextensionsclientset.Clientset, error) {
	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	apiextensionsClientset, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return clientset, apiextensionsClientset, nil
}
func main() {
	_, extensionsClientSet, err := InitializeKubernetesClientForExternalUse()
	if err != nil {
		fmt.Println("Failed to initialize Kubernetes client:", err)
		os.Exit(1)
	}

	crds := []string{
		"hostedclusters.hypershift.openshift.io",
		"nodepools.hypershift.openshift.io",
	}

	for _, crdName := range crds {
		crd, err := extensionsClientSet.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), crdName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Failed to fetch CRD: %s\n", err)
			continue
		}

		var schema *apiextensionsv1.JSONSchemaProps
		for _, version := range crd.Spec.Versions {
			if version.Name == "v1alpha1" {
				schema = version.Schema.OpenAPIV3Schema
				break
			}
		}

		if schema == nil {
			fmt.Printf("Schema not found for the specified version of CRD: %s\n", crdName)
			continue
		}

		// Prepare OpenAPI spec
		openAPISpec := map[string]interface{}{
			"openapi": "3.0.0",
			"info": map[string]string{
				"title":   "Generated API for CRD",
				"version": "1.0.0",
			},
			"paths": map[string]interface{}{},
			"components": map[string]interface{}{
				"schemas": map[string]interface{}{
					crdName: schema,
				},
			},
		}

		// Serialize to YAML
		data, err := yaml.Marshal(openAPISpec)
		if err != nil {
			fmt.Printf("Failed to marshal OpenAPI spec: %s\n", err)
			continue
		}

		// Write to a file
		fileName := fmt.Sprintf("%s_openapi_spec.yaml", crdName)
		err = os.WriteFile(fileName, data, 0644)
		if err != nil {
			fmt.Printf("Failed to write OpenAPI spec to file: %s\n", err)
			continue
		}

		fmt.Printf("OpenAPI spec for %s created successfully.\n", crdName)
	}
}
