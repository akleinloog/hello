# hello

A simple HTTP Server in Go that responds with a Hello to any request it gets.  
It includes a counter, the server host name and the relative request URI in the response.

Multiple versions are available to facilitate Kubernetes experiments. The only difference is the language of the greeting.

The server listens for requests on port 80 for requests and in addition, liveness checks and readiness checks are availabe on port 8080.

## Run locally

To start a local instance, simply run
```bash
go run main.go serve
```

## Docker

You can start an instance using docker run, for example:
```bash
docker run -p 10080:80 -p 10880:8080 akleinloog/hello:v1
```
The server will then be available at [http://localhost:10080](http://localhost:10080),
the liveness check at [http://localhost:10880/alive](http://localhost:10880/alive)
and the readiness check at [http://localhost:10880/ready](http://localhost:10880/ready).

To facilitate Kubernetes experiments you can toggle the liveness check by:
[http://localhost:10880/toggleAlive](http://localhost:10880/toggleAlive)

And the readiness check by:
[http://localhost:10880/toggleReady](http://localhost:10880/toggleReady)

To customize the greeting:
```bash
docker run -p 10080:80 -p 10880:8080 -e GREETING='Hi' akleinloog/hello:v4
```

## Docker-compose

You can also use a docker compose file, for example:
```yaml
version: '3.7'

services:
  hello:
    container_name: hello
    image: akleinloog/hello:v1
    ports:
      - 10080:80
      - 10880:8080
```

## Kubernetes Deployment

Deploy it to a Kubernetes cluster using kubectl:
```bash
kubectl create deployment helloapp --image=akleinloog/hello:v1

kubectl expose deployment helloapp --port=80 --type=LoadBalancer
```

Or use a Kubernetes deployment file like this:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hello-app
  template:
    metadata:
      labels:
        app: hello-app
        env: dev
    spec:
      containers:
      - name: hello-app
        image: akleinloog/hello:v1
        imagePullPolicy: Always
        resources:
          limits:
            memory: "128Mi"
            cpu: "200m"
        ports:
        - containerPort: 80
```

Don't forget the service:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-service
spec:
  selector:
    app: hello-app
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
  type: NodePort
```

## Versions

v1: Hello

v2: Bonjour

v3: Aloha

v4: Customizable from Environment Variable GREETING, defaults to Hello

## Version tagging

```
git tag -a v1 -m "Version 1 - English"

git push origin --tags
```