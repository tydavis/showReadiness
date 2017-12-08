# showReadiness

A simple golang application to explore [liveness and readiness checks](http://kubernetes.io/docs/user-guide/pod-states/) in Kubernetes.

## To build
In the cloned directory, run the following:
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags '-w -s' .
```

## To build with only Docker
Ensure [docker](https://www.docker.com/get-docker) is installed. In the cloned directory, run the folowing: 

```
docker run --rm -v $PWD:/go/src/github.com/tydavis/showreadiness -w /go/src/github.com/tydavis/showreadiness golang:alpine /bin/sh -c "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags '-w -s' . "
```

