#!/bin/bash

# Default values
WAIT_TIME=60
BUILD=true
SEED_URL=""
WORKER_COUNT=1000
MAX_URLS=10000
STORAGE_DRIVER="postgres"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --wait-time)
            WAIT_TIME="$2"
            shift
            shift
            ;;
        --build)
            BUILD=true
            shift
            ;;
        --seed-url)
            SEED_URL="$2"
            shift
            shift
            ;;
        --worker-count)
            WORKER_COUNT="$2"
            shift
            shift
            ;;
        --max-urls)
            MAX_URLS="$2"
            shift
            shift
            ;;
        --storage-driver)
            STORAGE_DRIVER="$2"
            shift
            shift
            ;;
        *)
            echo "Unknown parameter passed: $1"
            exit 1
            ;;
    esac
done

# Validate required arguments
if [ -z "$SEED_URL" ]; then
    echo "Error: --seed-url is required."
    exit 1
fi

# Output parsed parameters
echo "Parameters:"
echo "WAIT_TIME: $WAIT_TIME"
echo "BUILD: $BUILD"
echo "SEED_URL: $SEED_URL"
echo "WORKER_COUNT: $WORKER_COUNT"
echo "MAX_URLS: $MAX_URLS"
echo "STORAGE_DRIVER: $STORAGE_DRIVER"

# Step 1: Start Minikube if it's not already running
echo "Starting Minikube..."
if ! minikube status | grep -q "host: Running"; then
    minikube start
fi

# Step 2: Set up Docker environment to build images directly in Minikube
echo "Configuring Docker to use Minikube's Docker environment..."
eval $(minikube docker-env)

# Step 3: Build Docker image if the build flag is set
if $BUILD; then
    echo "Building Docker image for web crawler..."
    docker build -t web-crawler:latest .
fi

# Step 4: Create or update the ConfigMap with dynamic SEED_URL and WORKER_COUNT
echo "Creating or updating Kubernetes ConfigMap with dynamic values..."
kubectl create configmap crawler-config \
    --from-literal=STORAGE_DRIVER=$STORAGE_DRIVER \
    --from-literal=SEED_URL=$SEED_URL \
    --from-literal=WORKER_COUNT=$WORKER_COUNT \
    --from-literal=MAX_URLS=$MAX_URLS \
    --dry-run=client -o yaml | kubectl apply -f -

# Step 5: Deploy Redis
echo "Deploying Redis..."
kubectl apply -f k8s/redis-deployment.yaml

echo "Waiting for Redis to be ready..."
REDIS_STATUS=""
while [ "$REDIS_STATUS" != "Running" ]; do
    sleep 5
    REDIS_STATUS=$(kubectl get pod -l app=redis -o jsonpath="{.items[0].status.phase}")
    echo "Redis status: $REDIS_STATUS"
done
echo "Redis is ready."

# Step 6: Deploy the web crawler
echo "Deploying Web Crawler..."
kubectl apply -f k8s/crawler-deployment.yaml
sleep 5

# Step 7: Wait for the web crawler pods to complete
echo "Waiting for $WAIT_TIME seconds for web crawler pods to complete..."
sleep $WAIT_TIME

# Step 8: Retrieve logs of web crawler pods
POD_NAMES=$(kubectl get pods -l app=web-crawler -o jsonpath="{.items[*].metadata.name}")
for POD_NAME in $POD_NAMES; do
    echo "Printing logs for pod $POD_NAME..."
    kubectl logs "$POD_NAME"
done

# Cleanup
echo "Cleaning up Kubernetes resources..."
kubectl delete -f k8s/crawler-deployment.yaml
kubectl delete -f k8s/redis-deployment.yaml

echo "Deployment completed."
