# Instruction for running k8s for monitoring localy

## Minikube

Install minikube, see [docs](https://minikube.sigs.k8s.io/docs/start/):
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```

Run minikube:
```bash
minikube start
```

Stop minikube:
```bash
minikube stop
```

Delete minikube:
```bash
minikube delete --all
```

## Install prometheus stack

Add helm repo
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

Install prometheus stack:
```bash
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack -n monitoring --create-namespace
```

Delete prometheus:
```bash
helm uninstall -n monitoring prometheus
```

After installing, you can add node port for grafana or start port-forwarding:
```bash
kubectl port-forward -n monitoring svc/prometheus-grafana 30000:80
```
Grafana credentials: admin:prom-operator.

## Start collecting metrics
 
Add monitoring role in your db:
```bash
CREATE USER monitoring WITH LOGIN PASSWORD '1234';
ALTER ROLE monitoring SET search_path = mtm, monitoring, pg_catalog, public;
GRANT CONNECT ON DATABASE mydb TO monitoring;
GRANT USAGE ON SCHEMA mtm TO monitoring;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA mtm TO monitoring;
GRANT pg_read_all_settings TO monitoring;
GRANT pg_read_all_stats TO monitoring;
GRANT SELECT ON mtm.cluster_nodes TO monitoring;
```

You need to change **connection URL** for node in [config map](metrics/mtm-metrics-config-cm.yaml)

Create namespace and config that will be mounted to container:
```bash
kubectl create namespace mtm
kubectl apply -f metrics/mtm-metrics-config-cm.yaml -n mtm
```

Deploy mtm-metrics:
```bash
kubectl apply -f metrics/mtm-metrics-deployment.yaml -n mtm
```

Deploy PodMonitor for mtm-metrics if you installed Prometheus stack from above:
```bash
kubectl apply -f metrics/mtm-metrics-pod-monitor.yaml -n mtm
```

**If you edited _config_ and want new changes to take effect:**
```bash
kubectl apply -f metrics/mtm-metrics-config-cm.yaml -n mtm
kubectl rollout restart deployment mtm-metrics-deployment -n mtm
```

## Grafana dashboard

In Grafana GUI you need to isntall plugin **Statusmap** by _flant_.
Import [dashboard](../../metrics/grafana/nodes.json)

