# showReadiness

A simple golang application to explore [liveness and readiness checks][1] in
Kubernetes.

## To build

In the cloned directory, run the following:

Linux/OSX:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags '-w -s' .
```

Windows:

```powershell
$env:CGO_ENABLED=0; $env:GOOS="linux"; $env:GOARCH="amd64"; go build -a -tags netgo -installsuffix netgo -ldflags "-s"
```

Or type your

## To build with only Docker

If [docker][2] is installed, run the folowing in the checked-out directory:

```bash
docker run --rm -v $PWD:/go/src/github.com/tydavis/showreadiness \
-w /go/src/github.com/tydavis/showreadiness golang:alpine /bin/sh \
-c "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags '-w -s' . "
```

[1]:http://kubernetes.io/docs/user-guide/pod-states/
[2]:https://www.docker.com/get-docker