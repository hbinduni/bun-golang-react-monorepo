# Kubernetes Deployment

This directory contains Kubernetes manifests for deploying the Bun + Hono + React monorepo application.

## Directory Structure

```
k8s/
├── namespace.yaml              # Namespace for the application
├── configmap.yaml              # Application configuration
├── secret.yaml.template        # Secret template (DO NOT commit actual values)
├── imagepullsecret.yaml.template # GitHub Container Registry secret template
├── server-deployment.yaml      # Server deployment
├── server-service.yaml         # Server service
├── client-deployment.yaml      # Client deployment
├── client-service.yaml         # Client service (LoadBalancer)
├── ingress.yaml                # Ingress for external access
└── README.md                   # This file
```

## Prerequisites

1. **Kubernetes cluster** (minikube, kind, GKE, EKS, AKS, etc.)
2. **kubectl** configured to access your cluster
3. **Docker images** pushed to GitHub Container Registry (ghcr.io)
4. **Ingress controller** (optional, for ingress.yaml)

## Quick Start

### Using Makefile (Recommended)

```bash
# 1. Generate secrets
make k8s-generate-secret

# 2. Deploy everything
make k8s-deploy

# 3. Check status
make k8s-status

# 4. View logs
make k8s-logs-server
make k8s-logs-client
```

### Manual Deployment

```bash
# 1. Create namespace
kubectl apply -f k8s/namespace.yaml

# 2. Create configmap
kubectl apply -f k8s/configmap.yaml

# 3. Create secret (edit the template first!)
cp k8s/secret.yaml.template k8s/secret.yaml
# Edit k8s/secret.yaml with your actual values
kubectl apply -f k8s/secret.yaml

# 4. Deploy server
kubectl apply -f k8s/server-deployment.yaml
kubectl apply -f k8s/server-service.yaml

# 5. Deploy client
kubectl apply -f k8s/client-deployment.yaml
kubectl apply -f k8s/client-service.yaml

# 6. (Optional) Apply ingress
kubectl apply -f k8s/ingress.yaml
```

## Configuration

### Environment Variables

Edit `k8s/configmap.yaml` to modify application configuration:
- `NODE_ENV`: Node environment (production/development)
- `SERVER_PORT`: Server port
- `VITE_API_URL`: API URL for the client

### Secrets

Edit `k8s/secret.yaml` (created from template) for sensitive data:
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: JWT signing secret
- Other sensitive configuration

### Image Registry

If using private images from GitHub Container Registry:

1. Uncomment `imagePullSecrets` in deployment files
2. Create the image pull secret:
   ```bash
   make k8s-create-image-pull-secret
   ```

### Updating Images

Update the image in deployment files:
```yaml
image: ghcr.io/YOUR_GITHUB_USER/bun-hono-react-monorepo-server:latest
```

## Makefile Commands

### Deployment
- `make k8s-deploy` - Deploy all resources
- `make k8s-deploy-server` - Deploy only server
- `make k8s-deploy-client` - Deploy only client
- `make k8s-update` - Update deployments (rolling update)

### Management
- `make k8s-status` - Check deployment status
- `make k8s-pods` - List all pods
- `make k8s-services` - List all services
- `make k8s-describe` - Describe all resources

### Logs & Monitoring
- `make k8s-logs-server` - View server logs
- `make k8s-logs-client` - View client logs
- `make k8s-logs-server-follow` - Follow server logs
- `make k8s-logs-client-follow` - Follow client logs

### Restart/Reload
- `make k8s-reload-server` - Restart server pods
- `make k8s-reload-client` - Restart client pods
- `make k8s-reload` - Restart all pods

### Scaling
- `make k8s-scale-server REPLICAS=3` - Scale server
- `make k8s-scale-client REPLICAS=3` - Scale client

### Secrets
- `make k8s-generate-secret` - Generate secret from .env
- `make k8s-create-image-pull-secret` - Create GitHub CR secret

### Cleanup
- `make k8s-stop` - Stop all deployments (scale to 0)
- `make k8s-delete` - Delete all resources
- `make k8s-delete-namespace` - Delete entire namespace

## Exposing Services

### LoadBalancer (Cloud Providers)

The client service is configured as `LoadBalancer` by default. On cloud providers (GKE, EKS, AKS), this will provision an external load balancer.

```bash
# Get external IP
kubectl get svc client-service -n bun-hono-react
```

### NodePort (Local/Bare Metal)

Edit `k8s/client-service.yaml`:
```yaml
spec:
  type: NodePort
  ports:
  - port: 80
    targetPort: 80
    nodePort: 30080  # Add this line
```

Access: `http://<node-ip>:30080`

### Ingress (Recommended for Production)

1. Install ingress controller:
   ```bash
   # NGINX Ingress Controller
   kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml
   ```

2. Update `k8s/ingress.yaml` with your domain

3. Apply ingress:
   ```bash
   kubectl apply -f k8s/ingress.yaml
   ```

## Monitoring

### Check Pod Status
```bash
kubectl get pods -n bun-hono-react -w
```

### View Logs
```bash
# Server logs
kubectl logs -f deployment/server-deployment -n bun-hono-react

# Client logs
kubectl logs -f deployment/client-deployment -n bun-hono-react
```

### Execute Commands in Pod
```bash
# Server
kubectl exec -it deployment/server-deployment -n bun-hono-react -- /bin/sh

# Client
kubectl exec -it deployment/client-deployment -n bun-hono-react -- /bin/sh
```

## Troubleshooting

### Pods not starting
```bash
# Check pod events
kubectl describe pod <pod-name> -n bun-hono-react

# Check logs
kubectl logs <pod-name> -n bun-hono-react
```

### ImagePullBackOff
- Verify image exists in registry
- Check image pull secret if using private registry
- Verify image name and tag in deployment

### CrashLoopBackOff
- Check application logs
- Verify environment variables and secrets
- Check resource limits

### Service not accessible
```bash
# Check service endpoints
kubectl get endpoints -n bun-hono-react

# Test from within cluster
kubectl run test-pod --rm -it --image=curlimages/curl -n bun-hono-react -- sh
# Then: curl http://server-service:3000
```

## Production Recommendations

1. **Resource Limits**: Adjust resource requests/limits based on your workload
2. **Replica Count**: Scale replicas for high availability
3. **Health Checks**: Tune liveness/readiness probe timings
4. **Secrets Management**: Use external secret management (e.g., Sealed Secrets, External Secrets)
5. **TLS/HTTPS**: Enable TLS in ingress with cert-manager
6. **Monitoring**: Set up Prometheus and Grafana
7. **Logging**: Centralize logs with ELK/EFK stack
8. **Autoscaling**: Implement HPA (Horizontal Pod Autoscaler)

## Further Reading

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
