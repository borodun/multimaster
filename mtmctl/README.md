# Toolbox for Postgres multimaster

Usage:
```shell
make
./mtmctl <cmd>
./mtmctl --help
```

[Config example](config.yaml)

### Building container

```shell
docker build . -t borodun/mtmctl

### Running container

```shell
docker run -p 2222:22 -it -v $HOME/.ssh/:/home/mmts/.ssh/ -v $HOME/.pgpass:/home/mmts/.pgpass -v $(pwd)/config.yaml:/home/mmts/config.yaml -v /etc/hosts:/etc/hosts borodun/mtmctl:latest sh
# inside container
~ $ ./mtmctl --help
```

### Running in Kubernetes

```shell
kubectl create namespace mmts
kubectl create secret generic mtmctl-pgpass -n mmts --from-file=.pgpass=$HOME/.pgpass
kubectl create secret generic mtmctl-ssh-keys -n mmts --from-file=id_rsa=$HOME/.ssh/id_rsa --from-file=id_rsa.pub=$HOME/.ssh/id_rsa.pub 
kubectl create secret generic mtmctl-config -n mmts --from-file=config.yaml=config.yaml
kubectl create secret generic mtmctl-hosts -n mmts --from-file=hosts=/etc/hosts
kubectl apply -f k8s/mtmctl-deployment.yaml -n mmts
```