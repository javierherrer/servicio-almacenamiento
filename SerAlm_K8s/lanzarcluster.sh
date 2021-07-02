#/bin/sh
# Creacion de los binarios
raiz=$(pwd)
docker="$raiz/docker"
fuentes="$raiz/replicasys"

echo "Creando binarios $(pwd)"
cd $fuentes && ./crearejecutables.sh && cd $raiz


echo "Docker builds de todas las imagenes"
cd $docker/Gestor && docker build . -t localhost:5000/gestorvistas:latest
cd $docker/Test && docker build . -t localhost:5000/testalm:latest
cd $docker/Almacenamiento && docker build . -t localhost:5000/seralm:latest
cd $raiz

echo "Push de las imagenes al registro local"
docker push localhost:5000/gestorvistas:latest
docker push localhost:5000/testalm:latest
docker push localhost:5000/seralm:latest


echo "Borramos los pods y deployments de anteriores ejecuciones"
kubectl delete -f replicasys.yaml
for pod in $(kubectl get pods --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')
do
    kubectl delete pod $pod
done
#kubectl delete pod gestorvistas
#kubectl delete pod testalm
#kubectl delete deployment seralmreplica

#kubectl delete service prueba
sleep 2

echo "Eliminados los pods anteriores"

kubectl create -f replicasys.yaml
