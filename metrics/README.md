# Multimaster metrics

Usage:
```shell
make
./mtm-metrics
./mtm-metrics --help
```

[Config example](config.yaml)

### Building container

```shell
docker build . -t borodun/mmts-toolbox
```

### Running container

```shell
docker run -p 8080:8080 -v $(pwd)/test/test-config.yaml:/home/mmts/config.yaml -v $HOME/.pgpass:/home/mmts/.pgpass borodun/mtm-metrics:latest
```

### Running in Kubernetes

```shell
kubectl create namespace mmts
kubectl create secret generic mtm-metrics-pgpass -n mmts --from-file=.pgpass=$HOME/.pgpass
kubectl create secret generic mtm-metrics-config -n mmts --from-file=config.yaml=config.yaml
kubectl create secret generic mtm-metrics-hosts -n mmts --from-file=hosts=/etc/hosts
kubectl apply -f k8s/mtm-metrics-deployment.yaml -n mmts
```

Metrics are based on [postgres_exporter](https://github.com/ContaAzul/postgresql_exporter)