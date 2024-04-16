# Generating OpenAPI Spec for Custom Resource Definitions

This Go program initializes a Kubernetes client and retrieves Custom Resource Definitions (CRDs) specified in the `crds` slice. For each CRD, it fetches the schema for a specific version (e.g., v1alpha1) and generates an OpenAPI specification in YAML format. The generated spec includes the basic information required for the API, such as title, version, paths, and component schemas.

## Prerequisites
- Go programming environment
- Kubernetes cluster access
- `KUBECONFIG` environment variable set

## Setup
1. Set the `KUBECONFIG` environment variable to specify the kubeconfig file path.
2. Run the program to fetch CRDs, generate OpenAPI specs, and save them to YAML files.

## Usage
1. Clone the repository and navigate to the project directory.
2. Build and run the program using the following command:
   ```
   go run main.go
   ```

## Generated Output
For each CRD specified in the `crds` slice, an OpenAPI spec YAML file is created in the format `<CRD Name>_openapi_spec.yaml`.
The file contains the OpenAPI spec with the schema details for the specified version of the CRD.

## Note
- Ensure that the necessary RBAC permissions are granted to the Kubernetes client for fetching CRD information.
- Customize the program to handle errors, logging, and additional functionalities as needed.

--- 

Feel free to reach out if you have any questions or need further assistance with the OpenAPI spec generation for CRDs.
