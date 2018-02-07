# showReadiness

A simple golang application to explore [liveness and readiness checks][1] in
Kubernetes.

## Explanation

The `showreadiness` program will respond to requests on port 80 locally, port
6080 externally, with text indicating it is handling traffic. With that
response, there is a header called `responding-pod` which contains the name of
the pod (its hostname).

Calling `/makeNotReady` will disable traffic processing for whichever host
handled the request. Using this endpoint, we can stop traffic to each pod, and
then restore it using the `/makeReady` endpoint directly (via accessing the pod
over `localhost`).

## To build

In the cloned directory, run the following:

### Linux/OSX

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags '-w -s' .
```

### Windows

```powershell
$env:CGO_ENABLED=0; $env:GOOS="linux"; $env:GOARCH="amd64";
go build -a -tags netgo -installsuffix netgo -ldflags "-s"
```

Don't forget to unset the variables when you're done!

```powershell
Remove-Item env:CGO_ENABLED; Remove-Item env:GOOS; Remove-Item env:GOARCH;
```

### Using only Docker (no Go compiler)

Assuming [docker][2] is installed, run the folowing in the checked-out directory:

```bash
docker run --rm -v $PWD:/go/src/github.com/tydavis/showreadiness \
-w /go/src/github.com/tydavis/showreadiness golang:alpine /bin/sh \
-c "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags '-w -s' . "
``` 

[1]:http://kubernetes.io/docs/user-guide/pod-states/
[2]:https://www.docker.com/get-docker