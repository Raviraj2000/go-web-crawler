param (
    [int]$WaitTime = 60, # Default wait time in seconds if no argument is provided
    [bool]$build = $false, # Default value for the build parameter
    [string]$SeedUrl = "", # Seed URL (required)
    [string]$WorkerCount = "1000", # Default worker count
    [string]$MaxUrls = "10000",
    [string]$StorageDriver = "postgres"
)

Write-Output "Scraping $MaxUrls urls ...."

# Check if SeedUrl is provided
if ([string]::IsNullOrEmpty($SeedUrl)) {
    Write-Output "Error: Seed URL is required. Please provide a valid Seed URL with the -SeedUrl parameter."
    exit 1
}

# Step 1: Start Minikube if it's not already running
Write-Output "Starting Minikube..."
if (-not (minikube status | Select-String "host: Running")) {
    minikube start
}

# Step 2: Set up Docker environment to build images directly in Minikube
Write-Output "Configuring Docker to use Minikube's Docker environment..."
& minikube -p minikube docker-env --shell powershell | Invoke-Expression

if ($build) {
    Write-Output "Building Docker image for web crawler..."
    try {
        docker build -t web-crawler:latest .
        Write-Output "Docker image built successfully."
    }
    catch {
        Write-Output "Error building Docker image: $_"
        exit 1
    }
}

# Step 3: Create or update the ConfigMap with dynamic SEED_URL and WORKER_COUNT
Write-Output "Creating or updating Kubernetes ConfigMap with dynamic values..."
kubectl create configmap crawler-config --from-literal=STORAGE_DRIVER=$StorageDriver --from-literal=SEED_URL=$SeedUrl --from-literal=WORKER_COUNT=$WorkerCount --from-literal=MAX_URLS=$MaxUrls --dry-run=client -o yaml | kubectl apply -f -

# Step 4: Deploy Redis
Write-Output "Deploying Redis..."
kubectl apply -f k8s/redis-deployment.yaml

Write-Output "Waiting for Redis to be ready..."
$redisStatus = ""
while ($redisStatus -ne "Running") {
    Start-Sleep -Seconds 5
    $redisStatus = kubectl get pod -l app=redis -o jsonpath="{.items[0].status.phase}"
    Write-Output "Redis status: $redisStatus"
}

Write-Output "Redis is ready. Proceeding to Postgres setup."
kubectl apply -f k8s/postgres-service.yaml
# Step 5: Set up port-forwarding for PostgreSQL
Write-Output "Setting up port-forwarding for PostgreSQL..."
Start-Process -NoNewWindow -FilePath "kubectl" -ArgumentList "port-forward svc/postgres-service 5432:5432"
Start-Sleep -Seconds 5  # Wait a moment to ensure port-forwarding is ready

# Step 6: Deploy the web crawler
Write-Output "Deploying Web Crawler..."
kubectl apply -f k8s/crawler-deployment.yaml
Start-Sleep -Seconds 5

$POD_NAMES = kubectl get pods -l app=web-crawler -o jsonpath="{.items[*].metadata.name}" | ForEach-Object { $_ -split " " }

# Step 7: Wait for the web crawler pods to complete
Write-Output "Waiting for $WaitTime seconds for web crawler pods to complete..."
Start-Sleep -Seconds $WaitTime

# Cleanup
Write-Output "Cleaning up Kubernetes resources..."
kubectl delete -f k8s/crawler-deployment.yaml
kubectl delete -f k8s/redis-deployment.yaml
# Stop the port-forwarding process
Write-Output "Stopping port-forwarding for PostgreSQL..."
$portForwardProcess = Get-Process kubectl -ErrorAction SilentlyContinue | Where-Object { $_.StartInfo.Arguments -match "port-forward svc/postgres-service 5432:5432" }
if ($portForwardProcess) {
    $portForwardProcess | Stop-Process
    Write-Output "Port-forwarding stopped."
} else {
    Write-Output "No port-forwarding process found."
}
Write-Output "Deployment completed."
