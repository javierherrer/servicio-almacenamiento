#!/bin/sh
raiz=$(pwd)
CGO_ENABLED=0 go build -o ./cmd/server/distmutex ./cmd/server/raftserver.go
CGO_ENABLED=0 go build -o ./cmd/testclient/distclient ./cmd/testclient/rafttestclient.go

echo "-----------Docker builds de todas las imagenes"
cd ./cmd/server/ && docker build --no-cache . -t localhost:5000/semser:latest
cd "$raiz"
cd ./cmd/testclient/ && docker build --no-cache . -t localhost:5000/clientetestsem:latest
cd "$raiz"

echo "-----------Push de las imagenes al registro local"
docker push localhost:5000/semser:latest
docker push localhost:5000/clientetestsem:latest

echo "-----------Borramos los pods y deployments de anteriores ejecuciones"
kubectl delete -f distributedsemaphore.yaml
for pod in $(kubectl get pods --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')
do
    kubectl delete pod "$pod"
done

sleep 2

echo "Eliminados los pods anteriores"

kubectl create -f distributedsemaphore.yaml