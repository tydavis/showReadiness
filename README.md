# showReadiness

A simple golang application to explore [liveness and readiness checks](http://kubernetes.io/docs/user-guide/pod-states/) in Kubernetes.

## To build
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags '-w' .
```
