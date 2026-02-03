# Kubernetes Deployment Guide

This guide covers deploying the Bun + Hono + React monorepo application to Kubernetes.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Deployment](#deployment)
- [Management](#management)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

## Prerequisites

1. **Kubernetes Cluster**
   - Minikube (local development)
   - Kind (local development)
   - GKE, EKS, AKS (cloud providers)
   - Self-hosted cluster

2. **Tools**
   - `kubectl` - Kubernetes CLI
   - `make` - Build automation tool
   - `docker` - Container runtime
   - `openssl` - For generating secrets

3. **Docker Images**
   - Build and push images to GitHub Container Registry
   - See [DOCKER.md](./DOCKER.md) for details

4. **Environment Variables**
   ```bash
   export GITHUB_USER=your-github-username
   export GITHUB_TOKEN=your-github-token
   ```

## Quick Start

### 1. Build and Push Docker Images

```bash
# Login to GitHub Container Registry
make login

# Build images
make build-all

# Push images to registry
make push-all

# Or do all at once
make deploy
```

### 2. Configure Kubernetes

```bash
# Copy and edit environment file
cp .env.k8s.example .env.k8s
# Edit .env.k8s with your values
```

### 3. Deploy to Kubernetes

```bash
# Generate secret from environment variables
make k8s-generate-secret

# Deploy all resources
make k8s-deploy

# Apply the secret
make k8s-apply-secret
```

### 4. Check Status

```bash
# Check deployment status
make k8s-status

# Get external IP
make k8s-get-external-ip
```

## Configuration

### Environment Variables

Edit `.env.k8s`:

```bash
GITHUB_USER=your-github-username
GITHUB_TOKEN=ghp_your_token
DATABASE_URL=postgresql://user:pass@host:5432/db
JWT_SECRET=your-jwt-secret
```

### Kubernetes Manifests

All Kubernetes manifests are in the `k8s/` directory:

```
k8s/
├── namespace.yaml              # Namespace definition
├── configmap.yaml              # Non-sensitive configuration
├── secret.yaml.template        # Secret template
├── server-deployment.yaml      # Server deployment
├── server-service.yaml         # Server service
├── client-deployment.yaml      # Client deployment
├── client-service.yaml         # Client service (LoadBalancer)
└── ingress.yaml                # Ingress (optional)
```

### Update Image References

Before deploying, update the image references in deployment files:

```bash
# Automatically update with your GITHUB_USER
make k8s-update-images
```

Or manually edit `k8s/server-deployment.yaml` and `k8s/client-deployment.yaml`:

```yaml
image: ghcr.io/YOUR_GITHUB_USER/bun-hono-react-monorepo-server:latest
```

## Deployment

### Full Deployment Workflow

```bash
# Complete workflow: build, push, and deploy
make k8s-full-deploy
```

### Step-by-Step Deployment

```bash
# 1. Create namespace
make k8s-apply-namespace

# 2. Generate and apply secret
make k8s-generate-secret
make k8s-apply-secret

# 3. Apply ConfigMap
make k8s-apply-configmap

# 4. Deploy server
make k8s-deploy-server

# 5. Deploy client
make k8s-deploy-client

# 6. (Optional) Deploy ingress
make k8s-deploy-ingress
```

### Deploying Individual Components

```bash
# Deploy only server
make k8s-deploy-server

# Deploy only client
make k8s-deploy-client
```

## Management

### Status and Information

```bash
# Check deployment status
make k8s-status

# List pods
make k8s-pods

# List services
make k8s-services

# Describe resources
make k8s-describe

# Show deployment information
make k8s-info

# Check current context
make k8s-check-context
```

### Logs

```bash
# View server logs (last 100 lines)
make k8s-logs-server

# View client logs (last 100 lines)
make k8s-logs-client

# Follow server logs (real-time)
make k8s-logs-server-follow

# Follow client logs (real-time)
make k8s-logs-client-follow

# View all logs
make k8s-logs-all
```

### Restart/Reload

```bash
# Restart server pods
make k8s-reload-server

# Restart client pods
make k8s-reload-client

# Restart all pods
make k8s-reload
```

### Scaling

```bash
# Scale server to 3 replicas
make k8s-scale-server REPLICAS=3

# Scale client to 5 replicas
make k8s-scale-client REPLICAS=5

# Stop all deployments (scale to 0)
make k8s-stop

# Start all deployments (scale to default)
make k8s-start
```

### Updates

```bash
# Update deployments (rolling update)
make k8s-update
```

### Rollback

```bash
# Rollback server deployment
make k8s-rollback-server

# Rollback client deployment
make k8s-rollback-client

# View deployment history
make k8s-history-server
make k8s-history-client
```

### Secrets Management

```bash
# Generate secret from environment variables
make k8s-generate-secret

# Create image pull secret for private registry
make k8s-create-image-pull-secret
```

## Monitoring

### Resource Usage

```bash
# Show pod resource usage
make k8s-top-pods

# Show node resource usage
make k8s-top-nodes
```

### Events

```bash
# View recent events
make k8s-events
```

### Port Forwarding

```bash
# Forward server to localhost:3000
make k8s-port-forward-server

# Forward client to localhost:8080
make k8s-port-forward-client
```

### Execute Commands in Pods

```bash
# Open shell in server pod
make k8s-exec-server

# Open shell in client pod
make k8s-exec-client
```

## Exposing Services

### LoadBalancer (Cloud Providers)

The client service is configured as `LoadBalancer` by default:

```bash
# Get external IP
make k8s-get-external-ip
```

### NodePort (Local/Bare Metal)

Edit `k8s/client-service.yaml`:

```yaml
spec:
  type: NodePort
  ports:
  - port: 80
    targetPort: 80
    nodePort: 30080
```

Access at: `http://<node-ip>:30080`

### Ingress (Recommended for Production)

1. Install ingress controller:

```bash
# NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml
```

2. Update `k8s/ingress.yaml` with your domain

3. Deploy ingress:

```bash
make k8s-deploy-ingress
```

## Troubleshooting

### Pods Not Starting

```bash
# Describe server deployment
make k8s-describe-server

# Describe client deployment
make k8s-describe-client

# View events
make k8s-events
```

### ImagePullBackOff

1. Verify image exists:
   ```bash
   docker pull ghcr.io/YOUR_USER/bun-hono-react-monorepo-server:latest
   ```

2. Create image pull secret if using private registry:
   ```bash
   make k8s-create-image-pull-secret
   ```

3. Uncomment `imagePullSecrets` in deployment files

### CrashLoopBackOff

```bash
# Check logs
make k8s-logs-server
make k8s-logs-client

# Check events
make k8s-events
```

### Service Not Accessible

```bash
# Check service endpoints
kubectl get endpoints -n bun-hono-react

# Test from within cluster
kubectl run test-pod --rm -it --image=curlimages/curl -n bun-hono-react -- sh
# Then: curl http://server-service:3000
```

### Database Connection Issues

1. Verify DATABASE_URL in secret:
   ```bash
   kubectl get secret monorepo-secret -n bun-hono-react -o yaml
   ```

2. Check if database is accessible from cluster

3. Update secret:
   ```bash
   export DATABASE_URL="your-new-database-url"
   make k8s-generate-secret
   make k8s-apply-secret
   make k8s-reload
   ```

## Cleanup

```bash
# Delete all resources (keep namespace)
make k8s-delete

# Delete entire namespace (WARNING: deletes everything)
make k8s-delete-namespace
```

## Production Recommendations

### Security

1. **Use Secrets Management**
   - [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets)
   - [External Secrets Operator](https://external-secrets.io/)
   - Cloud provider secret management (AWS Secrets Manager, GCP Secret Manager, etc.)

2. **Enable RBAC**
   - Create service accounts with minimal permissions
   - Use pod security policies

3. **Network Policies**
   - Restrict pod-to-pod communication
   - Limit egress traffic

### High Availability

1. **Multiple Replicas**
   ```bash
   make k8s-scale-server REPLICAS=3
   make k8s-scale-client REPLICAS=3
   ```

2. **Pod Disruption Budgets**
   - Ensure minimum number of pods during updates

3. **Anti-Affinity Rules**
   - Distribute pods across different nodes

### Resource Management

1. **Set Resource Requests/Limits**
   - Edit deployment files with appropriate values
   - Monitor actual usage and adjust

2. **Horizontal Pod Autoscaling**
   ```yaml
   apiVersion: autoscaling/v2
   kind: HorizontalPodAutoscaler
   metadata:
     name: server-hpa
   spec:
     scaleTargetRef:
       apiVersion: apps/v1
       kind: Deployment
       name: server-deployment
     minReplicas: 2
     maxReplicas: 10
     metrics:
     - type: Resource
       resource:
         name: cpu
         target:
           type: Utilization
           averageUtilization: 70
   ```

### Monitoring and Logging

1. **Prometheus & Grafana**
   - Monitor metrics
   - Set up alerts

2. **ELK/EFK Stack**
   - Centralize logs
   - Create dashboards

3. **Distributed Tracing**
   - Jaeger or Zipkin
   - OpenTelemetry

### TLS/HTTPS

1. **Install cert-manager**
   ```bash
   kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
   ```

2. **Create ClusterIssuer**
   ```yaml
   apiVersion: cert-manager.io/v1
   kind: ClusterIssuer
   metadata:
     name: letsencrypt-prod
   spec:
     acme:
       server: https://acme-v02.api.letsencrypt.org/directory
       email: your-email@example.com
       privateKeySecretRef:
         name: letsencrypt-prod
       solvers:
       - http01:
           ingress:
             class: nginx
   ```

3. **Update ingress** with TLS configuration

## Available Make Commands

Run `make help` to see all available Docker commands.
Run `make k8s-info` to see all available Kubernetes commands.

### Kubernetes Commands Summary

| Command | Description |
|---------|-------------|
| `k8s-deploy` | Deploy all resources |
| `k8s-deploy-server` | Deploy only server |
| `k8s-deploy-client` | Deploy only client |
| `k8s-update` | Update deployments |
| `k8s-status` | Check status |
| `k8s-pods` | List pods |
| `k8s-services` | List services |
| `k8s-logs-server` | View server logs |
| `k8s-logs-client` | View client logs |
| `k8s-logs-server-follow` | Follow server logs |
| `k8s-logs-client-follow` | Follow client logs |
| `k8s-reload-server` | Restart server |
| `k8s-reload-client` | Restart client |
| `k8s-reload` | Restart all |
| `k8s-scale-server` | Scale server |
| `k8s-scale-client` | Scale client |
| `k8s-stop` | Stop all (scale to 0) |
| `k8s-start` | Start all |
| `k8s-delete` | Delete resources |
| `k8s-generate-secret` | Generate secret |
| `k8s-full-deploy` | Build, push, deploy |

For a complete list, see the [Makefile](./Makefile) or run `make help`.

## Further Reading

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
- [Production Best Practices](https://learnk8s.io/production-best-practices)
