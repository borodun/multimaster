# Info

API server for connecting with multimaster.

### Usage

```shell
./mtm-connector -u "postgresql://user@node1:5432/mydb?sslmode=disable"
```

### Running in Kubernetes

Create namespace:

```shell
kubectl create namespace mtm
```

If needed, create _pgpass_ and _hosts_ that will be mounted to container:

```shell
kubectl create secret generic mtm-connector-pgpass -n mtm --from-file=.pgpass=$HOME/.pgpass
```

Deploy _mtm-connector_ and node port:

```shell
kubectl apply -f k8s/mtm-connector-deployment.yaml -n mtm
kubectl apply -f k8s/mtm-connector-nodeport.yaml -n mtm
```