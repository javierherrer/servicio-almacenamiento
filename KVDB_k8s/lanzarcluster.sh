#/bin/sh
# Creacion de los binarios
echo "Creando binarios $(pwd)"
CGO_ENABLED=0 go build -o kvdb main.go

echo "Docker builds de todas las imagenes"
docker build . -t localhost:5000/kvdb:latest

echo "Push de las imagenes al registro local"
docker push localhost:5000/kvdb:latest


echo "Borramos los pods y deployments de anteriores ejecuciones"
kubectl delete -f kvdb_stateful.yaml 
for pod in $(kubectl get pods --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')
do
    kubectl delete pod $pod
done

sleep 2

echo "Eliminados los pods anteriores"

kubectl create -f kvdb_stateful.yaml 
