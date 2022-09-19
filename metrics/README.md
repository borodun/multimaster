# Multimaster metrics

Metrics for Postgres multimaster cluster.
Metrics are based on [postgres_exporter](https://github.com/ContaAzul/postgresql_exporter)

## Usage

You can [compile it yourself](#building-from-source), run in [container](#running-container)
or [Kubernetes](#running-in-kubernetes). To run it you need to write a config, see [example](config.yaml). Specify at
least one node that is in the cluster and online. If password isn't present in config then it will try to get it
from ~/.pgpass. Metrics will be available on [localhost:8080/metrics](http://localhost:8080/metrics)

#### Current limitations

1. For now, it works only in local network with multimaster, because it needs to discover other nodes that are in the
   cluster
2. To avoid naming conflicts, name of the database in config should be in format 'node{**id**}' (example: node**1**),
   where id is the id of this node in the cluster (sql: SELECT my_node_id FROM mtm.status())

### Building from source

```shell
make
./mtm-metrics
```

### Running container

```shell
docker run -p 8080:8080 -v $(pwd)/config.yaml:/home/mmts/config.yaml -v $HOME/.pgpass:/home/mmts/.pgpass -v /etc/hosts:/etc/hosts borodun/mtm-metrics:latest
```

### Running in Kubernetes

You need to have K8s with Prometheus and Grafana,
see [how to make one](https://github.com/borodun/k8s-manifests#bare-metal-kubernetes-for-working). You need to
install [Statusmap](https://grafana.com/grafana/plugins/flant-statusmap-panel/) plugin for Grafana.

Create namespace and config that will be mounted to container:

```shell
kubectl create namespace mtm
kubectl create configmap mtm-metrics-config -n mtm --from-file=config.yaml=config.yaml
```

If needed, create _pgpass_ and _hosts_ that will be mounted to container:

```shell
kubectl create secret generic mtm-metrics-pgpass -n mtm --from-file=.pgpass=$HOME/.pgpass
kubectl create configmap mtm-metrics-hosts -n mtm --from-file=hosts=/etc/hosts
```

Deploy _mtm-metrics_:

```shell
kubectl apply -f k8s/mtm-metrics-deployment.yaml -n mtm
```

Deploy PodMonitor for _mtm-metrics_ if you
installed [Prometheus operator](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack):

```shell
kubectl apply -f k8s/mtm-metrics-pod-monitor.yaml -n mtm
```

For additional configuration
see [PodMonitor docs](https://docs.openshift.com/container-platform/4.11/rest_api/monitoring_apis/podmonitor-monitoring-coreos-com-v1.html)

Open Grafana in browser and import [dashboard](grafana/nodes.json)