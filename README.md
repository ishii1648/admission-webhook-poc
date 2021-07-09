# admission-webhook-poc

admission-webhook-poc is admission webhook to inject web server(nginx) as sidecar.
nginx reverse proxy to backend(php-fpm).

## Usage

### 1. start cluster

```
$ kind create cluster
```

### 2. install cert manager

```
$ kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml
```

### 3. deploy admission webhook

```
$ make docker-build
$ make deploy
```

### 4. deploy sample manifests

Let's deploy php-fpm as deployment, then admission webhook automatically inject nginx as sidecar and nginx reverse proxy to php-fpm.

```
$ kubectl apply -f sample/
```

### 5. port forward to webapp

```
kubectl port-forward service/webapp 8080:80
```
