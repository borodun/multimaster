# Multimaster metrics

Metrics for Postgres multimaster cluster.

## Usage

You can [compile it yourself](#building-from-source), run in [container](#running-container)
or [Kubernetes](#running-in-kubernetes). . To run it you need to write a config, see [example](config.yaml). In config
you need to specify at least one node that is in the cluster and online. If password isn't present in config then it
will try to get it from ~/.pgpass

#### Restrictions

1. For now, it works only in local network with multimaster, because it needs to discover new nodes that will join the
   cluster
2. For now, name of the database in config should be in format 'node{id}', where id is the id of this node in the
   cluster (sql: SELECT my_node_id FROM mtm.status())

### Building from source

```shell
make
./mtm-metrics
./mtm-metrics -c path/to/config.yaml
```

### Running container

```shell
docker run -p 8080:8080 -v $(pwd)/config.yaml:/home/mmts/config.yaml -v $HOME/.pgpass:/home/mmts/.pgpass -v /etc/hosts:/etc/hosts borodun/mtm-metrics:latest
```

### Running in Kubernetes

```shell
kubectl create namespace mtm
kubectl create secret generic mtm-metrics-pgpass -n mtm --from-file=.pgpass=$HOME/.pgpass
kubectl create configmap mtm-metrics-config -n mtm --from-file=config.yaml=config.yaml
kubectl create configmap mtm-metrics-hosts -n mtm --from-file=hosts=/etc/hosts
kubectl apply -f k8s/mtm-metrics-deployment.yaml -n mtm
```

Metrics are based on [postgres_exporter](https://github.com/ContaAzul/postgresql_exporter)