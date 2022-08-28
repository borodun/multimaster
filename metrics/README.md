# Multimaster metrics

Usage:
```shell
make
./mmts-metrics
./mmts-metrics --help
```

[Config example](config.yaml)

### Building container

```shell
docker build . -t borodun/mmts-toolbox
```

### Running container

```shell
docker run -p 8080:8080 -v $(pwd)/test/test-config.yaml:/home/mmts/config.yaml -v $HOME/.pgpass:/home/mmts/.pgpass borodun/mmts-metrics:latest
```

### Running in Kubernetes

```shell
kubectl create namespace mmts
kubectl create secret generic mmts-metrics-pgpass -n mmts --from-file=.pgpass=$HOME/.pgpass
kubectl create secret generic mmts-metrics-config -n mmts --from-file=config.yaml=config.yaml
kubectl create secret generic mmts-metrics-hosts -n mmts --from-file=hosts=/etc/hosts
kubectl apply -f k8s/mmts-metrics-deployment.yaml -n mmts
```

Metrics are based on [postgres_exporter](https://github.com/ContaAzul/postgresql_exporter)