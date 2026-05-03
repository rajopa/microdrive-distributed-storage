🚪 Microdrive API Gateway

The API Gateway is the entry point for the Microdrive ecosystem. It handles incoming client traffic and orchestrates communication between internal microservices via gRPC.
🛠 Specifications

    Port: 55056 (Internal) / 80 (External via Service)

    Runtime: Debian Bookworm Slim (Optimized Multi-stage Build)

    Orchestration: Kubernetes (Namespace: microdrive)

🚀 Deployment Guide
1. Build the Image

Since your deployment manifest uses imagePullPolicy: Never, you need to build the image locally so your cluster (Minikube/Kind) can see it:
Bash

docker build -t yegor222/microdrive-gateway:latest .

2. Prepare the Namespace

Ensure the microdrive namespace exists in your cluster:
Bash

kubectl create namespace microdrive

3. Apply Manifests

Deploy the Gateway and its Service:
Bash

kubectl apply -f deployment.yaml

4. Accessing the Gateway

The service is configured as a LoadBalancer.

On a local Arch Linux setup (Minikube):
Expose the service to get an external IP:
Bash

minikube service gateway-service -n microdrive

Alternative (Port Forwarding):
If you want to access it directly on localhost:8080:
Bash

kubectl port-forward svc/gateway-service 8080:80 -n microdrive

📦 Resource Management

The service is lightweight and production-ready with the following limits:

    Memory: 64Mi (request) / 128Mi (limit)

    CPU: 100m (request) / 200m (limit)